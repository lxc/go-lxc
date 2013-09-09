/*
 * lxc.go: Go bindings for lxc
 *
 * Copyright © 2013, S.Çağlar Onur
 *
 * Authors:
 * S.Çağlar Onur <caglar@10ur.org>
 *
 * This library is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 2.1 of the License, or (at your option) any later version.

 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.

 * You should have received a copy of the GNU Lesser General Public
 * License along with this library; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301  USA
 */

// Package lxc provides Go (golang) Bindings for LXC (Linux Containers) C API.
package lxc

// #cgo linux LDFLAGS: -llxc -lutil
// #include <lxc/lxc.h>
// #include <lxc/lxccontainer.h>
// #include "lxc.h"
import "C"

import (
	"os"
	"path/filepath"
	"unsafe"
)

const (
	// WaitForever timeout
	WaitForever int = iota - 1
	// DontWait timeout
	DontWait
)

const (
	// CloneKeepName means don't edit the rootfs to change the hostname.
	CloneKeepName int = 1 << iota
	// CloneCopyHooks means copy all hooks into the container directory.
	CloneCopyHooks
	// CloneKeepMACAddr means don't change the mac address on network interfaces.
	CloneKeepMACAddr
	// CloneSnapshot means snapshot the original filesystem(s).
	CloneSnapshot
)

func init() {
	if os.Geteuid() != 0 {
		panic("Running as non-root.")
	}
}

// NewContainer returns a new container struct
func NewContainer(name string) *Container {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	return &Container{container: C.lxc_container_new(cname, nil), verbosity: Quiet}
}

// GetContainer increments reference counter of the container object
func GetContainer(lxc *Container) bool {
	return C.lxc_container_get(lxc.container) == 1
}

// PutContainer decrements reference counter of the container object
func PutContainer(lxc *Container) bool {
	return C.lxc_container_put(lxc.container) == 1
}

// Version returns LXC version
func Version() string {
	return C.GoString(C.lxc_get_version())
}

// DefaultConfigPath returns default config path
func DefaultConfigPath() string {
	return C.GoString(C.lxc_get_default_config_path())
}

// DefaultLvmVg returns default LVM volume group
func DefaultLvmVg() string {
	return C.GoString(C.lxc_get_default_lvm_vg())
}

// DefaultZfsRoot returns default ZFS root
func DefaultZfsRoot() string {
	return C.GoString(C.lxc_get_default_zfs_root())
}

// ContainerNames returns the names of containers on the system.
func ContainerNames(paths ...string) []string {
	if paths == nil {
		matches, err := filepath.Glob(filepath.Join(DefaultConfigPath(), "/*/config"))
		if err != nil {
			return nil
		}

		for i, v := range matches {
			matches[i] = filepath.Base(filepath.Dir(v))
		}
		return matches
	}
	// FIXME: Support custom config paths
	return nil
}

// Containers returns the containers on the system.
func Containers() []Container {
	var containers []Container

	for _, v := range ContainerNames() {
		containers = append(containers, *NewContainer(v))
	}
	return containers
}
