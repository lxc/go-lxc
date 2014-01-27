// Copyright © 2013, S.Çağlar Onur
// Use of this source code is governed by a LGPLv2.1
// license that can be found in the LICENSE file.
//
// Authors:
// S.Çağlar Onur <caglar@10ur.org>

// +build linux

package lxc

const (
	// WaitForever timeout
	WaitForever int = iota - 1
	// DontWait timeout
	DontWait
)

const (
	isDefined = 1 << iota
	isNotDefined
	isRunning
	isNotRunning
)

const (
	errAddDeviceNodeFailed           = "adding device %s to container %q failed"
	errAlreadyDefined                = "container %q already defined"
	errAlreadyFrozen                 = "container %q is already frozen"
	errAlreadyRunning                = "container %q is already running"
	errAttachFailed                  = "attaching to the container %q failed"
	errBlkioUsage                    = "BlkioUsage for the container %q failed"
	errClearingCgroupItemFailed      = "clearing cgroup item for the container %q failed (key: %s)"
	errCloneFailed                   = "cloning the container %q failed"
	errCloseAllFdsFailed             = "setting close_all_fds flag for container %q failed"
	errCreateFailed                  = "creating the container %q failed"
	errCreateSnapshotFailed          = "snapshotting the container %q failed"
	errDaemonizeFailed               = "setting daemonize flag for container %q failed"
	errDestroyFailed                 = "destroying the container %q failed"
	errDestroySnapshotFailed         = "destroying the snapshot %q failed"
	errExecuteFailed                 = "executing the command in a temporary container %q failed"
	errFreezeFailed                  = "freezing the container %q failed"
	errInsufficientNumberOfArguments = "insufficient number of arguments were supplied"
	errInterfaces                    = "getting interface names for the container %q failed"
	errIPAddresses                   = "getting IP addresses of the container %q failed"
	errIPAddress                     = "getting IP address on the interface %s of the container %q failed"
	errIPv4Addresses                 = "getting IPv4 addresses of the container %q failed"
	errIPv6Addresses                 = "getting IPv6 addresses of the container %q failed"
	errKMemLimit                     = "your kernel does not support cgroup memory controller"
	errLoadConfigFailed              = "loading config file for the container %q failed (path: %s)"
	errMemLimit                      = "your kernel does not support cgroup memory controller"
	errNewFailed                     = "allocating the container %q failed"
	errNotDefined                    = "there is no container named %q"
	errNotFrozen                     = "container %q is not frozen"
	errNotRunning                    = "container %q is not running"
	errRebootFailed                  = "rebooting the container %q failed"
	errRemoveDeviceNodeFailed        = "removing device %s from container %q failed"
	errRenameFailed                  = "renaming the container %q failed"
	errRestoreSnapshotFailed         = "restoring the container %q failed"
	errSaveConfigFailed              = "saving config file for the container %q failed (path: %s)"
	errSettingCgroupItemFailed       = "setting cgroup item for the container %q failed (key: %s, value: %s)"
	errSettingConfigItemFailed       = "setting config item for the container %q failed (key: %s, value: %s)"
	errSettingConfigPathFailed       = "setting config file for the container %q failed (path: %s)"
	errSettingKMemoryLimitFailed     = "setting kernel memory limit for the container %q failed"
	errSettingMemoryLimitFailed      = "setting memory limit for the container %q failed"
	errSettingMemorySwapLimitFailed  = "setting memroy+swap limit for the container %q failed"
	errShutdownFailed                = "shutting down the container %q failed"
	errStartFailed                   = "starting the container %q failed"
	errStopFailed                    = "stopping the container %q failed"
	errMemorySwapLimit               = "your kernel does not support cgroup swap controller"
	errUnfreezeFailed                = "unfreezing the container %q failed"
)
