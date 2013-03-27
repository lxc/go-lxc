/*
 * lxc.go: Go bindings for lxc
 *
 * Copyright © 2013, S.Çağlar Onur
 *
 * Authors:
 * S.Çağlar Onur <caglar@10ur.org>
 *
 * This library is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2, as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License along
 * with this program; if not, write to the Free Software Foundation, Inc.,
 * 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
 */

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

func makeArgs(args []string) []*C.char {
	ret := make([]*C.char, len(args))
	for i, s := range args {
		ret[i] = C.CString(s)
	}
	return ret
}

func freeArgs(cArgs []*C.char) {
	for _, s := range cArgs {
		C.free(unsafe.Pointer(s))
	}
}

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

func (c *Container) Create(template string, args []string) bool {
	ctemplate := C.CString(template)
	defer C.free(unsafe.Pointer(ctemplate))
	if args != nil {
		cargs := makeArgs(args)
		defer freeArgs(cargs)
		return bool(C.container_create(c.container, ctemplate, &cargs[0]))
	}
	return bool(C.container_create(c.container, ctemplate, nil))
}

func (c *Container) Start(useinit bool, args []string) bool {
	cuseinit := 0
	if useinit {
		cuseinit = 1
	}
	if args != nil {
		cargs := makeArgs(args)
		defer freeArgs(cargs)
		return bool(C.container_start(c.container, C.int(cuseinit), &cargs[0]))
	}
	return bool(C.container_start(c.container, C.int(cuseinit), nil))
}

func (c *Container) Stop() bool {
	return bool(C.container_stop(c.container))
}

func (c *Container) Shutdown(timeout int) bool {
	return bool(C.container_shutdown(c.container, C.int(timeout)))
}

func (c *Container) Destroy() bool {
	return bool(C.container_destroy(c.container))
}

func (c *Container) ConfigFileName() string {
	return C.GoString(C.container_config_file_name(c.container))
}
