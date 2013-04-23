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

//Go (golang) Bindings for LXC (Linux Containers)
//
//This package implements Go bindings for the LXC C API.
package lxc

// #cgo linux LDFLAGS: -llxc -lutil
// #include <lxc/lxc.h>
// #include <lxc/lxccontainer.h>
// #include "lxc.h"
import "C"

import (
	"path/filepath"
	"unsafe"
)

const (
	// Timeout
	WAIT_FOREVER int = iota - 1
	DONT_WAIT

	LXC_NETWORK_KEY = "lxc.network"
)

func NewContainer(name string) *Container {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return &Container{container: C.lxc_container_new(cname, nil)}
}

// Returns LXC version
func Version() string {
	return C.GoString(C.lxc_get_version())
}

// Returns default config path
func DefaultConfigPath() string {
	return C.GoString(C.lxc_get_default_config_path())
}

// Returns the names of containers on the system.
func ContainerNames() []string {
	// FIXME: Support custom config paths
	matches, err := filepath.Glob(filepath.Join(DefaultConfigPath(), "/*/config"))
	if err != nil {
		return nil
	}

	for i, v := range matches {
		matches[i] = filepath.Base(filepath.Dir(v))
	}
	return matches
}

// Returns the containers on the system.
func Containers() []Container {
	var containers []Container

	for _, v := range ContainerNames() {
		containers = append(containers, *NewContainer(v))
	}
	return containers
}
