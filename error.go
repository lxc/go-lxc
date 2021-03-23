// Copyright © 2013, 2014, The Go-LXC Authors. All rights reserved.
// Use of this source code is governed by a LGPLv2.1
// license that can be found in the LICENSE file.

// +build linux,cgo

package lxc

const (
	// ErrAddDeviceNodeFailed - adding device to container failed
	ErrAddDeviceNodeFailed = lxcError("adding device to container failed")

	// ErrAllocationFailed - allocating memory failed
	ErrAllocationFailed = lxcError("allocating memory failed")

	// ErrAlreadyDefined - container already defined
	ErrAlreadyDefined = lxcError("container already defined")

	// ErrAlreadyFrozen - container is already frozen
	ErrAlreadyFrozen = lxcError("container is already frozen")

	// ErrAlreadyRunning - container is already running
	ErrAlreadyRunning = lxcError("container is already running")

	// ErrAttachFailed - attaching to the container failed
	ErrAttachFailed = lxcError("attaching to the container failed")

	// ErrAttachInterfaceFailed - attaching specified netdev to the container failed
	ErrAttachInterfaceFailed = lxcError("attaching specified netdev to the container failed")

	// ErrBlkioUsage - BlkioUsage for the container failed
	ErrBlkioUsage = lxcError("BlkioUsage for the container failed")

	// ErrCheckpointFailed - checkpoint failed
	ErrCheckpointFailed = lxcError("checkpoint failed")

	// ErrClearingConfigItemFailed - clearing config item for the container failed
	ErrClearingConfigItemFailed = lxcError("clearing config item for the container failed")

	// ErrClearingCgroupItemFailed - clearing cgroup item for the container failed
	ErrClearingCgroupItemFailed = lxcError("clearing cgroup item for the container failed")

	// ErrCloneFailed - cloning the container failed
	ErrCloneFailed = lxcError("cloning the container failed")

	// ErrCloseAllFdsFailed - setting close_all_fds flag for container failed
	ErrCloseAllFdsFailed = lxcError("setting close_all_fds flag for container failed")

	// ErrCreateFailed - creating the container failed
	ErrCreateFailed = lxcError("creating the container failed")

	// ErrCreateSnapshotFailed - snapshotting the container failed
	ErrCreateSnapshotFailed = lxcError("snapshotting the container failed")

	// ErrDaemonizeFailed - setting daemonize flag for container failed
	ErrDaemonizeFailed = lxcError("setting daemonize flag for container failed")

	// ErrDestroyAllSnapshotsFailed - destroying all snapshots failed
	ErrDestroyAllSnapshotsFailed = lxcError("destroying all snapshots failed")

	// ErrDestroyFailed - destroying the container failed
	ErrDestroyFailed = lxcError("destroying the container failed")

	// ErrDestroySnapshotFailed - destroying the snapshot failed
	ErrDestroySnapshotFailed = lxcError("destroying the snapshot failed")

	// ErrDestroyWithAllSnapshotsFailed - destroying the container with all snapshots failed
	ErrDestroyWithAllSnapshotsFailed = lxcError("destroying the container with all snapshots failed")

	// ErrDetachInterfaceFailed - detaching specified netdev to the container failed
	ErrDetachInterfaceFailed = lxcError("detaching specified netdev to the container failed")

	// ErrExecuteFailed - executing the command in a temporary container failed
	ErrExecuteFailed = lxcError("executing the command in a temporary container failed")

	// ErrFreezeFailed - freezing the container failed
	ErrFreezeFailed = lxcError("freezing the container failed")

	// ErrInsufficientNumberOfArguments - insufficient number of arguments were supplied
	ErrInsufficientNumberOfArguments = lxcError("insufficient number of arguments were supplied")

	// ErrInterfaces - getting interface names for the container failed
	ErrInterfaces = lxcError("getting interface names for the container failed")

	// ErrIPAddresses - getting IP addresses of the container failed
	ErrIPAddresses = lxcError("getting IP addresses of the container failed")

	// ErrIPAddress - getting IP address on the interface of the container failed
	ErrIPAddress = lxcError("getting IP address on the interface of the container failed")

	// ErrIPv4Addresses - getting IPv4 addresses of the container failed
	ErrIPv4Addresses = lxcError("getting IPv4 addresses of the container failed")

	// ErrIPv6Addresses - getting IPv6 addresses of the container failed
	ErrIPv6Addresses = lxcError("getting IPv6 addresses of the container failed")

	// ErrKMemLimit - your kernel does not support cgroup kernel memory controller
	ErrKMemLimit = lxcError("your kernel does not support cgroup kernel memory controller")

	// ErrLoadConfigFailed - loading config file for the container failed
	ErrLoadConfigFailed = lxcError("loading config file for the container failed")

	// ErrMemLimit - your kernel does not support cgroup memory controller
	ErrMemLimit = lxcError("your kernel does not support cgroup memory controller")

	// ErrMemorySwapLimit - your kernel does not support cgroup swap controller
	ErrMemorySwapLimit = lxcError("your kernel does not support cgroup swap controller")

	// ErrMethodNotAllowed - the requested method is not currently supported with unprivileged containers
	ErrMethodNotAllowed = lxcError("the requested method is not currently supported with unprivileged containers")

	// ErrNewFailed - allocating the container failed
	ErrNewFailed = lxcError("allocating the container failed")

	// ErrNoSnapshot - container has no snapshot
	ErrNoSnapshot = lxcError("container has no snapshot")

	// ErrNotDefined - container is not defined
	ErrNotDefined = lxcError("container is not defined")

	// ErrNotFrozen - container is not frozen
	ErrNotFrozen = lxcError("container is not frozen")

	// ErrNotRunning - container is not running
	ErrNotRunning = lxcError("container is not running")

	// ErrNotSupported - method is not supported by this LXC version
	ErrNotSupported = lxcError("method is not supported by this LXC version")

	// ErrRebootFailed - rebooting the container failed
	ErrRebootFailed = lxcError("rebooting the container failed")

	// ErrRemoveDeviceNodeFailed - removing device from container failed
	ErrRemoveDeviceNodeFailed = lxcError("removing device from container failed")

	// ErrRenameFailed - renaming the container failed
	ErrRenameFailed = lxcError("renaming the container failed")

	// ErrRestoreFailed - restore failed
	ErrRestoreFailed = lxcError("restore failed")

	// ErrRestoreSnapshotFailed - restoring the container failed
	ErrRestoreSnapshotFailed = lxcError("restoring the container failed")

	// ErrSaveConfigFailed - saving config file for the container failed
	ErrSaveConfigFailed = lxcError("saving config file for the container failed")

	// ErrSettingCgroupItemFailed - setting cgroup item for the container failed
	ErrSettingCgroupItemFailed = lxcError("setting cgroup item for the container failed")

	// ErrSettingConfigItemFailed - setting config item for the container failed
	ErrSettingConfigItemFailed = lxcError("setting config item for the container failed")

	// ErrSettingConfigPathFailed - setting config file for the container failed
	ErrSettingConfigPathFailed = lxcError("setting config file for the container failed")

	// ErrSettingKMemoryLimitFailed - setting kernel memory limit for the container failed
	ErrSettingKMemoryLimitFailed = lxcError("setting kernel memory limit for the container failed")

	// ErrSettingMemoryLimitFailed - setting memory limit for the container failed
	ErrSettingMemoryLimitFailed = lxcError("setting memory limit for the container failed")

	// ErrSettingMemorySwapLimitFailed - setting memory+swap limit for the container failed
	ErrSettingMemorySwapLimitFailed = lxcError("setting memory+swap limit for the container failed")

	// ErrSettingSoftMemoryLimitFailed - setting soft memory limit for the container failed
	ErrSettingSoftMemoryLimitFailed = lxcError("setting soft memory limit for the container failed")

	// ErrShutdownFailed - shutting down the container failed
	ErrShutdownFailed = lxcError("shutting down the container failed")

	// ErrSoftMemLimit - your kernel does not support cgroup memory controller
	ErrSoftMemLimit = lxcError("your kernel does not support cgroup memory controller")

	// ErrStartFailed - starting the container failed
	ErrStartFailed = lxcError("starting the container failed")

	// ErrStopFailed - stopping the container failed
	ErrStopFailed = lxcError("stopping the container failed")

	// ErrTemplateNotAllowed - unprivileged users only allowed to use "download" template
	ErrTemplateNotAllowed = lxcError("unprivileged users only allowed to use \"download\" template")

	// ErrUnfreezeFailed - unfreezing the container failed
	ErrUnfreezeFailed = lxcError("unfreezing the container failed")

	// ErrUnknownBackendStore - unknown backend type
	ErrUnknownBackendStore = lxcError("unknown backend type")

	// ErrReleaseFailed - releasing the container failed
	ErrReleaseFailed = lxcError("releasing the container failed")
)

type lxcError string

func (e lxcError) Error() string {
	return string(e)
}
