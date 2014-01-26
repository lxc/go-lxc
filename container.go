// Copyright © 2013, S.Çağlar Onur
// Use of this source code is governed by a LGPLv2.1
// license that can be found in the LICENSE file.
//
// Authors:
// S.Çağlar Onur <caglar@10ur.org>

// +build linux

package lxc

// #include <lxc/lxccontainer.h>
// #include "lxc.h"
import "C"

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"
)

// Container struct
type Container struct {
	container *C.struct_lxc_container
	mu        sync.RWMutex

	name      string
	verbosity Verbosity
}

// Snapshot struct
type Snapshot struct {
	Name        string
	CommentPath string
	Timestamp   string
	Path        string
}

func (c *Container) makeSure(flags int) error {
	if flags&isDefined != 0 && !c.Defined() {
		return fmt.Errorf(errNotDefined, c.name)
	}

	if flags&isNotDefined != 0 && c.Defined() {
		return fmt.Errorf(errAlreadyDefined, c.name)
	}

	if flags&isRunning != 0 && !c.Running() {
		return fmt.Errorf(errNotRunning, c.name)
	}

	if flags&isNotRunning != 0 && c.Running() {
		return fmt.Errorf(errAlreadyRunning, c.name)
	}
	return nil
}

// Name returns the name of the container
func (c *Container) Name() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return C.GoString(c.container.name)
}

// Defined returns true if the container is already defined
func (c *Container) Defined() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return bool(C.lxc_defined(c.container))
}

// Running returns true if the container is already running
func (c *Container) Running() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return bool(C.lxc_running(c.container))
}

// Controllable returns true if the caller may control the container
func (c *Container) Controllable() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return bool(C.lxc_may_control(c.container))
}

// CreateSnapshot creates a new snapshot
func (c *Container) CreateSnapshot() (*Snapshot, error) {
	if err := c.makeSure(isDefined | isNotRunning); err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	ret := int(C.lxc_snapshot(c.container))
	if ret < 0 {
		return nil, fmt.Errorf(errCreateSnapshotFailed, c.name)
	}
	return &Snapshot{Name: fmt.Sprintf("snap%d", ret)}, nil
}

// RestoreSnapshot creates a new container based on a snapshot
func (c *Container) RestoreSnapshot(snapshot Snapshot, name string) error {
	if err := c.makeSure(isDefined); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	csnapname := C.CString(snapshot.Name)
	defer C.free(unsafe.Pointer(csnapname))

	if !bool(C.lxc_snapshot_restore(c.container, csnapname, cname)) {
		return fmt.Errorf(errRestoreSnapshotFailed, c.name)
	}
	return nil
}

// DestroySnapshot destroys the specified snapshot
func (c *Container) DestroySnapshot(snapshot Snapshot) error {
	if err := c.makeSure(isDefined); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	csnapname := C.CString(snapshot.Name)
	defer C.free(unsafe.Pointer(csnapname))

	if !bool(C.lxc_snapshot_destroy(c.container, csnapname)) {
		return fmt.Errorf(errDestroySnapshotFailed, c.name)
	}
	return nil
}

// Snapshots returns the list of container snapshots
func (c *Container) Snapshots() ([]Snapshot, error) {
	if err := c.makeSure(isDefined); err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	var csnapshots *C.struct_lxc_snapshot
	var snapshots []Snapshot

	size := int(C.lxc_snapshot_list(c.container, &csnapshots))
	defer freeSnapshots(csnapshots, size)

	if size < 1 {
		return nil, fmt.Errorf("%s has no snapshots", c.name)
	}

	p := uintptr(unsafe.Pointer(csnapshots))
	for i := 0; i < size; i++ {
		z := (*C.struct_lxc_snapshot)(unsafe.Pointer(p))
		s := &Snapshot{Name: C.GoString(z.name), Timestamp: C.GoString(z.timestamp),
			CommentPath: C.GoString(z.comment_pathname), Path: C.GoString(z.lxcpath)}
		snapshots = append(snapshots, *s)
		p += unsafe.Sizeof(*csnapshots)
	}
	return snapshots, nil
}

// State returns the state of the container
func (c *Container) State() State {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return stateMap[C.GoString(C.lxc_state(c.container))]
}

// InitPID returns the process ID of the container's init process seen from outside the container
func (c *Container) InitPID() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return int(C.lxc_init_pid(c.container))
}

// Daemonize returns whether the container wished to be daemonized
func (c *Container) Daemonize() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return bool(c.container.daemonize)
}

// WantDaemonize sets the daemonize flag for the container
func (c *Container) WantDaemonize(state bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !bool(C.lxc_want_daemonize(c.container, C.bool(state))) {
		return fmt.Errorf(errDaemonizeFailed, c.name)
	}
	return nil
}

// WantCloseAllFds sets the close_all_fds flag for the container
func (c *Container) WantCloseAllFds(state bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !bool(C.lxc_want_close_all_fds(c.container, C.bool(state))) {
		return fmt.Errorf(errCloseAllFdsFailed, c.name)
	}
	return nil
}

// SetVerbosity sets the verbosity level of some API calls
func (c *Container) SetVerbosity(verbosity Verbosity) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.verbosity = verbosity
}

// Freeze freezes the running container
func (c *Container) Freeze() error {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	if c.State() == FROZEN {
		return fmt.Errorf(errAlreadyFrozen, c.name)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !bool(C.lxc_freeze(c.container)) {
		return fmt.Errorf(errFreezeFailed, c.name)
	}

	return nil
}

// Unfreeze thaws the frozen container
func (c *Container) Unfreeze() error {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	if c.State() != FROZEN {
		return fmt.Errorf(errNotFrozen, c.name)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !bool(C.lxc_unfreeze(c.container)) {
		return fmt.Errorf(errUnfreezeFailed, c.name)
	}

	return nil
}

// CreateUsing creates the container using given template and arguments with specified backend
func (c *Container) CreateUsing(template string, backend BackendStore, args ...string) error {
	// FIXME: Support bdevtype and bdev_specs
	// bdevtypes:
	// "btrfs", "zfs", "lvm", "dir"
	//
	// best tries to find the best backing store type
	//
	// bdev_specs:
	// zfs requires zfsroot
	// lvm requires lvname/vgname/thinpool as well as fstype and fssize
	// btrfs requires nothing
	// dir requires nothing
	if err := c.makeSure(isNotDefined); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	ctemplate := C.CString(template)
	defer C.free(unsafe.Pointer(ctemplate))

	cbackend := C.CString(backend.String())
	defer C.free(unsafe.Pointer(cbackend))

	ret := false
	if args != nil {
		cargs := makeNullTerminatedArgs(args)
		defer freeNullTerminatedArgs(cargs, len(args))

		ret = bool(C.lxc_create(c.container, ctemplate, cbackend, C.int(c.verbosity), cargs))
	} else {
		ret = bool(C.lxc_create(c.container, ctemplate, cbackend, C.int(c.verbosity), nil))
	}

	if !ret {
		return fmt.Errorf(errCreateFailed, c.name)
	}
	return nil
}

// Create creates the container using given template and arguments with Directory backend
func (c *Container) Create(template string, args ...string) error {
	return c.CreateUsing(template, Directory, args...)
}

// CreateAsUser creates the container using given template and arguments with Directory backend as an unprivileged user
func (c *Container) CreateAsUser(distro string, release string, arch string, args ...string) error {
	// required parameters
	nargs := []string{"-d", distro, "-r", release, "-a", arch}
	// optional arguments
	nargs = append(nargs, args...)

	return c.CreateUsing("download", Directory, nargs...)
}

// Start starts the container
func (c *Container) Start() error {
	if err := c.makeSure(isDefined | isNotRunning); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !bool(C.lxc_start(c.container, 0, nil)) {
		return fmt.Errorf(errStartFailed, c.name)
	}
	return nil
}

// Execute executes the given command in a temporary container
func (c *Container) Execute(args ...string) ([]byte, error) {
	if err := c.makeSure(isNotDefined); err != nil {
		return nil, err
	}

	cargs := []string{"lxc-execute", "-n", c.Name(), "-P", c.ConfigPath(), "--"}
	cargs = append(cargs, args...)

	c.mu.Lock()
	defer c.mu.Unlock()

	/*
	 * FIXME: Go runtime and src/c/start.c signal_handler are not playing nice together so use c-execute for now
	 */
	output, err := exec.Command(cargs[0], cargs[1:]...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf(errExecuteFailed, c.name)
	}

	return output, nil
	/*
		cargs := makeNullTerminatedArgs(args)
		defer freeNullTerminatedArgs(cargs, len(args))

		if !bool(C.lxc_start(c.container, 1, cargs)) {
			return fmt.Errorf(errExecuteFailed, c.name)
		}
		return nil
	*/
}

// Stop stops the container
func (c *Container) Stop() error {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !bool(C.lxc_stop(c.container)) {
		return fmt.Errorf(errStopFailed, c.name)
	}
	return nil
}

// Reboot reboots the container
func (c *Container) Reboot() error {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !bool(C.lxc_reboot(c.container)) {
		return fmt.Errorf(errRebootFailed, c.name)
	}
	return nil
}

// Shutdown shutdowns the container
func (c *Container) Shutdown(timeout int) error {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !bool(C.lxc_shutdown(c.container, C.int(timeout))) {
		return fmt.Errorf(errShutdownFailed, c.name)
	}
	return nil
}

// Destroy destroys the container
func (c *Container) Destroy() error {
	if err := c.makeSure(isDefined | isNotRunning); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !bool(C.lxc_destroy(c.container)) {
		return fmt.Errorf(errDestroyFailed, c.name)
	}
	return nil
}

// CloneUsing clones the container using specified backend
func (c *Container) CloneUsing(name string, backend BackendStore, flags CloneFlags) error {
	// FIXME: support lxcpath, bdevtype, bdevdata, newsize and hookargs
	//
	// bdevtypes:
	// "btrfs", "zfs", "lvm", "dir" "overlayfs"
	//
	// bdevdata:
	// zfs requires zfsroot
	// lvm requires lvname/vgname/thinpool as well as fstype and fssize
	// btrfs requires nothing
	// dir requires nothing
	//
	// newsize: for blockdev-backed backingstores
	//
	// hookargs: additional arguments to pass to the clone hook script
	if err := c.makeSure(isDefined | isNotRunning); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cbackend := C.CString(backend.String())
	defer C.free(unsafe.Pointer(cbackend))

	if !bool(C.lxc_clone(c.container, cname, C.int(flags), cbackend)) {
		return fmt.Errorf(errCloneFailed, c.name)
	}
	return nil
}

// Clone clones the container using the Directory backendstore
func (c *Container) Clone(name string) error {
	return c.CloneUsing(name, Directory, 0)
}

// Rename renames the container
func (c *Container) Rename(name string) error {
	if err := c.makeSure(isDefined | isNotRunning); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	if !bool(C.lxc_rename(c.container, cname)) {
		return fmt.Errorf(errRenameFailed, c.name)
	}
	return nil
}

// Wait waits for container to reach a given state or timeouts
func (c *Container) Wait(state State, timeout int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	cstate := C.CString(state.String())
	defer C.free(unsafe.Pointer(cstate))

	return bool(C.lxc_wait(c.container, cstate, C.int(timeout)))
}

// ConfigFileName returns the container's configuration file's name
func (c *Container) ConfigFileName() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// allocated in c.c
	configFileName := C.lxc_config_file_name(c.container)
	defer C.free(unsafe.Pointer(configFileName))

	return C.GoString(configFileName)
}

// ConfigItem returns the value of the given config item
func (c *Container) ConfigItem(key string) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	// allocated in c.c
	configItem := C.lxc_get_config_item(c.container, ckey)
	defer C.free(unsafe.Pointer(configItem))

	ret := strings.TrimSpace(C.GoString(configItem))
	return strings.Split(ret, "\n")
}

// SetConfigItem sets the value of the given config item
func (c *Container) SetConfigItem(key string, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	if !bool(C.lxc_set_config_item(c.container, ckey, cvalue)) {
		return fmt.Errorf(errSettingConfigItemFailed, c.name, key, value)
	}
	return nil
}

// CgroupItem returns the value of the given cgroup subsystem value
func (c *Container) CgroupItem(key string) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	// allocated in c.c
	cgroupItem := C.lxc_get_cgroup_item(c.container, ckey)
	defer C.free(unsafe.Pointer(cgroupItem))

	ret := strings.TrimSpace(C.GoString(cgroupItem))
	return strings.Split(ret, "\n")
}

// SetCgroupItem sets the value of given cgroup subsystem value
func (c *Container) SetCgroupItem(key string, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	if !bool(C.lxc_set_cgroup_item(c.container, ckey, cvalue)) {
		return fmt.Errorf(errSettingCgroupItemFailed, c.name, key, value)
	}
	return nil
}

// ClearConfigItem clears the value of given config item
func (c *Container) ClearConfigItem(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	if !bool(C.lxc_clear_config_item(c.container, ckey)) {
		return fmt.Errorf(errClearingCgroupItemFailed, c.name, key)
	}
	return nil
}

// ConfigKeys returns the names of the config items
func (c *Container) ConfigKeys(key ...string) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var keys *_Ctype_char

	if key != nil && len(key) == 1 {
		ckey := C.CString(key[0])
		defer C.free(unsafe.Pointer(ckey))

		// allocated in c.c
		keys = C.lxc_get_keys(c.container, ckey)
		defer C.free(unsafe.Pointer(keys))
	} else {
		// allocated in c.c
		keys = C.lxc_get_keys(c.container, nil)
		defer C.free(unsafe.Pointer(keys))
	}
	ret := strings.TrimSpace(C.GoString(keys))
	return strings.Split(ret, "\n")
}

// LoadConfigFile loads the configuration file from given path
func (c *Container) LoadConfigFile(path string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	if !bool(C.lxc_load_config(c.container, cpath)) {
		return fmt.Errorf(errLoadConfigFailed, c.name, path)
	}
	return nil
}

// SaveConfigFile saves the configuration file to given path
func (c *Container) SaveConfigFile(path string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	if !bool(C.lxc_save_config(c.container, cpath)) {
		return fmt.Errorf(errSaveConfigFailed, c.name, path)
	}
	return nil
}

// ConfigPath returns the configuration file's path
func (c *Container) ConfigPath() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return C.GoString(C.lxc_get_config_path(c.container))
}

// SetConfigPath sets the configuration file's path
func (c *Container) SetConfigPath(path string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	if !bool(C.lxc_set_config_path(c.container, cpath)) {
		return fmt.Errorf(errSettingConfigPathFailed, c.name, path)
	}
	return nil
}

// MemoryUsage returns memory usage of the container in bytes
func (c *Container) MemoryUsage() (ByteSize, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return -1, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	memUsed, err := strconv.ParseFloat(c.CgroupItem("memory.usage_in_bytes")[0], 64)
	if err != nil {
		return -1, fmt.Errorf(errMemLimit)
	}
	return ByteSize(memUsed), err
}

// KernelMemoryUsage returns current kernel memory allocation of the container in bytes
func (c *Container) KernelMemoryUsage() (ByteSize, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return -1, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	kmemUsed, err := strconv.ParseFloat(c.CgroupItem("memory.kmem.usage_in_bytes")[0], 64)
	if err != nil {
		return -1, fmt.Errorf(errKMemLimit)
	}
	return ByteSize(kmemUsed), err
}

// SwapUsage returns swap usage of the container in bytes
func (c *Container) SwapUsage() (ByteSize, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return -1, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	swapUsed, err := strconv.ParseFloat(c.CgroupItem("memory.memsw.usage_in_bytes")[0], 64)
	if err != nil {
		return -1, fmt.Errorf(errSwapLimit)
	}
	return ByteSize(swapUsed), err
}

// BlkioUsage returns number of bytes transferred to/from the disk by the container
func (c *Container) BlkioUsage() (ByteSize, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return -1, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, v := range c.CgroupItem("blkio.throttle.io_service_bytes") {
		b := strings.Split(v, " ")
		if b[0] == "Total" {
			blkioUsed, err := strconv.ParseFloat(b[1], 64)
			if err != nil {
				return -1, err
			}
			return ByteSize(blkioUsed), err
		}
	}
	return -1, fmt.Errorf(errBlkioUsage, c.name)
}

// MemoryLimit returns memory limit of the container in bytes
func (c *Container) MemoryLimit() (ByteSize, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return -1, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	memLimit, err := strconv.ParseFloat(c.CgroupItem("memory.limit_in_bytes")[0], 64)
	if err != nil {
		return -1, fmt.Errorf(errMemLimit)
	}
	return ByteSize(memLimit), err
}

// SetMemoryLimit sets memory limit of the container in bytes
func (c *Container) SetMemoryLimit(limit ByteSize) error {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	if err := c.SetCgroupItem("memory.limit_in_bytes", fmt.Sprintf("%.f", limit)); err != nil {
		return fmt.Errorf(errSettingMemoryLimitFailed, c.name)
	}
	return nil
}

// KernelMemoryLimit returns kernel memory limit of the container in bytes
func (c *Container) KernelMemoryLimit() (ByteSize, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return -1, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	kmemLimit, err := strconv.ParseFloat(c.CgroupItem("memory.kmem.limit_in_bytes")[0], 64)
	if err != nil {
		return -1, fmt.Errorf(errKMemLimit)
	}
	return ByteSize(kmemLimit), err
}

// SetKernelMemoryLimit sets kernel memory limit of the container in bytes
func (c *Container) SetKernelMemoryLimit(limit ByteSize) error {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	if err := c.SetCgroupItem("memory.kmem.limit_in_bytes", fmt.Sprintf("%.f", limit)); err != nil {
		return fmt.Errorf(errSettingKMemoryLimitFailed, c.name)
	}
	return nil
}

// SwapLimit returns the swap limit of the container in bytes
func (c *Container) SwapLimit() (ByteSize, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return -1, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	swapLimit, err := strconv.ParseFloat(c.CgroupItem("memory.memsw.limit_in_bytes")[0], 64)
	if err != nil {
		return -1, fmt.Errorf(errSwapLimit)
	}
	return ByteSize(swapLimit), err
}

// SetSwapLimit sets memory limit of the container in bytes
func (c *Container) SetSwapLimit(limit ByteSize) error {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	if err := c.SetCgroupItem("memory.memsw.limit_in_bytes", fmt.Sprintf("%.f", limit)); err != nil {
		return fmt.Errorf(errSettingSwapLimitFailed, c.name)
	}
	return nil
}

// CPUTime returns the total CPU time (in nanoseconds) consumed by all tasks in this cgroup (including tasks lower in the hierarchy).
func (c *Container) CPUTime() (time.Duration, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return -1, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	cpuUsage, err := strconv.ParseInt(c.CgroupItem("cpuacct.usage")[0], 10, 64)
	if err != nil {
		return -1, err
	}
	return time.Duration(cpuUsage), err
}

// CPUTimePerCPU returns the CPU time (in nanoseconds) consumed on each CPU by
// all tasks in this cgroup (including tasks lower in the hierarchy).
func (c *Container) CPUTimePerCPU() ([]time.Duration, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return nil, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	var cpuTimes []time.Duration

	for _, v := range strings.Split(c.CgroupItem("cpuacct.usage_percpu")[0], " ") {
		cpuUsage, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		cpuTimes = append(cpuTimes, time.Duration(cpuUsage))
	}
	return cpuTimes, nil
}

// CPUStats returns the number of CPU cycles (in the units defined by USER_HZ on the system)
// consumed by tasks in this cgroup and its children in both user mode and system (kernel) mode.
func (c *Container) CPUStats() ([]int64, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return nil, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	cpuStat := c.CgroupItem("cpuacct.stat")
	user, err := strconv.ParseInt(strings.Split(cpuStat[0], " ")[1], 10, 64)
	if err != nil {
		return nil, err
	}
	system, err := strconv.ParseInt(strings.Split(cpuStat[1], " ")[1], 10, 64)
	if err != nil {
		return nil, err
	}
	return []int64{user, system}, nil
}

// ConsoleGetFD allocates a console tty from container
// ttynum: tty number to attempt to allocate or -1 to allocate the first available tty
//
// Returns "ttyfd" on success, -1 on failure. The returned "ttyfd" is
// used to keep the tty allocated. The caller should close "ttyfd" to
// indicate that it is done with the allocated console so that it can
// be allocated by another caller.
func (c *Container) ConsoleGetFD(ttynum int) (int, error) {
	// FIXME: Make idiomatic
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return -1, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	ret := int(C.lxc_console_getfd(c.container, C.int(ttynum)))
	if ret < 0 {
		return -1, fmt.Errorf(errAttachFailed, c.name)
	}
	return ret, nil
}

// Console allocates and runs a console tty from container
// ttynum: tty number to attempt to allocate, -1 to allocate the first available tty, or 0 to allocate the console
// stdinfd: fd to read input from
// stdoutfd: fd to write output to
// stderrfd: fd to write error output to
// escape: he escape character (1 == 'a', 2 == 'b', ...)
//
// This function will not return until the console has been exited by the user.
func (c *Container) Console(ttynum, stdinfd, stdoutfd, stderrfd, escape int) error {
	// FIXME: Make idiomatic
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !bool(C.lxc_console(c.container, C.int(ttynum), C.int(stdinfd),
		C.int(stdoutfd), C.int(stderrfd), C.int(escape))) {
		return fmt.Errorf(errAttachFailed, c.name)
	}
	return nil
}

// AttachShell runs a shell inside the container
func (c *Container) AttachShell() error {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if int(C.lxc_attach(c.container, false)) < 0 {
		return fmt.Errorf(errAttachFailed, c.name)
	}
	return nil
}

// AttachShellWithClearEnvironment runs a shell inside the container and
// clears all environment variables before attaching
func (c *Container) AttachShellWithClearEnvironment() error {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if int(C.lxc_attach(c.container, true)) < 0 {
		return fmt.Errorf(errAttachFailed, c.name)
	}
	return nil
}

// RunCommand runs the user specified command inside the container and waits it to exit
func (c *Container) RunCommand(args ...string) error {
	if args == nil {
		return fmt.Errorf(errInsufficientNumberOfArguments)
	}

	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	cargs := makeNullTerminatedArgs(args)
	defer freeNullTerminatedArgs(cargs, len(args))

	if int(C.lxc_attach_run_wait(c.container, false, cargs)) < 0 {
		return fmt.Errorf(errAttachFailed, c.name)
	}
	return nil
}

// RunCommandWithClearEnvironment runs the user specified command inside the container
// and waits it to exit. It also clears all environment variables before attaching.
func (c *Container) RunCommandWithClearEnvironment(args ...string) error {
	if args == nil {
		return fmt.Errorf(errInsufficientNumberOfArguments)
	}

	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	cargs := makeNullTerminatedArgs(args)
	defer freeNullTerminatedArgs(cargs, len(args))

	if int(C.lxc_attach_run_wait(c.container, true, cargs)) < 0 {
		return fmt.Errorf(errAttachFailed, c.name)
	}
	return nil
}

// Interfaces returns the name of the network interfaces from the container
func (c *Container) Interfaces() ([]string, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return nil, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	result := C.lxc_get_interfaces(c.container)
	if result == nil {
		return nil, fmt.Errorf(errInterfaces, c.name)
	}
	return convertArgs(result), nil
}

// IPAddress returns the IP address of the given network interface from the container
func (c *Container) IPAddress(interfaceName string) ([]string, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return nil, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	cinterface := C.CString(interfaceName)
	defer C.free(unsafe.Pointer(cinterface))

	result := C.lxc_get_ips(c.container, cinterface, nil, 0)
	if result == nil {
		return nil, fmt.Errorf(errIPAddress, interfaceName, c.name)
	}
	return convertArgs(result), nil
}

// IPAddresses returns all IP addresses from the container
func (c *Container) IPAddresses() ([]string, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return nil, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	result := C.lxc_get_ips(c.container, nil, nil, 0)
	if result == nil {
		return nil, fmt.Errorf(errIPAddresses, c.name)
	}
	return convertArgs(result), nil

}

// IPv4Addresses returns all IPv4 addresses from the container
func (c *Container) IPv4Addresses() ([]string, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return nil, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	cfamily := C.CString("inet")
	defer C.free(unsafe.Pointer(cfamily))

	result := C.lxc_get_ips(c.container, nil, cfamily, 0)
	if result == nil {
		return nil, fmt.Errorf(errIPv4Addresses, c.name)
	}
	return convertArgs(result), nil
}

// IPv6Addresses returns all IPv6 addresses from the container
func (c *Container) IPv6Addresses() ([]string, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return nil, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	cfamily := C.CString("inet6")
	defer C.free(unsafe.Pointer(cfamily))

	result := C.lxc_get_ips(c.container, nil, cfamily, 0)
	if result == nil {
		return nil, fmt.Errorf(errIPv6Addresses, c.name)
	}
	return convertArgs(result), nil
}

// LogFile returns the name of the logfile
func (c *Container) LogFile() string {
	return c.ConfigItem("lxc.logfile")[0]
}

// SetLogFile sets the logfile to given filename
func (c *Container) SetLogFile(filename string) error {
	if err := c.SetConfigItem("lxc.logfile", filename); err != nil {
		return err
	}
	return nil
}

// LogLevel returns the name of the logfile
func (c *Container) LogLevel() LogLevel {
	return logLevelMap[c.ConfigItem("lxc.loglevel")[0]]
}

// SetLogLevel sets the logfile to given filename
func (c *Container) SetLogLevel(level LogLevel) error {
	if err := c.SetConfigItem("lxc.loglevel", level.String()); err != nil {
		return err
	}
	return nil
}

// AddDeviceNode adds specified device to the container
func (c *Container) AddDeviceNode(source string, destination ...string) error {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	csource := C.CString(source)
	defer C.free(unsafe.Pointer(csource))

	if destination != nil && len(destination) == 1 {
		cdestination := C.CString(destination[0])
		defer C.free(unsafe.Pointer(cdestination))

		if !bool(C.lxc_add_device_node(c.container, csource, cdestination)) {
			return fmt.Errorf(errAddDeviceNodeFailed, source, c.name)
		}
		return nil
	}

	if !bool(C.lxc_add_device_node(c.container, csource, nil)) {
		return fmt.Errorf(errAddDeviceNodeFailed, source, c.name)
	}
	return nil

}

// RemoveDeviceNode removes the specified device from the container
func (c *Container) RemoveDeviceNode(source string, destination ...string) error {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	csource := C.CString(source)
	defer C.free(unsafe.Pointer(csource))

	if destination != nil && len(destination) == 1 {
		cdestination := C.CString(destination[0])
		defer C.free(unsafe.Pointer(cdestination))

		if !bool(C.lxc_remove_device_node(c.container, csource, cdestination)) {
			return fmt.Errorf(errRemoveDeviceNodeFailed, source, c.name)
		}
		return nil
	}

	if !bool(C.lxc_remove_device_node(c.container, csource, nil)) {
		return fmt.Errorf(errRemoveDeviceNodeFailed, source, c.name)
	}
	return nil
}
