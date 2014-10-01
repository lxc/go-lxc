// Copyright © 2013, 2014, S.Çağlar Onur
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
	"io/ioutil"
	"os/exec"
	"reflect"
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
		return ErrNotDefined
	}

	if flags&isNotDefined != 0 && c.Defined() {
		return ErrAlreadyDefined
	}

	if flags&isRunning != 0 && !c.Running() {
		return ErrNotRunning
	}

	if flags&isNotRunning != 0 && c.Running() {
		return ErrAlreadyRunning
	}
	return nil
}

func (c *Container) cgroupItemAsByteSize(filename string, missing error) (ByteSize, error) {
	size, err := strconv.ParseFloat(c.CgroupItem(filename)[0], 64)
	if err != nil {
		return -1, missing
	}
	return ByteSize(size), nil
}

func (c *Container) setCgroupItemWithByteSize(filename string, limit ByteSize, missing error) error {
	if err := c.SetCgroupItem(filename, fmt.Sprintf("%.f", limit)); err != nil {
		return missing
	}
	return nil
}

// Name returns the name of the container.
func (c *Container) Name() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return C.GoString(c.container.name)
}

// Defined returns true if the container is already defined.
func (c *Container) Defined() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return bool(C.go_lxc_defined(c.container))
}

// Running returns true if the container is already running.
func (c *Container) Running() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return bool(C.go_lxc_running(c.container))
}

// Controllable returns true if the caller can control the container.
func (c *Container) Controllable() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return bool(C.go_lxc_may_control(c.container))
}

// CreateSnapshot creates a new snapshot.
func (c *Container) CreateSnapshot() (*Snapshot, error) {
	if err := c.makeSure(isDefined | isNotRunning); err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	ret := int(C.go_lxc_snapshot(c.container))
	if ret < 0 {
		return nil, ErrCreateSnapshotFailed
	}
	return &Snapshot{Name: fmt.Sprintf("snap%d", ret)}, nil
}

// RestoreSnapshot creates a new container based on a snapshot.
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

	if !bool(C.go_lxc_snapshot_restore(c.container, csnapname, cname)) {
		return ErrRestoreSnapshotFailed
	}
	return nil
}

// DestroySnapshot destroys the specified snapshot.
func (c *Container) DestroySnapshot(snapshot Snapshot) error {
	if err := c.makeSure(isDefined); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	csnapname := C.CString(snapshot.Name)
	defer C.free(unsafe.Pointer(csnapname))

	if !bool(C.go_lxc_snapshot_destroy(c.container, csnapname)) {
		return ErrDestroySnapshotFailed
	}
	return nil
}

// Snapshots returns the list of container snapshots.
func (c *Container) Snapshots() ([]Snapshot, error) {
	if err := c.makeSure(isDefined); err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	var csnapshots *C.struct_lxc_snapshot

	size := int(C.go_lxc_snapshot_list(c.container, &csnapshots))
	defer freeSnapshots(csnapshots, size)

	if size < 1 {
		return nil, ErrNoSnapshot
	}

	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(csnapshots)),
		Len:  size,
		Cap:  size,
	}
	gosnapshots := *(*[]C.struct_lxc_snapshot)(unsafe.Pointer(&hdr))

	snapshots := make([]Snapshot, size, size)
	for i := 0; i < size; i++ {
		snapshots[i] = Snapshot{
			Name:        C.GoString(gosnapshots[i].name),
			Timestamp:   C.GoString(gosnapshots[i].timestamp),
			CommentPath: C.GoString(gosnapshots[i].comment_pathname),
			Path:        C.GoString(gosnapshots[i].lxcpath),
		}
	}

	return snapshots, nil
}

// State returns the state of the container.
func (c *Container) State() State {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return stateMap[C.GoString(C.go_lxc_state(c.container))]
}

// InitPid returns the process ID of the container's init process
// seen from outside the container.
func (c *Container) InitPid() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return int(C.go_lxc_init_pid(c.container))
}

// Daemonize returns true if the container wished to be daemonized.
func (c *Container) Daemonize() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return bool(c.container.daemonize)
}

// WantDaemonize determines if the container wants to run daemonized.
func (c *Container) WantDaemonize(state bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !bool(C.go_lxc_want_daemonize(c.container, C.bool(state))) {
		return ErrDaemonizeFailed
	}
	return nil
}

// WantCloseAllFds determines whether container wishes all file descriptors
// to be closed on startup.
func (c *Container) WantCloseAllFds(state bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !bool(C.go_lxc_want_close_all_fds(c.container, C.bool(state))) {
		return ErrCloseAllFdsFailed
	}
	return nil
}

// SetVerbosity sets the verbosity level of some API calls
func (c *Container) SetVerbosity(verbosity Verbosity) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.verbosity = verbosity
}

// Freeze freezes the running container.
func (c *Container) Freeze() error {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	if c.State() == FROZEN {
		return ErrAlreadyFrozen
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !bool(C.go_lxc_freeze(c.container)) {
		return ErrFreezeFailed
	}

	return nil
}

// Unfreeze thaws the frozen container.
func (c *Container) Unfreeze() error {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	if c.State() != FROZEN {
		return ErrNotFrozen
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !bool(C.go_lxc_unfreeze(c.container)) {
		return ErrUnfreezeFailed
	}

	return nil
}

// CreateUsing creates the container using given arguments with specified backend.
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
		if cargs == nil {
			return ErrAllocationFailed
		}
		defer freeNullTerminatedArgs(cargs, len(args))

		ret = bool(C.go_lxc_create(c.container, ctemplate, cbackend, C.int(c.verbosity), cargs))
	} else {
		ret = bool(C.go_lxc_create(c.container, ctemplate, cbackend, C.int(c.verbosity), nil))
	}

	if !ret {
		return ErrCreateFailed
	}
	return nil
}

// Create creates the container using the Directory backendstore.
func (c *Container) Create(template string, args ...string) error {
	return c.CreateUsing(template, Directory, args...)
}

// CreateAsUser creates the container as "unprivileged user" using the Directory backendstore.
func (c *Container) CreateAsUser(distro string, release string, arch string, args ...string) error {
	// required parameters
	nargs := []string{"-d", distro, "-r", release, "-a", arch}
	// optional arguments
	nargs = append(nargs, args...)

	return c.CreateUsing("download", Directory, nargs...)
}

// Start starts the container.
func (c *Container) Start() error {
	if err := c.makeSure(isDefined | isNotRunning); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !bool(C.go_lxc_start(c.container, 0, nil)) {
		return ErrStartFailed
	}
	return nil
}

// Execute executes the given command in a temporary container.
func (c *Container) Execute(args ...string) ([]byte, error) {
	if err := c.makeSure(isNotDefined); err != nil {
		return nil, err
	}

	cargs := []string{"lxc-execute", "-n", c.Name(), "-P", c.ConfigPath(), "--"}
	cargs = append(cargs, args...)

	c.mu.Lock()
	defer c.mu.Unlock()

	/*
	 * FIXME: Go runtime and src/lxc/start.c signal_handler are not playing nice together so use lxc-execute for now
	 * go-nuts thread: https://groups.google.com/forum/#!msg/golang-nuts/h9GbvfYv83w/5Ly_jvOr86wJ
	 */
	output, err := exec.Command(cargs[0], cargs[1:]...).CombinedOutput()
	if err != nil {
		return nil, ErrExecuteFailed
	}

	return output, nil
	/*
	   cargs := makeNullTerminatedArgs(args)
	   if cargs == nil {
		   return ErrAllocationFailed
	   }
	   defer freeNullTerminatedArgs(cargs, len(args))

	   if !bool(C.go_lxc_start(c.container, 1, cargs)) {
		   return ErrExecuteFailed
	   }
	   return nil
	*/
}

// Stop stops the container.
func (c *Container) Stop() error {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !bool(C.go_lxc_stop(c.container)) {
		return ErrStopFailed
	}
	return nil
}

// Reboot reboots the container.
func (c *Container) Reboot() error {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !bool(C.go_lxc_reboot(c.container)) {
		return ErrRebootFailed
	}
	return nil
}

// Shutdown shuts down the container.
func (c *Container) Shutdown(timeout time.Duration) error {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if !bool(C.go_lxc_shutdown(c.container, C.int(timeout.Seconds()))) {
		return ErrShutdownFailed
	}
	return nil
}

// Destroy destroys the container.
func (c *Container) Destroy() error {
	if err := c.makeSure(isDefined | isNotRunning); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !bool(C.go_lxc_destroy(c.container)) {
		return ErrDestroyFailed
	}
	return nil
}

// CloneUsing clones the container using given arguments with specified backend.
//
// Additional flags to change the cloning behaviour:
// CloneKeepName, CloneKeepMACAddr, CloneSnapshot and CloneMaybeSnapshot
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

	if !bool(C.go_lxc_clone(c.container, cname, C.int(flags), cbackend)) {
		return ErrCloneFailed
	}
	return nil
}

// Clone clones the container using the Directory backendstore.
func (c *Container) Clone(name string) error {
	return c.CloneUsing(name, Directory, 0)
}

// Rename renames the container.
func (c *Container) Rename(name string) error {
	if err := c.makeSure(isDefined | isNotRunning); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	if !bool(C.go_lxc_rename(c.container, cname)) {
		return ErrRenameFailed
	}
	return nil
}

// Wait waits for container to reach a particular state.
func (c *Container) Wait(state State, timeout time.Duration) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	cstate := C.CString(state.String())
	defer C.free(unsafe.Pointer(cstate))

	return bool(C.go_lxc_wait(c.container, cstate, C.int(timeout.Seconds())))
}

// ConfigFileName returns the container's configuration file's name.
func (c *Container) ConfigFileName() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// allocated in lxc.c
	configFileName := C.go_lxc_config_file_name(c.container)
	defer C.free(unsafe.Pointer(configFileName))

	return C.GoString(configFileName)
}

// ConfigItem returns the value of the given config item.
func (c *Container) ConfigItem(key string) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	// allocated in lxc.c
	configItem := C.go_lxc_get_config_item(c.container, ckey)
	defer C.free(unsafe.Pointer(configItem))

	ret := strings.TrimSpace(C.GoString(configItem))
	return strings.Split(ret, "\n")
}

// SetConfigItem sets the value of the given config item.
func (c *Container) SetConfigItem(key string, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	if !bool(C.go_lxc_set_config_item(c.container, ckey, cvalue)) {
		return ErrSettingConfigItemFailed
	}
	return nil
}

// RunningConfigItem returns the value of the given config item.
func (c *Container) RunningConfigItem(key string) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	// allocated in lxc.c
	configItem := C.go_lxc_get_running_config_item(c.container, ckey)
	defer C.free(unsafe.Pointer(configItem))

	ret := strings.TrimSpace(C.GoString(configItem))
	return strings.Split(ret, "\n")
}

// CgroupItem returns the value of the given cgroup subsystem value.
func (c *Container) CgroupItem(key string) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	// allocated in lxc.c
	cgroupItem := C.go_lxc_get_cgroup_item(c.container, ckey)
	defer C.free(unsafe.Pointer(cgroupItem))

	ret := strings.TrimSpace(C.GoString(cgroupItem))
	return strings.Split(ret, "\n")
}

// SetCgroupItem sets the value of given cgroup subsystem value.
func (c *Container) SetCgroupItem(key string, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	if !bool(C.go_lxc_set_cgroup_item(c.container, ckey, cvalue)) {
		return ErrSettingCgroupItemFailed
	}
	return nil
}

// ClearConfig completely clears the containers in-memory configuration.
func (c *Container) ClearConfig() {
	c.mu.Lock()
	defer c.mu.Unlock()

	C.go_lxc_clear_config(c.container)
}

// ClearConfigItem clears the value of given config item.
func (c *Container) ClearConfigItem(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	if !bool(C.go_lxc_clear_config_item(c.container, ckey)) {
		return ErrClearingCgroupItemFailed
	}
	return nil
}

// ConfigKeys returns the names of the config items.
func (c *Container) ConfigKeys(key ...string) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var keys *_Ctype_char

	if key != nil && len(key) == 1 {
		ckey := C.CString(key[0])
		defer C.free(unsafe.Pointer(ckey))

		// allocated in lxc.c
		keys = C.go_lxc_get_keys(c.container, ckey)
		defer C.free(unsafe.Pointer(keys))
	} else {
		// allocated in lxc.c
		keys = C.go_lxc_get_keys(c.container, nil)
		defer C.free(unsafe.Pointer(keys))
	}
	ret := strings.TrimSpace(C.GoString(keys))
	return strings.Split(ret, "\n")
}

// LoadConfigFile loads the configuration file from given path.
func (c *Container) LoadConfigFile(path string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	if !bool(C.go_lxc_load_config(c.container, cpath)) {
		return ErrLoadConfigFailed
	}
	return nil
}

// SaveConfigFile saves the configuration file to given path.
func (c *Container) SaveConfigFile(path string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	if !bool(C.go_lxc_save_config(c.container, cpath)) {
		return ErrSaveConfigFailed
	}
	return nil
}

// ConfigPath returns the configuration file's path.
func (c *Container) ConfigPath() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return C.GoString(C.go_lxc_get_config_path(c.container))
}

// SetConfigPath sets the configuration file's path.
func (c *Container) SetConfigPath(path string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	if !bool(C.go_lxc_set_config_path(c.container, cpath)) {
		return ErrSettingConfigPathFailed
	}
	return nil
}

// MemoryUsage returns memory usage of the container in bytes.
func (c *Container) MemoryUsage() (ByteSize, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return -1, err
	}

	return c.cgroupItemAsByteSize("memory.usage_in_bytes", ErrMemLimit)
}

// MemoryLimit returns memory limit of the container in bytes.
func (c *Container) MemoryLimit() (ByteSize, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return -1, err
	}

	return c.cgroupItemAsByteSize("memory.limit_in_bytes", ErrMemLimit)
}

// SetMemoryLimit sets memory limit of the container in bytes.
func (c *Container) SetMemoryLimit(limit ByteSize) error {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	return c.setCgroupItemWithByteSize("memory.limit_in_bytes", limit, ErrSettingMemoryLimitFailed)
}

// SoftMemoryLimit returns soft memory limit of the container in bytes.
func (c *Container) SoftMemoryLimit() (ByteSize, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return -1, err
	}

	return c.cgroupItemAsByteSize("memory.soft_limit_in_bytes", ErrSoftMemLimit)
}

// SetSoftMemoryLimit sets soft  memory limit of the container in bytes.
func (c *Container) SetSoftMemoryLimit(limit ByteSize) error {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	return c.setCgroupItemWithByteSize("memory.soft_limit_in_bytes", limit, ErrSettingSoftMemoryLimitFailed)
}

// KernelMemoryUsage returns current kernel memory allocation of the container in bytes.
func (c *Container) KernelMemoryUsage() (ByteSize, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return -1, err
	}

	return c.cgroupItemAsByteSize("memory.kmem.usage_in_bytes", ErrKMemLimit)
}

// KernelMemoryLimit returns kernel memory limit of the container in bytes.
func (c *Container) KernelMemoryLimit() (ByteSize, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return -1, err
	}

	return c.cgroupItemAsByteSize("memory.kmem.usage_in_bytes", ErrKMemLimit)
}

// SetKernelMemoryLimit sets kernel memory limit of the container in bytes.
func (c *Container) SetKernelMemoryLimit(limit ByteSize) error {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	return c.setCgroupItemWithByteSize("memory.kmem.limit_in_bytes", limit, ErrSettingKMemoryLimitFailed)
}

// MemorySwapUsage returns memory+swap usage of the container in bytes.
func (c *Container) MemorySwapUsage() (ByteSize, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return -1, err
	}

	return c.cgroupItemAsByteSize("memory.memsw.usage_in_bytes", ErrMemorySwapLimit)
}

// MemorySwapLimit returns the memory+swap limit of the container in bytes.
func (c *Container) MemorySwapLimit() (ByteSize, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return -1, err
	}

	return c.cgroupItemAsByteSize("memory.memsw.limit_in_bytes", ErrMemorySwapLimit)
}

// SetMemorySwapLimit sets memory+swap limit of the container in bytes.
func (c *Container) SetMemorySwapLimit(limit ByteSize) error {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	return c.setCgroupItemWithByteSize("memory.memsw.limit_in_bytes", limit, ErrSettingMemorySwapLimitFailed)
}

// BlkioUsage returns number of bytes transferred to/from the disk by the container.
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
			return ByteSize(blkioUsed), nil
		}
	}
	return -1, ErrBlkioUsage
}

// CPUTime returns the total CPU time (in nanoseconds) consumed by all tasks
// in this cgroup (including tasks lower in the hierarchy).
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
	return time.Duration(cpuUsage), nil
}

// CPUTimePerCPU returns the CPU time (in nanoseconds) consumed on each CPU by
// all tasks in this cgroup (including tasks lower in the hierarchy).
func (c *Container) CPUTimePerCPU() (map[int]time.Duration, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return nil, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	cpuTimes := make(map[int]time.Duration)
	for i, v := range strings.Split(c.CgroupItem("cpuacct.usage_percpu")[0], " ") {
		cpuUsage, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		cpuTimes[i] = time.Duration(cpuUsage)
	}
	return cpuTimes, nil
}

// CPUStats returns the number of CPU cycles (in the units defined by USER_HZ on the system)
// consumed by tasks in this cgroup and its children in both user mode and system (kernel) mode.
func (c *Container) CPUStats() (map[string]int64, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return nil, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	cpuStat := c.CgroupItem("cpuacct.stat")
	user, err := strconv.ParseInt(strings.Split(cpuStat[0], "user ")[1], 10, 64)
	if err != nil {
		return nil, err
	}
	system, err := strconv.ParseInt(strings.Split(cpuStat[1], "system ")[1], 10, 64)
	if err != nil {
		return nil, err
	}

	return map[string]int64{"user": user, "system": system}, nil
}

// ConsoleFd allocates a console tty from container
// ttynum: tty number to attempt to allocate or -1 to allocate the first available tty
//
// Returns "ttyfd" on success, -1 on failure. The returned "ttyfd" is
// used to keep the tty allocated. The caller should close "ttyfd" to
// indicate that it is done with the allocated console so that it can
// be allocated by another caller.
func (c *Container) ConsoleFd(ttynum int) (int, error) {
	// FIXME: Make idiomatic
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return -1, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	ret := int(C.go_lxc_console_getfd(c.container, C.int(ttynum)))
	if ret < 0 {
		return -1, ErrAttachFailed
	}
	return ret, nil
}

// Console allocates and runs a console tty from container
// ttynum: tty number to attempt to allocate, -1 to allocate the first available tty, or 0 to allocate the console
// stdinfd: fd to read input from
// stdoutfd: fd to write output to
// stderrfd: fd to write error output to
// escape: the escape character (1 == 'a', 2 == 'b', ...)
//
// This function will not return until the console has been exited by the user.
func (c *Container) Console(ttynum int, stdinfd, stdoutfd, stderrfd uintptr, escape int) error {
	// FIXME: Make idiomatic
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !bool(C.go_lxc_console(c.container, C.int(ttynum), C.int(stdinfd),
		C.int(stdoutfd), C.int(stderrfd), C.int(escape))) {
		return ErrAttachFailed
	}
	return nil
}

// AttachShell attaches a shell to the container.
func (c *Container) AttachShell() error {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if int(C.go_lxc_attach(c.container, false)) < 0 {
		return ErrAttachFailed
	}
	return nil
}

// AttachShellWithClearEnvironment attaches a shell to the container.
// It clears all environment variables before attaching.
func (c *Container) AttachShellWithClearEnvironment() error {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if int(C.go_lxc_attach(c.container, true)) < 0 {
		return ErrAttachFailed
	}
	return nil
}

type Command struct {
	Stdinfd uintptr
	Stdoutfd uintptr
	Stderrfd uintptr
	Args []string
	Env []string
	Cwd string
	ClearEnv bool
}

func (cmd *Command) Run(c *Container) (bool, error) {
	if len(cmd.Args) == 0 {
		return false, ErrInsufficientNumberOfArguments
	}

	if err := c.makeSure(isDefined | isRunning); err != nil {
		return false, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	cargs := makeNullTerminatedArgs(cmd.Args)
	if cargs == nil {
		return false, ErrAllocationFailed
	}
	defer freeNullTerminatedArgs(cargs, len(cmd.Args))

	cenv := makeNullTerminatedArgs(cmd.Env)
	if cenv == nil {
		return false, ErrAllocationFailed
	}
	defer freeNullTerminatedArgs(cenv, len(cmd.Env))

	cwd := C.CString(cmd.Cwd)
	defer C.free(unsafe.Pointer(cwd))

	ret := int(C.go_lxc_attach_run_wait(
		c.container,
		C.bool(cmd.ClearEnv),
		C.int(cmd.Stdinfd),
		C.int(cmd.Stdoutfd),
		C.int(cmd.Stderrfd),
		*cwd,
		*cenv,
		cargs,
	))

	if ret < 0 {
		return false, ErrAttachFailed
	}
	return ret == 0, nil
}

// RunCommand runs the user specified command inside the container and waits for it to exit.
// stdinfd: fd to read input from
// stdoutfd: fd to write output to
// stderrfd: fd to write error output to
func (c *Container) RunCommand(stdinfd, stdoutfd, stderrfd uintptr, args ...string) (bool, error) {
	cmd := &Command{
		Stdinfd: stdinfd,
		Stdoutfd: stdoutfd,
		Stderrfd: stderrfd,
		Args: args,
		ClearEnv: false,
	}
	return cmd.Run(c)
}

// RunCommandWithClearEnvironment runs the user specified command inside the container
// and waits for it to exit. It clears all environment variables before running.
// stdinfd: fd to read input from
// stdoutfd: fd to write output to
// stderrfd: fd to write error output to
func (c *Container) RunCommandWithClearEnvironment(stdinfd, stdoutfd, stderrfd uintptr, args ...string) (bool, error) {
	cmd := &Command{
		Stdinfd: stdinfd,
		Stdoutfd: stdoutfd,
		Stderrfd: stderrfd,
		Args: args,
		ClearEnv: true,
	}
	return cmd.Run(c)
}

// Interfaces returns the names of the network interfaces.
func (c *Container) Interfaces() ([]string, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return nil, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	result := C.go_lxc_get_interfaces(c.container)
	if result == nil {
		return nil, ErrInterfaces
	}
	return convertArgs(result), nil
}

// InterfaceStats returns the stats about container's network interfaces
func (c *Container) InterfaceStats() (map[string]map[string]ByteSize, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return nil, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	var interfaceName string

	statistics := make(map[string]map[string]ByteSize)

	for i := 0; i < len(c.ConfigItem("lxc.network")); i++ {
		interfaceType := c.RunningConfigItem(fmt.Sprintf("lxc.network.%d.type", i))
		if interfaceType == nil {
			continue
		}

		if interfaceType[0] == "veth" {
			interfaceName = c.RunningConfigItem(fmt.Sprintf("lxc.network.%d.veth.pair", i))[0]
		} else {
			interfaceName = c.RunningConfigItem(fmt.Sprintf("lxc.network.%d.link", i))[0]
		}

		for _, v := range []string{"rx", "tx"} {
			/* tx and rx are reversed from the host vs container */
			content, err := ioutil.ReadFile(fmt.Sprintf("/sys/class/net/%s/statistics/%s_bytes", interfaceName, v))
			if err != nil {
				return nil, err
			}

			bytes, err := strconv.ParseInt(strings.Split(string(content), "\n")[0], 10, 64)
			if err != nil {
				return nil, err
			}

			if statistics[interfaceName] == nil {
				statistics[interfaceName] = make(map[string]ByteSize)
			}
			statistics[interfaceName][v] = ByteSize(bytes)
		}
	}

	return statistics, nil
}

// IPAddress returns the IP address of the given network interface.
func (c *Container) IPAddress(interfaceName string) ([]string, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return nil, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	cinterface := C.CString(interfaceName)
	defer C.free(unsafe.Pointer(cinterface))

	result := C.go_lxc_get_ips(c.container, cinterface, nil, 0)
	if result == nil {
		return nil, ErrIPAddress
	}
	return convertArgs(result), nil
}

// IPAddresses returns all IP addresses.
func (c *Container) IPAddresses() ([]string, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return nil, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	result := C.go_lxc_get_ips(c.container, nil, nil, 0)
	if result == nil {
		return nil, ErrIPAddresses
	}
	return convertArgs(result), nil

}

// IPv4Addresses returns all IPv4 addresses.
func (c *Container) IPv4Addresses() ([]string, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return nil, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	cfamily := C.CString("inet")
	defer C.free(unsafe.Pointer(cfamily))

	result := C.go_lxc_get_ips(c.container, nil, cfamily, 0)
	if result == nil {
		return nil, ErrIPv4Addresses
	}
	return convertArgs(result), nil
}

// IPv6Addresses returns all IPv6 addresses.
func (c *Container) IPv6Addresses() ([]string, error) {
	if err := c.makeSure(isDefined | isRunning); err != nil {
		return nil, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	cfamily := C.CString("inet6")
	defer C.free(unsafe.Pointer(cfamily))

	result := C.go_lxc_get_ips(c.container, nil, cfamily, 0)
	if result == nil {
		return nil, ErrIPv6Addresses
	}
	return convertArgs(result), nil
}

// LogFile returns the name of the logfile.
func (c *Container) LogFile() string {
	return c.ConfigItem("lxc.logfile")[0]
}

// SetLogFile sets the name of the logfile.
func (c *Container) SetLogFile(filename string) error {
	if err := c.SetConfigItem("lxc.logfile", filename); err != nil {
		return err
	}
	return nil
}

// LogLevel returns the level of the logfile.
func (c *Container) LogLevel() LogLevel {
	return logLevelMap[c.ConfigItem("lxc.loglevel")[0]]
}

// SetLogLevel sets the level of the logfile.
func (c *Container) SetLogLevel(level LogLevel) error {
	if err := c.SetConfigItem("lxc.loglevel", level.String()); err != nil {
		return err
	}
	return nil
}

// AddDeviceNode adds specified device to the container.
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

		if !bool(C.go_lxc_add_device_node(c.container, csource, cdestination)) {
			return ErrAddDeviceNodeFailed
		}
		return nil
	}

	if !bool(C.go_lxc_add_device_node(c.container, csource, nil)) {
		return ErrAddDeviceNodeFailed
	}
	return nil

}

// RemoveDeviceNode removes the specified device from the container.
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

		if !bool(C.go_lxc_remove_device_node(c.container, csource, cdestination)) {
			return ErrRemoveDeviceNodeFailed
		}
		return nil
	}

	if !bool(C.go_lxc_remove_device_node(c.container, csource, nil)) {
		return ErrRemoveDeviceNodeFailed
	}
	return nil
}
