package lxc

var (
	ErrAddDeviceNodeFailed           = NewError("adding device to container failed")
	ErrAlreadyDefined                = NewError("container already defined")
	ErrAlreadyFrozen                 = NewError("container is already frozen")
	ErrAlreadyRunning                = NewError("container is already running")
	ErrAttachFailed                  = NewError("attaching to the container failed")
	ErrBlkioUsage                    = NewError("BlkioUsage for the container failed")
	ErrClearingCgroupItemFailed      = NewError("clearing cgroup item for the container failed")
	ErrCloneFailed                   = NewError("cloning the container failed")
	ErrCloseAllFdsFailed             = NewError("setting close_all_fds flag for container failed")
	ErrCreateFailed                  = NewError("creating the container failed")
	ErrCreateSnapshotFailed          = NewError("snapshotting the container failed")
	ErrDaemonizeFailed               = NewError("setting daemonize flag for container failed")
	ErrDestroyFailed                 = NewError("destroying the container failed")
	ErrDestroySnapshotFailed         = NewError("destroying the snapshot failed")
	ErrExecuteFailed                 = NewError("executing the command in a temporary container failed")
	ErrFreezeFailed                  = NewError("freezing the container failed")
	ErrInsufficientNumberOfArguments = NewError("insufficient number of arguments were supplied")
	ErrInterfaces                    = NewError("getting interface names for the container failed")
	ErrIPAddresses                   = NewError("getting IP addresses of the container failed")
	ErrIPAddress                     = NewError("getting IP address on the interface of the container failed")
	ErrIPv4Addresses                 = NewError("getting IPv4 addresses of the container failed")
	ErrIPv6Addresses                 = NewError("getting IPv6 addresses of the container failed")
	ErrKMemLimit                     = NewError("your kernel does not support cgroup memory controller")
	ErrLoadConfigFailed              = NewError("loading config file for the container failed")
	ErrMemLimit                      = NewError("your kernel does not support cgroup memory controller")
	ErrNewFailed                     = NewError("allocating the container failed")
	ErrNotDefined                    = NewError("container is not defined")
	ErrNotFrozen                     = NewError("container is not frozen")
	ErrNotRunning                    = NewError("container is not running")
	ErrRebootFailed                  = NewError("rebooting the container failed")
	ErrRemoveDeviceNodeFailed        = NewError("removing device from container failed")
	ErrRenameFailed                  = NewError("renaming the container failed")
	ErrRestoreSnapshotFailed         = NewError("restoring the container failed")
	ErrSaveConfigFailed              = NewError("saving config file for the container failed")
	ErrSettingCgroupItemFailed       = NewError("setting cgroup item for the container failed")
	ErrSettingConfigItemFailed       = NewError("setting config item for the container failed")
	ErrSettingConfigPathFailed       = NewError("setting config file for the container failed")
	ErrSettingKMemoryLimitFailed     = NewError("setting kernel memory limit for the container failed")
	ErrSettingMemoryLimitFailed      = NewError("setting memory limit for the container failed")
	ErrSettingMemorySwapLimitFailed  = NewError("setting memroy+swap limit for the container failed")
	ErrShutdownFailed                = NewError("shutting down the container failed")
	ErrStartFailed                   = NewError("starting the container failed")
	ErrStopFailed                    = NewError("stopping the container failed")
	ErrMemorySwapLimit               = NewError("your kernel does not support cgroup swap controller")
	ErrUnfreezeFailed                = NewError("unfreezing the container failed")
)

type Error struct {
	Message string
}

func NewError(msg string) *Error {
	return &Error{
		Message: msg,
	}
}

func (e Error) Error() string {
	return e.Message
}
