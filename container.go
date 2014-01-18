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
	verbosity Verbosity
	name      string
	mu        sync.RWMutex
}

// Snapshot struct
type Snapshot struct {
	Name        string
	CommentPath string
	Timestamp   string
	Path        string
}

func (lxc *Container) ensureDefinedAndRunning() error {
	if !lxc.Defined() {
		return fmt.Errorf(errNotDefined, lxc.name)
	}

	if !lxc.Running() {
		return fmt.Errorf(errNotRunning, lxc.name)
	}
	return nil
}

func (lxc *Container) ensureDefinedButNotRunning() error {
	if !lxc.Defined() {
		return fmt.Errorf(errNotDefined, lxc.name)
	}

	if lxc.Running() {
		return fmt.Errorf(errAlreadyRunning, lxc.name)
	}
	return nil
}

// Name returns the name of the container
func (lxc *Container) Name() string {
	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	return C.GoString(lxc.container.name)
}

// Defined returns true if the container is already defined
func (lxc *Container) Defined() bool {
	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	return bool(C.lxc_container_defined(lxc.container))
}

// Running returns true if the container is already running
func (lxc *Container) Running() bool {
	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	return bool(C.lxc_container_running(lxc.container))
}

// MayControl returns true if the caller may control the container
func (lxc *Container) MayControl() bool {
	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	return bool(C.lxc_container_may_control(lxc.container))
}

// CreateSnapshot creates a new snapshot
func (lxc *Container) CreateSnapshot() (*Snapshot, error) {
	if err := lxc.ensureDefinedButNotRunning(); err != nil {
		return nil, err
	}

	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	ret := int(C.lxc_container_snapshot(lxc.container))
	if ret < 0 {
		return nil, fmt.Errorf(errCreateSnapshotFailed, lxc.name)
	}
	return &Snapshot{Name: fmt.Sprintf("snap%d", ret)}, nil
}

// RestoreSnapshot creates a new container based on a snapshot
func (lxc *Container) RestoreSnapshot(snapshot Snapshot, name string) error {
	if !lxc.Defined() {
		return fmt.Errorf(errNotDefined, lxc.name)
	}

	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	csnapname := C.CString(snapshot.Name)
	defer C.free(unsafe.Pointer(csnapname))

	if !bool(C.lxc_container_snapshot_restore(lxc.container, csnapname, cname)) {
		return fmt.Errorf(errRestoreSnapshotFailed, lxc.name)
	}
	return nil
}

// DestroySnapshot destroys the specified snapshot
func (lxc *Container) DestroySnapshot(snapshot Snapshot) error {
	if !lxc.Defined() {
		return fmt.Errorf(errNotDefined, lxc.name)
	}

	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	csnapname := C.CString(snapshot.Name)
	defer C.free(unsafe.Pointer(csnapname))

	if !bool(C.lxc_container_snapshot_destroy(lxc.container, csnapname)) {
		return fmt.Errorf(errDestroySnapshotFailed, lxc.name)
	}
	return nil
}

// Snapshots returns the list of container snapshots
func (lxc *Container) Snapshots() ([]Snapshot, error) {
	if !lxc.Defined() {
		return nil, fmt.Errorf(errNotDefined, lxc.name)
	}

	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	var csnapshots *C.struct_lxc_snapshot
	var snapshots []Snapshot

	size := int(C.lxc_container_snapshot_list(lxc.container, &csnapshots))
	defer freeSnapshots(csnapshots, size)

	if size < 1 {
		return nil, fmt.Errorf("%s has no snapshots", lxc.name)
	}

	p := uintptr(unsafe.Pointer(csnapshots))
	for i := 0; i < size; i++ {
		z := (*C.struct_lxc_snapshot)(unsafe.Pointer(p))
		s := &Snapshot{Name: C.GoString(z.name), Timestamp: C.GoString(z.timestamp), CommentPath: C.GoString(z.comment_pathname), Path: C.GoString(z.lxcpath)}
		snapshots = append(snapshots, *s)
		p += unsafe.Sizeof(*csnapshots)
	}
	return snapshots, nil
}

// State returns the state of the container
func (lxc *Container) State() State {
	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	return stateMap[C.GoString(C.lxc_container_state(lxc.container))]
}

// InitPID returns the process ID of the container's init process seen from outside the container
func (lxc *Container) InitPID() int {
	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	return int(C.lxc_container_init_pid(lxc.container))
}

// Daemonize returns whether the container wished to be daemonized
func (lxc *Container) Daemonize() bool {
	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	return bool(lxc.container.daemonize)
}

// WantDaemonize sets the daemonize flag for the container
func (lxc *Container) WantDaemonize(state bool) error {
	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	if !bool(C.lxc_container_want_daemonize(lxc.container, C.bool(state))) {
		return fmt.Errorf(errDaemonizeFailed, lxc.name)
	}
	return nil
}

// WantCloseAllFds sets the close_all_fds flag for the container
func (lxc *Container) WantCloseAllFds(state bool) error {
	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	if !bool(C.lxc_container_want_close_all_fds(lxc.container, C.bool(state))) {
		return fmt.Errorf(errCloseAllFdsFailed, lxc.name)
	}
	return nil
}

// SetVerbosity sets the verbosity level of some API calls
func (lxc *Container) SetVerbosity(verbosity Verbosity) {
	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	lxc.verbosity = verbosity
}

// Freeze freezes the running container
func (lxc *Container) Freeze() error {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	if lxc.State() == FROZEN {
		return fmt.Errorf(errAlreadyFrozen, lxc.name)
	}

	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	if !bool(C.lxc_container_freeze(lxc.container)) {
		return fmt.Errorf(errFreezeFailed, lxc.name)
	}

	return nil
}

// Unfreeze thaws the frozen container
func (lxc *Container) Unfreeze() error {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	if lxc.State() != FROZEN {
		return fmt.Errorf(errNotFrozen, lxc.name)
	}

	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	if !bool(C.lxc_container_unfreeze(lxc.container)) {
		return fmt.Errorf(errUnfreezeFailed, lxc.name)
	}

	return nil
}

// CreateUsing creates the container using given template and arguments with specified backend
func (lxc *Container) CreateUsing(template string, backend BackendStore, args ...string) error {
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
	if lxc.Defined() {
		return fmt.Errorf(errAlreadyDefined, lxc.name)
	}

	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	ctemplate := C.CString(template)
	defer C.free(unsafe.Pointer(ctemplate))

	cbackend := C.CString(backend.String())
	defer C.free(unsafe.Pointer(cbackend))

	ret := false
	if args != nil {
		cargs := makeNullTerminatedArgs(args)
		defer freeNullTerminatedArgs(cargs, len(args))

		ret = bool(C.lxc_container_create(lxc.container, ctemplate, cbackend, C.int(lxc.verbosity), cargs))
	} else {
		ret = bool(C.lxc_container_create(lxc.container, ctemplate, cbackend, C.int(lxc.verbosity), nil))
	}

	if !ret {
		return fmt.Errorf(errCreateFailed, lxc.name)
	}
	return nil
}

// Create creates the container using given template and arguments with Best backend
func (lxc *Container) Create(template string, args ...string) error {
	return lxc.CreateUsing(template, Best, args...)
}

// CreateAsUser creates the container using given template and arguments with Best backend as an unprivileged user
func (lxc *Container) CreateAsUser(distro string, release string, arch string, args ...string) error {
	// required parameters
	nargs := []string{"-d", distro, "-r", release, "-a", arch}
	// optional arguments
	nargs = append(nargs, args...)

	return lxc.CreateUsing("download", Best, nargs...)
}

// Start starts the container
func (lxc *Container) Start() error {
	if err := lxc.ensureDefinedButNotRunning(); err != nil {
		return err
	}

	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	if !bool(C.lxc_container_start(lxc.container, 0, nil)) {
		return fmt.Errorf(errStartFailed, lxc.name)
	}
	return nil
}

// Execute executes the given command in a temporary container
func (lxc *Container) Execute(args ...string) ([]byte, error) {
	if lxc.Defined() {
		return nil, fmt.Errorf(errAlreadyDefined, lxc.name)
	}

	cargs := []string{"lxc-execute", "-n", lxc.Name(), "-P", lxc.ConfigPath(), "--"}
	cargs = append(cargs, args...)

	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	/*
	 * FIXME: Go runtime and src/lxc/start.c signal_handler are not playing nice together so use lxc-execute for now
	 */
	output, err := exec.Command(cargs[0], cargs[1:]...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf(errExecuteFailed, lxc.name)
	}

	return output, nil
	/*
		cargs := makeNullTerminatedArgs(args)
		defer freeNullTerminatedArgs(cargs, len(args))

		if !bool(C.lxc_container_start(lxc.container, 1, cargs)) {
			return fmt.Errorf(errExecuteFailed, lxc.name)
		}
		return nil
	*/
}

// Stop stops the container
func (lxc *Container) Stop() error {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	if !bool(C.lxc_container_stop(lxc.container)) {
		return fmt.Errorf(errStopFailed, lxc.name)
	}
	return nil
}

// Reboot reboots the container
func (lxc *Container) Reboot() error {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	if !bool(C.lxc_container_reboot(lxc.container)) {
		return fmt.Errorf(errRebootFailed, lxc.name)
	}
	return nil
}

// Shutdown shutdowns the container
func (lxc *Container) Shutdown(timeout int) error {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	if !bool(C.lxc_container_shutdown(lxc.container, C.int(timeout))) {
		return fmt.Errorf(errShutdownFailed, lxc.name)
	}
	return nil
}

// Destroy destroys the container
func (lxc *Container) Destroy() error {
	if err := lxc.ensureDefinedButNotRunning(); err != nil {
		return err
	}

	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	if !bool(C.lxc_container_destroy(lxc.container)) {
		return fmt.Errorf(errDestroyFailed, lxc.name)
	}
	return nil
}

// CloneUsing clones the container using specified backend
func (lxc *Container) CloneUsing(name string, backend BackendStore, flags CloneFlags) error {
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
	if err := lxc.ensureDefinedButNotRunning(); err != nil {
		return err
	}

	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cbackend := C.CString(backend.String())
	defer C.free(unsafe.Pointer(cbackend))

	if !bool(C.lxc_container_clone(lxc.container, cname, C.int(flags), cbackend)) {
		return fmt.Errorf(errCloneFailed, lxc.name)
	}
	return nil
}

// Clone clones the container using the Directory backendstore
func (lxc *Container) Clone(name string) error {
	return lxc.CloneUsing(name, Directory, 0)
}

// Rename renames the container
func (lxc *Container) Rename(name string) error {
	if err := lxc.ensureDefinedButNotRunning(); err != nil {
		return err
	}

	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	if !bool(C.lxc_container_rename(lxc.container, cname)) {
		return fmt.Errorf(errRenameFailed, lxc.name)
	}
	return nil
}

// Wait waits for container to reach a given state or timeouts
func (lxc *Container) Wait(state State, timeout int) bool {
	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	cstate := C.CString(state.String())
	defer C.free(unsafe.Pointer(cstate))

	return bool(C.lxc_container_wait(lxc.container, cstate, C.int(timeout)))
}

// ConfigFileName returns the container's configuration file's name
func (lxc *Container) ConfigFileName() string {
	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	// allocated in lxc.c
	configFileName := C.lxc_container_config_file_name(lxc.container)
	defer C.free(unsafe.Pointer(configFileName))

	return C.GoString(configFileName)
}

// ConfigItem returns the value of the given config item
func (lxc *Container) ConfigItem(key string) []string {
	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	// allocated in lxc.c
	configItem := C.lxc_container_get_config_item(lxc.container, ckey)
	defer C.free(unsafe.Pointer(configItem))

	ret := strings.TrimSpace(C.GoString(configItem))
	return strings.Split(ret, "\n")
}

// SetConfigItem sets the value of the given config item
func (lxc *Container) SetConfigItem(key string, value string) error {
	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	if !bool(C.lxc_container_set_config_item(lxc.container, ckey, cvalue)) {
		return fmt.Errorf(errSettingConfigItemFailed, lxc.name, key, value)
	}
	return nil
}

// CgroupItem returns the value of the given cgroup subsystem value
func (lxc *Container) CgroupItem(key string) []string {
	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	// allocated in lxc.c
	cgroupItem := C.lxc_container_get_cgroup_item(lxc.container, ckey)
	defer C.free(unsafe.Pointer(cgroupItem))

	ret := strings.TrimSpace(C.GoString(cgroupItem))
	return strings.Split(ret, "\n")
}

// SetCgroupItem sets the value of given cgroup subsystem value
func (lxc *Container) SetCgroupItem(key string, value string) error {
	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	if !bool(C.lxc_container_set_cgroup_item(lxc.container, ckey, cvalue)) {
		return fmt.Errorf(errSettingCgroupItemFailed, lxc.name, key, value)
	}
	return nil
}

// ClearConfigItem clears the value of given config item
func (lxc *Container) ClearConfigItem(key string) error {
	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	if !bool(C.lxc_container_clear_config_item(lxc.container, ckey)) {
		return fmt.Errorf(errClearingCgroupItemFailed, lxc.name, key)
	}
	return nil
}

// ConfigKeys returns the names of the config items
func (lxc *Container) ConfigKeys(key ...string) []string {
	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	var keys *_Ctype_char

	if key != nil && len(key) == 1 {
		ckey := C.CString(key[0])
		defer C.free(unsafe.Pointer(ckey))

		// allocated in lxc.c
		keys = C.lxc_container_get_keys(lxc.container, ckey)
		defer C.free(unsafe.Pointer(keys))
	} else {
		// allocated in lxc.c
		keys = C.lxc_container_get_keys(lxc.container, nil)
		defer C.free(unsafe.Pointer(keys))
	}
	ret := strings.TrimSpace(C.GoString(keys))
	return strings.Split(ret, "\n")
}

// LoadConfigFile loads the configuration file from given path
func (lxc *Container) LoadConfigFile(path string) error {
	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	if !bool(C.lxc_container_load_config(lxc.container, cpath)) {
		return fmt.Errorf(errLoadConfigFailed, lxc.name, path)
	}
	return nil
}

// SaveConfigFile saves the configuration file to given path
func (lxc *Container) SaveConfigFile(path string) error {
	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	if !bool(C.lxc_container_save_config(lxc.container, cpath)) {
		return fmt.Errorf(errSaveConfigFailed, lxc.name, path)
	}
	return nil
}

// ConfigPath returns the configuration file's path
func (lxc *Container) ConfigPath() string {
	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	return C.GoString(C.lxc_container_get_config_path(lxc.container))
}

// SetConfigPath sets the configuration file's path
func (lxc *Container) SetConfigPath(path string) error {
	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	if !bool(C.lxc_container_set_config_path(lxc.container, cpath)) {
		return fmt.Errorf(errSettingConfigPathFailed, lxc.name, path)
	}
	return nil
}

// MemoryUsage returns memory usage of the container in bytes
func (lxc *Container) MemoryUsage() (ByteSize, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return -1, err
	}

	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	memUsed, err := strconv.ParseFloat(lxc.CgroupItem("memory.usage_in_bytes")[0], 64)
	if err != nil {
		return -1, fmt.Errorf(errMemLimit)
	}
	return ByteSize(memUsed), err
}

// KernelMemoryUsage returns current kernel memory allocation of the container in bytes
func (lxc *Container) KernelMemoryUsage() (ByteSize, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return -1, err
	}

	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	kmemUsed, err := strconv.ParseFloat(lxc.CgroupItem("memory.kmem.usage_in_bytes")[0], 64)
	if err != nil {
		return -1, fmt.Errorf(errKMemLimit)
	}
	return ByteSize(kmemUsed), err
}

// SwapUsage returns swap usage of the container in bytes
func (lxc *Container) SwapUsage() (ByteSize, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return -1, err
	}

	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	swapUsed, err := strconv.ParseFloat(lxc.CgroupItem("memory.memsw.usage_in_bytes")[0], 64)
	if err != nil {
		return -1, fmt.Errorf(errSwapLimit)
	}
	return ByteSize(swapUsed), err
}

// BlkioUsage returns number of bytes transferred to/from the disk by the container
func (lxc *Container) BlkioUsage() (ByteSize, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return -1, err
	}

	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	for _, v := range lxc.CgroupItem("blkio.throttle.io_service_bytes") {
		b := strings.Split(v, " ")
		if b[0] == "Total" {
			blkioUsed, err := strconv.ParseFloat(b[1], 64)
			if err != nil {
				return -1, err
			}
			return ByteSize(blkioUsed), err
		}
	}
	return -1, fmt.Errorf(errBlkioUsage, lxc.name)
}

// MemoryLimit returns memory limit of the container in bytes
func (lxc *Container) MemoryLimit() (ByteSize, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return -1, err
	}

	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	memLimit, err := strconv.ParseFloat(lxc.CgroupItem("memory.limit_in_bytes")[0], 64)
	if err != nil {
		return -1, fmt.Errorf(errMemLimit)
	}
	return ByteSize(memLimit), err
}

// SetMemoryLimit sets memory limit of the container in bytes
func (lxc *Container) SetMemoryLimit(limit ByteSize) error {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	if err := lxc.SetCgroupItem("memory.limit_in_bytes", fmt.Sprintf("%.f", limit)); err != nil {
		return fmt.Errorf(errSettingMemoryLimitFailed, lxc.name)
	}
	return nil
}

// KernelMemoryLimit returns kernel memory limit of the container in bytes
func (lxc *Container) KernelMemoryLimit() (ByteSize, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return -1, err
	}

	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	kmemLimit, err := strconv.ParseFloat(lxc.CgroupItem("memory.kmem.limit_in_bytes")[0], 64)
	if err != nil {
		return -1, fmt.Errorf(errKMemLimit)
	}
	return ByteSize(kmemLimit), err
}

// SetKernelMemoryLimit sets kernel memory limit of the container in bytes
func (lxc *Container) SetKernelMemoryLimit(limit ByteSize) error {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	if err := lxc.SetCgroupItem("memory.kmem.limit_in_bytes", fmt.Sprintf("%.f", limit)); err != nil {
		return fmt.Errorf(errSettingKMemoryLimitFailed, lxc.name)
	}
	return nil
}

// SwapLimit returns the swap limit of the container in bytes
func (lxc *Container) SwapLimit() (ByteSize, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return -1, err
	}

	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	swapLimit, err := strconv.ParseFloat(lxc.CgroupItem("memory.memsw.limit_in_bytes")[0], 64)
	if err != nil {
		return -1, fmt.Errorf(errSwapLimit)
	}
	return ByteSize(swapLimit), err
}

// SetSwapLimit sets memory limit of the container in bytes
func (lxc *Container) SetSwapLimit(limit ByteSize) error {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	if err := lxc.SetCgroupItem("memory.memsw.limit_in_bytes", fmt.Sprintf("%.f", limit)); err != nil {
		return fmt.Errorf(errSettingSwapLimitFailed, lxc.name)
	}
	return nil
}

// CPUTime returns the total CPU time (in nanoseconds) consumed by all tasks in this cgroup (including tasks lower in the hierarchy).
func (lxc *Container) CPUTime() (time.Duration, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return -1, err
	}

	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	cpuUsage, err := strconv.ParseInt(lxc.CgroupItem("cpuacct.usage")[0], 10, 64)
	if err != nil {
		return -1, err
	}
	return time.Duration(cpuUsage), err
}

// CPUTimePerCPU returns the CPU time (in nanoseconds) consumed on each CPU by all tasks in this cgroup (including tasks lower in the hierarchy).
func (lxc *Container) CPUTimePerCPU() ([]time.Duration, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return nil, err
	}

	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	var cpuTimes []time.Duration

	for _, v := range strings.Split(lxc.CgroupItem("cpuacct.usage_percpu")[0], " ") {
		cpuUsage, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		cpuTimes = append(cpuTimes, time.Duration(cpuUsage))
	}
	return cpuTimes, nil
}

// CPUStats returns the number of CPU cycles (in the units defined by USER_HZ on the system) consumed by tasks in this cgroup and its children in both user mode and system (kernel) mode.
func (lxc *Container) CPUStats() ([]int64, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return nil, err
	}

	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	cpuStat := lxc.CgroupItem("cpuacct.stat")
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
func (lxc *Container) ConsoleGetFD(ttynum int) (int, error) {
	// FIXME: Make idiomatic
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return -1, err
	}

	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	ret := int(C.lxc_container_console_getfd(lxc.container, C.int(ttynum)))
	if ret < 0 {
		return -1, fmt.Errorf(errAttachFailed, lxc.name)
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
func (lxc *Container) Console(ttynum, stdinfd, stdoutfd, stderrfd, escape int) error {
	// FIXME: Make idiomatic
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	if !bool(C.lxc_container_console(lxc.container, C.int(ttynum), C.int(stdinfd), C.int(stdoutfd), C.int(stderrfd), C.int(escape))) {
		return fmt.Errorf(errAttachFailed, lxc.name)
	}
	return nil
}

// AttachShell runs a shell inside the container
func (lxc *Container) AttachShell() error {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	if int(C.lxc_container_attach(lxc.container, false)) < 0 {
		return fmt.Errorf(errAttachFailed, lxc.name)
	}
	return nil
}

// AttachShellWithClearEnvironment runs a shell inside the container and clears all environment variables before attaching
func (lxc *Container) AttachShellWithClearEnvironment() error {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	if int(C.lxc_container_attach(lxc.container, true)) < 0 {
		return fmt.Errorf(errAttachFailed, lxc.name)
	}
	return nil
}

// RunCommand runs the user specified command inside the container and waits it to exit
func (lxc *Container) RunCommand(args ...string) error {
	if args == nil {
		return fmt.Errorf(errInsufficientNumberOfArguments)
	}

	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	cargs := makeNullTerminatedArgs(args)
	defer freeNullTerminatedArgs(cargs, len(args))

	if int(C.lxc_container_attach_run_wait(lxc.container, false, cargs)) < 0 {
		return fmt.Errorf(errAttachFailed, lxc.name)
	}
	return nil
}

// RunCommandWithClearEnvironment runs the user specified command inside the container and waits it to exit. It also clears all environment variables before attaching.
func (lxc *Container) RunCommandWithClearEnvironment(args ...string) error {
	if args == nil {
		return fmt.Errorf(errInsufficientNumberOfArguments)
	}

	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	cargs := makeNullTerminatedArgs(args)
	defer freeNullTerminatedArgs(cargs, len(args))

	if int(C.lxc_container_attach_run_wait(lxc.container, true, cargs)) < 0 {
		return fmt.Errorf(errAttachFailed, lxc.name)
	}
	return nil
}

// Interfaces returns the name of the network interfaces from the container
func (lxc *Container) Interfaces() ([]string, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return nil, err
	}

	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	result := C.lxc_container_get_interfaces(lxc.container)
	if result == nil {
		return nil, fmt.Errorf(errInterfaces, lxc.name)
	}
	return convertArgs(result), nil
}

// IPAddress returns the IP address of the given network interface from the container
func (lxc *Container) IPAddress(interfaceName string) ([]string, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return nil, err
	}

	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	cinterface := C.CString(interfaceName)
	defer C.free(unsafe.Pointer(cinterface))

	result := C.lxc_container_get_ips(lxc.container, cinterface, nil, 0)
	if result == nil {
		return nil, fmt.Errorf(errIPAddress, interfaceName, lxc.name)
	}
	return convertArgs(result), nil
}

// IPAddresses returns all IP addresses from the container
func (lxc *Container) IPAddresses() ([]string, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return nil, err
	}

	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	result := C.lxc_container_get_ips(lxc.container, nil, nil, 0)
	if result == nil {
		return nil, fmt.Errorf(errIPAddresses, lxc.name)
	}
	return convertArgs(result), nil

}

// IPv4Addresses returns all IPv4 addresses from the container
func (lxc *Container) IPv4Addresses() ([]string, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return nil, err
	}

	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	cfamily := C.CString("inet")
	defer C.free(unsafe.Pointer(cfamily))

	result := C.lxc_container_get_ips(lxc.container, nil, cfamily, 0)
	if result == nil {
		return nil, fmt.Errorf(errIPv4Addresses, lxc.name)
	}
	return convertArgs(result), nil
}

// IPv6Addresses returns all IPv6 addresses from the container
func (lxc *Container) IPv6Addresses() ([]string, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return nil, err
	}

	lxc.mu.RLock()
	defer lxc.mu.RUnlock()

	cfamily := C.CString("inet6")
	defer C.free(unsafe.Pointer(cfamily))

	result := C.lxc_container_get_ips(lxc.container, nil, cfamily, 0)
	if result == nil {
		return nil, fmt.Errorf(errIPv6Addresses, lxc.name)
	}
	return convertArgs(result), nil
}

// LogFile returns the name of the logfile
func (lxc *Container) LogFile() string {
	return lxc.ConfigItem("lxc.logfile")[0]
}

// SetLogFile sets the logfile to given filename
func (lxc *Container) SetLogFile(filename string) error {
	if err := lxc.SetConfigItem("lxc.logfile", filename); err != nil {
		return err
	}
	return nil
}

// LogLevel returns the name of the logfile
func (lxc *Container) LogLevel() LogLevel {
	return logLevelMap[lxc.ConfigItem("lxc.loglevel")[0]]
}

// SetLogLevel sets the logfile to given filename
func (lxc *Container) SetLogLevel(level LogLevel) error {
	if err := lxc.SetConfigItem("lxc.loglevel", level.String()); err != nil {
		return err
	}
	return nil
}

// AddDeviceNode adds specified device to the container
func (lxc *Container) AddDeviceNode(source string, destination ...string) error {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	csource := C.CString(source)
	defer C.free(unsafe.Pointer(csource))

	if destination != nil && len(destination) == 1 {
		cdestination := C.CString(destination[0])
		defer C.free(unsafe.Pointer(cdestination))

		if !bool(C.lxc_container_add_device_node(lxc.container, csource, cdestination)) {
			return fmt.Errorf(errAddDeviceNodeFailed, source, lxc.name)
		}
		return nil
	}

	if !bool(C.lxc_container_add_device_node(lxc.container, csource, nil)) {
		return fmt.Errorf(errAddDeviceNodeFailed, source, lxc.name)
	}
	return nil

}

// RemoveDeviceNode removes the specified device from the container
func (lxc *Container) RemoveDeviceNode(source string, destination ...string) error {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	lxc.mu.Lock()
	defer lxc.mu.Unlock()

	csource := C.CString(source)
	defer C.free(unsafe.Pointer(csource))

	if destination != nil && len(destination) == 1 {
		cdestination := C.CString(destination[0])
		defer C.free(unsafe.Pointer(cdestination))

		if !bool(C.lxc_container_remove_device_node(lxc.container, csource, cdestination)) {
			return fmt.Errorf(errRemoveDeviceNodeFailed, source, lxc.name)
		}
		return nil
	}

	if !bool(C.lxc_container_remove_device_node(lxc.container, csource, nil)) {
		return fmt.Errorf(errRemoveDeviceNodeFailed, source, lxc.name)
	}
	return nil
}
