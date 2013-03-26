package lxc

/*
#cgo linux CFLAGS: -I/usr/include
#cgo linux LDFLAGS: -L/usr/lib/x86_64-linux-gnu -llxc -lutil
#include <lxc/lxc.h>
#include <lxc/lxccontainer.h>
#include "lxc.h"
*/
import "C"

import (
	"unsafe"
)

type Container struct {
	container *C.struct_lxc_container
}

func NewContainer(name string) Container {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return Container{C.lxc_container_new(cname)}
}

func (c *Container) Defined() bool {
	return bool(C.container_defined(c.container))
}

func (c *Container) Running() bool {
	return bool(C.container_running(c.container))
}

func (c *Container) State() string {
	return C.GoString(C.container_state(c.container))
}

func (c *Container) InitPID() int {
	return int(C.container_init_pid(c.container))
}

func (c *Container) Daemonize() bool {
	return bool(c.container.daemonize != 0)
}

func (c *Container) SetDaemonize() {
	C.container_want_daemonize(c.container)
}

func (c *Container) Freeze() bool {
	return bool(C.container_freeze(c.container))
}

func (c *Container) Unfreeze() bool {
	return bool(C.container_unfreeze(c.container))
}

func (c *Container) Create(template string) bool {
	ctemplate := C.CString(template)
	defer C.free(unsafe.Pointer(ctemplate))
	return bool(C.container_create(c.container, ctemplate, nil))
}

func (c *Container) Start(useinit bool) bool {
	cuseinit := 0
	if useinit {
		cuseinit = 1
	}
	return bool(C.container_start(c.container, C.int(cuseinit), nil))
}

func (c *Container) Stop() bool {
	return bool(C.container_stop(c.container))
}

func (c *Container) Shutdown(timeout int) bool {
	return bool(C.container_shutdown(c.container, C.int(timeout)))
}

func (c *Container) ConfigFileName() string {
	return C.GoString(C.container_config_file_name(c.container))
}
