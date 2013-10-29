// Copyright © 2013, S.Çağlar Onur
// Use of this source code is governed by a LGPLv2.1
// license that can be found in the LICENSE file.
//
// Authors:
// S.Çağlar Onur <caglar@10ur.org>

// +build linux

// Package lxc provides Go (golang) Bindings for LXC (Linux Containers) C API.
package lxc

// #cgo pkg-config: lxc
// #include <lxc/lxc.h>
// #include <lxc/lxccontainer.h>
// #include "lxc.h"
import "C"

import (
	"os"
	"unsafe"
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
		var cnames **C.char

		size := int(C.lxc_container_list_defined_containers(nil, &cnames))
		if size < 1 {
			return nil
		}
		return convertNArgs(cnames, size)
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

// ActiveContainerNames returns the names of the active containers on the system.
func ActiveContainerNames(paths ...string) []string {
	if paths == nil {
		var cnames **C.char

		size := int(C.lxc_container_list_active_containers(nil, &cnames))
		if size < 1 {
			return nil
		}
		return convertNArgs(cnames, size)
	}
	// FIXME: Support custom config paths
	return nil
}

// ActiveContainers returns the active containers on the system.
func ActiveContainers() []Container {
	var containers []Container

	for _, v := range ActiveContainerNames() {
		containers = append(containers, *NewContainer(v))
	}
	return containers
}
