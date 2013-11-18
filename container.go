// Copyright © 2013, S.Çağlar Onur
// Use of this source code is governed by a LGPLv2.1
// license that can be found in the LICENSE file.
//
// Authors:
// S.Çağlar Onur <caglar@10ur.org>

// +build linux

package lxc

// #include <lxc/lxc.h>
// #include <lxc/lxccontainer.h>
// #include "lxc.h"
import "C"

import (
	"fmt"
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
	sync.RWMutex
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
		return fmt.Errorf(errNotDefined, C.GoString(lxc.container.name))
	}

	if !lxc.Running() {
		return fmt.Errorf(errNotRunning, C.GoString(lxc.container.name))
	}
	return nil
}

func (lxc *Container) ensureDefinedButNotRunning() error {
	if !lxc.Defined() {
		return fmt.Errorf(errNotDefined, C.GoString(lxc.container.name))
	}

	if lxc.Running() {
		return fmt.Errorf(errAlreadyRunning, C.GoString(lxc.container.name))
	}
	return nil
}

// Name returns container's name
func (lxc *Container) Name() string {
	lxc.RLock()
	defer lxc.RUnlock()

	return C.GoString(lxc.container.name)
}

// Defined returns whether the container is already defined or not
func (lxc *Container) Defined() bool {
	lxc.RLock()
	defer lxc.RUnlock()

	return bool(C.lxc_container_defined(lxc.container))
}

// Running returns whether the container is already running or not
func (lxc *Container) Running() bool {
	lxc.RLock()
	defer lxc.RUnlock()

	return bool(C.lxc_container_running(lxc.container))
}

// MayControl returns whether the container is already running or not
func (lxc *Container) MayControl() bool {
	lxc.RLock()
	defer lxc.RUnlock()

	return bool(C.lxc_container_may_control(lxc.container))
}

// CreateSnapshot creates a new snapshot
func (lxc *Container) CreateSnapshot() error {
	if err := lxc.ensureDefinedButNotRunning(); err != nil {
		return err
	}

	lxc.Lock()
	defer lxc.Unlock()

	// FIXME: LXC C API returns the number of snapshots, should we return it as well?
	if int(C.lxc_container_snapshot(lxc.container)) < 0 {
		return fmt.Errorf(errCreateSnapshotFailed, C.GoString(lxc.container.name))
	}
	return nil
}

// RestoreSnapshot creates a new snapshot
func (lxc *Container) RestoreSnapshot(snapshot Snapshot, name string) error {
	if !lxc.Defined() {
		return fmt.Errorf(errNotDefined, C.GoString(lxc.container.name))
	}

	lxc.Lock()
	defer lxc.Unlock()

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	csnapname := C.CString(snapshot.Name)
	defer C.free(unsafe.Pointer(csnapname))

	if !bool(C.lxc_container_snapshot_restore(lxc.container, csnapname, cname)) {
		return fmt.Errorf(errRestoreSnapshotFailed, C.GoString(lxc.container.name))
	}
	return nil
}

// DestroySnapshot destroys the snapshot
func (lxc *Container) DestroySnapshot(snapshot Snapshot) error {
	if !lxc.Defined() {
		return fmt.Errorf(errNotDefined, C.GoString(lxc.container.name))
	}

	lxc.Lock()
	defer lxc.Unlock()

	csnapname := C.CString(snapshot.Name)
	defer C.free(unsafe.Pointer(csnapname))

	if !bool(C.lxc_container_snapshot_destroy(lxc.container, csnapname)) {
		return fmt.Errorf(errDestroySnapshotFailed, C.GoString(lxc.container.name))
	}
	return nil
}

// Snapshots lists the snapshot of given container
func (lxc *Container) Snapshots() ([]Snapshot, error) {
	if !lxc.Defined() {
		return nil, fmt.Errorf(errNotDefined, C.GoString(lxc.container.name))
	}

	lxc.Lock()
	defer lxc.Unlock()

	var csnapshots *C.struct_lxc_snapshot
	var snapshots []Snapshot

	size := int(C.lxc_container_snapshot_list(lxc.container, &csnapshots))
	defer freeSnapshots(csnapshots, size)

	if size < 1 {
		return nil, fmt.Errorf("%s has no snapshots", C.GoString(lxc.container.name))
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

// State returns the container's state
func (lxc *Container) State() State {
	lxc.RLock()
	defer lxc.RUnlock()

	return stateMap[C.GoString(C.lxc_container_state(lxc.container))]
}

// InitPID returns the container's PID
func (lxc *Container) InitPID() int {
	lxc.RLock()
	defer lxc.RUnlock()

	return int(C.lxc_container_init_pid(lxc.container))
}

// Daemonize returns whether the daemonize flag is set
func (lxc *Container) Daemonize() bool {
	lxc.RLock()
	defer lxc.RUnlock()

	return bool(lxc.container.daemonize != 0)
}

// SetDaemonize sets the daemonize flag
func (lxc *Container) SetDaemonize() error {
	lxc.Lock()
	defer lxc.Unlock()

	C.lxc_container_want_daemonize(lxc.container)
	if bool(lxc.container.daemonize == 0) {
		return fmt.Errorf(errDaemonizeFailed, C.GoString(lxc.container.name))
	}
	return nil
}

// SetCloseAllFds sets the close_all_fds flag for the container
func (lxc *Container) SetCloseAllFds() error {
	lxc.Lock()
	defer lxc.Unlock()

	if !bool(C.lxc_container_want_close_all_fds(lxc.container)) {
		return fmt.Errorf(errCloseAllFdsFailed, C.GoString(lxc.container.name))
	}
	return nil
}

// SetVerbosity sets the verbosity level of some API calls
func (lxc *Container) SetVerbosity(verbosity Verbosity) {
	lxc.Lock()
	defer lxc.Unlock()

	lxc.verbosity = verbosity
}

// Freeze freezes the running container
func (lxc *Container) Freeze() error {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	if lxc.State() == FROZEN {
		return fmt.Errorf(errAlreadyFrozen, C.GoString(lxc.container.name))
	}

	lxc.Lock()
	defer lxc.Unlock()

	if !bool(C.lxc_container_freeze(lxc.container)) {
		return fmt.Errorf(errFreezeFailed, C.GoString(lxc.container.name))
	}

	return nil
}

// Unfreeze unfreezes the frozen container
func (lxc *Container) Unfreeze() error {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	if lxc.State() != FROZEN {
		return fmt.Errorf(errNotFrozen, C.GoString(lxc.container.name))
	}

	lxc.Lock()
	defer lxc.Unlock()

	if !bool(C.lxc_container_unfreeze(lxc.container)) {
		return fmt.Errorf(errUnfreezeFailed, C.GoString(lxc.container.name))
	}

	return nil
}

// Create creates the container using given template and arguments
func (lxc *Container) Create(template string, args ...string) error {
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
		return fmt.Errorf(errAlreadyDefined, C.GoString(lxc.container.name))
	}

	lxc.Lock()
	defer lxc.Unlock()

	ctemplate := C.CString(template)
	defer C.free(unsafe.Pointer(ctemplate))

	cbdevtype := C.CString("dir")
	defer C.free(unsafe.Pointer(cbdevtype))

	ret := false
	if args != nil {
		cargs := makeNullTerminatedArgs(args)
		defer freeNullTerminatedArgs(cargs, len(args))

		ret = bool(C.lxc_container_create(lxc.container, ctemplate, cbdevtype, C.int(lxc.verbosity), cargs))
	} else {
		ret = bool(C.lxc_container_create(lxc.container, ctemplate, cbdevtype, C.int(lxc.verbosity), nil))
	}

	if !ret {
		return fmt.Errorf(errCreateFailed, C.GoString(lxc.container.name))
	}
	return nil
}

// Start starts the container
func (lxc *Container) Start() error {
	if err := lxc.ensureDefinedButNotRunning(); err != nil {
		return err
	}

	lxc.Lock()
	defer lxc.Unlock()

	if !bool(C.lxc_container_start(lxc.container, C.int(0), nil)) {
		return fmt.Errorf(errStartFailed, C.GoString(lxc.container.name))
	}
	return nil
}

// Execute executes the given argument in a temporary container
func (lxc *Container) Execute(args ...string) error {
	// FIXME: disable for now
	return fmt.Errorf("NOT SUPPORTED")
/*
	if lxc.Defined() || lxc.Running() {
		return fmt.Errorf(errAlreadyDefined, C.GoString(lxc.container.name))
	}

	lxc.Lock()
	defer lxc.Unlock()

	cargs := makeNullTerminatedArgs(args)
	defer freeNullTerminatedArgs(cargs, len(args))

	if !bool(C.lxc_container_start(lxc.container, C.int(1), cargs)) {
		return fmt.Errorf(errExecuteFailed, C.GoString(lxc.container.name))
	}
	return nil
*/
}

// Stop stops the container
func (lxc *Container) Stop() error {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	lxc.Lock()
	defer lxc.Unlock()

	if !bool(C.lxc_container_stop(lxc.container)) {
		return fmt.Errorf(errStopFailed, C.GoString(lxc.container.name))
	}
	return nil
}

// Reboot reboots the container
func (lxc *Container) Reboot() error {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	lxc.Lock()
	defer lxc.Unlock()

	if !bool(C.lxc_container_reboot(lxc.container)) {
		return fmt.Errorf(errRebootFailed, C.GoString(lxc.container.name))
	}
	return nil
}

// Shutdown shutdowns the container
func (lxc *Container) Shutdown(timeout int) error {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	lxc.Lock()
	defer lxc.Unlock()

	if !bool(C.lxc_container_shutdown(lxc.container, C.int(timeout))) {
		return fmt.Errorf(errShutdownFailed, C.GoString(lxc.container.name))
	}
	return nil
}

// Destroy destroys the container
func (lxc *Container) Destroy() error {
	if err := lxc.ensureDefinedButNotRunning(); err != nil {
		return err
	}

	lxc.Lock()
	defer lxc.Unlock()

	if !bool(C.lxc_container_destroy(lxc.container)) {
		return fmt.Errorf(errDestroyFailed, C.GoString(lxc.container.name))
	}
	return nil
}

// Clone clones the container
func (lxc *Container) Clone(name string, flags int, backend BackendStore) error {
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
	// flags: LXC_CLONE_SNAPSHOT || LXC_CLONE_KEEPNAME || LXC_CLONE_KEEPMACADDR || LXC_CLONE_COPYHOOKS
	//
	// newsize: for blockdev-backed backingstores
	//
	// hookargs: additional arguments to pass to the clone hook script
	if err := lxc.ensureDefinedButNotRunning(); err != nil {
		return err
	}

	lxc.Lock()
	defer lxc.Unlock()

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	if !bool(C.lxc_container_clone(lxc.container, cname, C.int(flags), C.CString(backend.String()))) {
		return fmt.Errorf(errCloneFailed, C.GoString(lxc.container.name))
	}
	return nil
}

// CloneToDirectory clones the container using Directory backendstore
func (lxc *Container) CloneToDirectory(name string) error {
	return lxc.Clone(name, 0, Directory)
}

// CloneToBtrFS clones the container using BtrFS backendstore
func (lxc *Container) CloneToBtrFS(name string) error {
	return lxc.Clone(name, 0, BtrFS)
}

// CloneToOverlayFS clones the container using OverlayFS backendstore
func (lxc *Container) CloneToOverlayFS(name string) error {
	return lxc.Clone(name, CloneSnapshot, OverlayFS)
}

// Wait waits till the container changes its state or timeouts
func (lxc *Container) Wait(state State, timeout int) bool {
	lxc.Lock()
	defer lxc.Unlock()

	cstate := C.CString(state.String())
	defer C.free(unsafe.Pointer(cstate))

	return bool(C.lxc_container_wait(lxc.container, cstate, C.int(timeout)))
}

// ConfigFileName returns the container's configuration file's name
func (lxc *Container) ConfigFileName() string {
	lxc.RLock()
	defer lxc.RUnlock()

	// allocated in lxc.c
	configFileName := C.lxc_container_config_file_name(lxc.container)
	defer C.free(unsafe.Pointer(configFileName))

	return C.GoString(configFileName)
}

// ConfigItem returns the value of the given key
func (lxc *Container) ConfigItem(key string) []string {
	lxc.RLock()
	defer lxc.RUnlock()

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	// allocated in lxc.c
	configItem := C.lxc_container_get_config_item(lxc.container, ckey)
	defer C.free(unsafe.Pointer(configItem))

	ret := strings.TrimSpace(C.GoString(configItem))
	return strings.Split(ret, "\n")
}

// SetConfigItem sets the value of given key
func (lxc *Container) SetConfigItem(key string, value string) error {
	lxc.Lock()
	defer lxc.Unlock()

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	if !bool(C.lxc_container_set_config_item(lxc.container, ckey, cvalue)) {
		return fmt.Errorf(errSettingConfigItemFailed, C.GoString(lxc.container.name), key, value)
	}
	return nil
}

// CgroupItem returns the value of the given key
func (lxc *Container) CgroupItem(key string) []string {
	lxc.RLock()
	defer lxc.RUnlock()

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	// allocated in lxc.c
	cgroupItem := C.lxc_container_get_cgroup_item(lxc.container, ckey)
	defer C.free(unsafe.Pointer(cgroupItem))

	ret := strings.TrimSpace(C.GoString(cgroupItem))
	return strings.Split(ret, "\n")
}

// SetCgroupItem sets the value of given key
func (lxc *Container) SetCgroupItem(key string, value string) error {
	lxc.Lock()
	defer lxc.Unlock()

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	if !bool(C.lxc_container_set_cgroup_item(lxc.container, ckey, cvalue)) {
		return fmt.Errorf(errSettingCgroupItemFailed, C.GoString(lxc.container.name), key, value)
	}
	return nil
}

// ClearConfigItem clears the value of given key
func (lxc *Container) ClearConfigItem(key string) error {
	lxc.Lock()
	defer lxc.Unlock()

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	if !bool(C.lxc_container_clear_config_item(lxc.container, ckey)) {
		return fmt.Errorf(errClearingCgroupItemFailed, C.GoString(lxc.container.name), key)
	}
	return nil
}

// ConfigKeys returns the name of the config keys
func (lxc *Container) ConfigKeys(key ...string) []string {
	lxc.RLock()
	defer lxc.RUnlock()

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
	lxc.Lock()
	defer lxc.Unlock()

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	if !bool(C.lxc_container_load_config(lxc.container, cpath)) {
		return fmt.Errorf(errLoadConfigFailed, C.GoString(lxc.container.name), path)
	}
	return nil
}

// SaveConfigFile saves the configuration file to given path
func (lxc *Container) SaveConfigFile(path string) error {
	lxc.Lock()
	defer lxc.Unlock()

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	if !bool(C.lxc_container_save_config(lxc.container, cpath)) {
		return fmt.Errorf(errSaveConfigFailed, C.GoString(lxc.container.name), path)
	}
	return nil
}

// ConfigPath returns the configuration file's path
func (lxc *Container) ConfigPath() string {
	lxc.RLock()
	defer lxc.RUnlock()

	return C.GoString(C.lxc_container_get_config_path(lxc.container))
}

// SetConfigPath sets the configuration file's path
func (lxc *Container) SetConfigPath(path string) error {
	lxc.Lock()
	defer lxc.Unlock()

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	if !bool(C.lxc_container_set_config_path(lxc.container, cpath)) {
		return fmt.Errorf(errSettingConfigPathFailed, C.GoString(lxc.container.name), path)
	}
	return nil
}

// MemoryUsage returns memory usage in bytes
func (lxc *Container) MemoryUsage() (ByteSize, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return -1, err
	}

	lxc.RLock()
	defer lxc.RUnlock()

	memUsed, err := strconv.ParseFloat(lxc.CgroupItem("memory.usage_in_bytes")[0], 64)
	if err != nil {
		return -1, fmt.Errorf(errMemLimit)
	}
	return ByteSize(memUsed), err
}

// SwapUsage returns swap usage in bytes
func (lxc *Container) SwapUsage() (ByteSize, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return -1, err
	}

	lxc.RLock()
	defer lxc.RUnlock()

	swapUsed, err := strconv.ParseFloat(lxc.CgroupItem("memory.memsw.usage_in_bytes")[0], 64)
	if err != nil {
		return -1, fmt.Errorf(errSwapLimit)
	}
	return ByteSize(swapUsed), err
}

// MemoryLimit returns memory limit in bytes
func (lxc *Container) MemoryLimit() (ByteSize, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return -1, err
	}

	lxc.RLock()
	defer lxc.RUnlock()

	memLimit, err := strconv.ParseFloat(lxc.CgroupItem("memory.limit_in_bytes")[0], 64)
	if err != nil {
		return -1, fmt.Errorf(errMemLimit)
	}
	return ByteSize(memLimit), err
}

// SetMemoryLimit sets memory limit in bytes
func (lxc *Container) SetMemoryLimit(limit ByteSize) error {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	if err := lxc.SetCgroupItem("memory.limit_in_bytes", fmt.Sprintf("%.f", limit)); err != nil {
		return fmt.Errorf(errSettingMemoryLimitFailed, C.GoString(lxc.container.name))
	}
	return nil
}

// SwapLimit returns the swap limit in bytes
func (lxc *Container) SwapLimit() (ByteSize, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return -1, err
	}

	lxc.RLock()
	defer lxc.RUnlock()

	swapLimit, err := strconv.ParseFloat(lxc.CgroupItem("memory.memsw.limit_in_bytes")[0], 64)
	if err != nil {
		return -1, fmt.Errorf(errSwapLimit)
	}
	return ByteSize(swapLimit), err
}

// SetSwapLimit sets memory limit in bytes
func (lxc *Container) SetSwapLimit(limit ByteSize) error {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	if err := lxc.SetCgroupItem("memory.memsw.limit_in_bytes", fmt.Sprintf("%.f", limit)); err != nil {
		return fmt.Errorf(errSettingSwapLimitFailed, C.GoString(lxc.container.name))
	}
	return nil
}

// CPUTime returns the total CPU time (in nanoseconds) consumed by all tasks in this cgroup (including tasks lower in the hierarchy).
func (lxc *Container) CPUTime() (time.Duration, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return -1, err
	}

	lxc.RLock()
	defer lxc.RUnlock()

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

	lxc.RLock()
	defer lxc.RUnlock()

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

	lxc.RLock()
	defer lxc.RUnlock()

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

	lxc.Lock()
	defer lxc.Unlock()

	ret := int(C.lxc_container_console_getfd(lxc.container, C.int(ttynum)))
	if ret < 0 {
		return -1, fmt.Errorf(errAttachFailed, C.GoString(lxc.container.name))
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

	lxc.Lock()
	defer lxc.Unlock()

	if !bool(C.lxc_container_console(lxc.container, C.int(ttynum), C.int(stdinfd), C.int(stdoutfd), C.int(stderrfd), C.int(escape))) {
		return fmt.Errorf(errAttachFailed, C.GoString(lxc.container.name))
	}
	return nil
}

// AttachRunShell runs a shell inside the container
func (lxc *Container) AttachRunShell() error {
	// FIXME: support lxc_attach_options_t, currently we use LXC_ATTACH_OPTIONS_DEFAULT
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	lxc.Lock()
	defer lxc.Unlock()

	ret := int(C.lxc_container_attach(lxc.container))
	if ret < 0 {
		return fmt.Errorf(errAttachFailed, C.GoString(lxc.container.name))
	}
	return nil
}

// AttachRunCommand runs user specified command inside the container and waits it to exit
func (lxc *Container) AttachRunCommand(args ...string) error {
	// FIXME: support lxc_attach_options_t, currently we use LXC_ATTACH_OPTIONS_DEFAULT
	if args == nil {
		return fmt.Errorf(errInsufficientNumberOfArguments)
	}

	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	lxc.Lock()
	defer lxc.Unlock()

	cargs := makeNullTerminatedArgs(args)
	defer freeNullTerminatedArgs(cargs, len(args))

	ret := int(C.lxc_container_attach_run_wait(lxc.container, cargs))
	if ret < 0 {
		return fmt.Errorf(errAttachFailed, C.GoString(lxc.container.name))
	}
	return nil
}

// Interfaces returns the name of the interfaces from the container
func (lxc *Container) Interfaces() ([]string, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return nil, err
	}

	lxc.RLock()
	defer lxc.RUnlock()

	result := C.lxc_container_get_interfaces(lxc.container)
	if result == nil {
		return nil, fmt.Errorf(errInterfaces, C.GoString(lxc.container.name))
	}
	return convertArgs(result), nil
}

// IPAddress returns the IP address of the given interface
func (lxc *Container) IPAddress(interfaceName string) ([]string, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return nil, err
	}

	lxc.RLock()
	defer lxc.RUnlock()

	cinterface := C.CString(interfaceName)
	defer C.free(unsafe.Pointer(cinterface))

	result := C.lxc_container_get_ips(lxc.container, cinterface, nil, 0)
	if result == nil {
		return nil, fmt.Errorf(errIPAddress, interfaceName, C.GoString(lxc.container.name))
	}
	return convertArgs(result), nil
}

// IPAddresses returns all IP addresses from the container
func (lxc *Container) IPAddresses() ([]string, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return nil, err
	}

	lxc.RLock()
	defer lxc.RUnlock()

	result := C.lxc_container_get_ips(lxc.container, nil, nil, 0)
	if result == nil {
		return nil, fmt.Errorf(errIPAddresses, C.GoString(lxc.container.name))
	}
	return convertArgs(result), nil

}

// IPv4Addresses returns all IPv4 addresses from the container
func (lxc *Container) IPv4Addresses() ([]string, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return nil, err
	}

	lxc.RLock()
	defer lxc.RUnlock()

	cfamily := C.CString("inet")
	defer C.free(unsafe.Pointer(cfamily))

	result := C.lxc_container_get_ips(lxc.container, nil, cfamily, 0)
	if result == nil {
		return nil, fmt.Errorf(errIPv4Addresses, C.GoString(lxc.container.name))
	}
	return convertArgs(result), nil
}

// IPv6Addresses returns all IPv6 addresses from the container
func (lxc *Container) IPv6Addresses() ([]string, error) {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return nil, err
	}

	lxc.RLock()
	defer lxc.RUnlock()

	cfamily := C.CString("inet6")
	defer C.free(unsafe.Pointer(cfamily))

	result := C.lxc_container_get_ips(lxc.container, nil, cfamily, 0)
	if result == nil {
		return nil, fmt.Errorf(errIPv6Addresses, C.GoString(lxc.container.name))
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

// AddDeviceNode adds the device node into given container
func (lxc *Container) AddDeviceNode(source string, destination ...string) error {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	lxc.Lock()
	defer lxc.Unlock()

	csource := C.CString(source)
	defer C.free(unsafe.Pointer(csource))

	if destination != nil && len(destination) == 1 {
		cdestination := C.CString(destination[0])
		defer C.free(unsafe.Pointer(cdestination))

		if !bool(C.lxc_container_add_device_node(lxc.container, csource, cdestination)) {
			return fmt.Errorf("adding device %s to container %q failed", source, C.GoString(lxc.container.name))
		}
		return nil
	}

	if !bool(C.lxc_container_add_device_node(lxc.container, csource, nil)) {
		return fmt.Errorf("adding device %s to container %q failed", source, C.GoString(lxc.container.name))
	}
	return nil

}

// RemoveDeviceNode removes the device node from given container
func (lxc *Container) RemoveDeviceNode(source string, destination ...string) error {
	if err := lxc.ensureDefinedAndRunning(); err != nil {
		return err
	}

	lxc.Lock()
	defer lxc.Unlock()

	csource := C.CString(source)
	defer C.free(unsafe.Pointer(csource))

	if destination != nil && len(destination) == 1 {
		cdestination := C.CString(destination[0])
		defer C.free(unsafe.Pointer(cdestination))

		if !bool(C.lxc_container_remove_device_node(lxc.container, csource, cdestination)) {
			return fmt.Errorf("adding device %s to container %q failed", source, C.GoString(lxc.container.name))
		}
		return nil
	}

	if !bool(C.lxc_container_remove_device_node(lxc.container, csource, nil)) {
		return fmt.Errorf("adding device %s to container %q failed", source, C.GoString(lxc.container.name))
	}
	return nil
}
