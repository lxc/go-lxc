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

// #cgo linux LDFLAGS: -llxc -lutil
// #include <lxc/lxc.h>
// #include <lxc/lxccontainer.h>
// #include "lxc.h"
import "C"

import (
	"strings"
	"unsafe"
)

type State int

const (
	// Timeout
	WAIT_FOREVER int = iota
	DONT_WAIT
	// State
	STOPPED   = C.STOPPED
	STARTING  = C.STARTING
	RUNNING   = C.RUNNING
	STOPPING  = C.STOPPING
	ABORTING  = C.ABORTING
	FREEZING  = C.FREEZING
	FROZEN    = C.FROZEN
	THAWED    = C.THAWED
	MAX_STATE = C.MAX_STATE
)

// State as string
func (t State) String() string {
	switch t {
	case STOPPED:
		return "STOPPED"
	case STARTING:
		return "STARTING"
	case RUNNING:
		return "RUNNING"
	case STOPPING:
		return "STOPPING"
	case ABORTING:
		return "ABORTING"
	case FREEZING:
		return "FREEZING"
	case FROZEN:
		return "FROZEN"
	case THAWED:
		return "THAWED"
	case MAX_STATE:
		return "MAX_STATE"
	}
	return "<INVALID>"
}

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

func (lxc *Container) GetName() string {
	return C.GoString(lxc.container.name)
}

func (lxc *Container) Defined() bool {
	return bool(C.lxc_container_defined(lxc.container))
}

func (lxc *Container) Running() bool {
	return bool(C.lxc_container_running(lxc.container))
}

func (lxc *Container) GetState() string {
	return C.GoString(C.lxc_container_state(lxc.container))
}

func (lxc *Container) GetInitPID() int {
	return int(C.lxc_container_init_pid(lxc.container))
}

func (lxc *Container) GetDaemonize() bool {
	return bool(lxc.container.daemonize != 0)
}

func (lxc *Container) SetDaemonize() {
	C.lxc_container_want_daemonize(lxc.container)
}

func (lxc *Container) Freeze() bool {
	return bool(C.lxc_container_freeze(lxc.container))
}

func (lxc *Container) Unfreeze() bool {
	return bool(C.lxc_container_unfreeze(lxc.container))
}

func (lxc *Container) Create(template string, args []string) bool {
	ctemplate := C.CString(template)
	defer C.free(unsafe.Pointer(ctemplate))
	if args != nil {
		cargs := makeArgs(args)
		defer freeArgs(cargs)
		return bool(C.lxc_container_create(lxc.container, ctemplate, &cargs[0]))
	}
	return bool(C.lxc_container_create(lxc.container, ctemplate, nil))
}

func (lxc *Container) Start(useinit bool, args []string) bool {
	cuseinit := 0
	if useinit {
		cuseinit = 1
	}
	if args != nil {
		cargs := makeArgs(args)
		defer freeArgs(cargs)
		return bool(C.lxc_container_start(lxc.container, C.int(cuseinit), &cargs[0]))
	}
	return bool(C.lxc_container_start(lxc.container, C.int(cuseinit), nil))
}

func (lxc *Container) Stop() bool {
	return bool(C.lxc_container_stop(lxc.container))
}

func (lxc *Container) Shutdown(timeout int) bool {
	return bool(C.lxc_container_shutdown(lxc.container, C.int(timeout)))
}

func (lxc *Container) Destroy() bool {
	return bool(C.lxc_container_destroy(lxc.container))
}

func (lxc *Container) Wait(state State, timeout int) bool {
	cstate := C.CString(state.String())
	defer C.free(unsafe.Pointer(cstate))
	return bool(C.lxc_container_wait(lxc.container, cstate, C.int(timeout)))
}

func (lxc *Container) GetConfigFileName() string {
	return C.GoString(C.lxc_container_config_file_name(lxc.container))
}

func (lxc *Container) GetConfigItem(key string) []string {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	return strings.Split(C.GoString(C.lxc_container_get_config_item(lxc.container, ckey)), "\n")
}

func (lxc *Container) SetConfigItem(key string, value string) bool {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	return bool(C.lxc_container_set_config_item(lxc.container, ckey, cvalue))
}

func (lxc *Container) ClearConfigItem(key string) bool {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	return bool(C.lxc_container_clear_config_item(lxc.container, ckey))
}

func (lxc *Container) GetKeys(key string) []string {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	return strings.Split(C.GoString(C.lxc_container_get_keys(lxc.container, ckey)), "\n")
}

func (lxc *Container) LoadConfigFile(path string) bool {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return bool(C.lxc_container_load_config(lxc.container, cpath))
}

func (lxc *Container) SaveConfigFile(path string) bool {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return bool(C.lxc_container_save_config(lxc.container, cpath))
}
