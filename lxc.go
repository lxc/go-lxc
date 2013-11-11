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
	"fmt"
	"os"
	"unsafe"
)

func init() {
	if os.Geteuid() != 0 {
		panic("Running as non-root.")
	}
}

// NewContainer returns a new container struct
func NewContainer(name string, lxcpath ...string) (*Container, error) {
	var container *C.struct_lxc_container

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	if lxcpath != nil && len(lxcpath) == 1 {
		clxcpath := C.CString(lxcpath[0])
		defer C.free(unsafe.Pointer(clxcpath))

		container = C.lxc_container_new(cname, clxcpath)
	} else {
		container = C.lxc_container_new(cname, nil)
	}

	if container == nil {
		return nil, fmt.Errorf(errNewFailed, name)
	}
	return &Container{container: container, verbosity: Quiet}, nil
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

// ContainerNames returns the names of defined and active containers on the system.
func ContainerNames(lxcpath ...string) []string {
	var size int
	var cnames **C.char

	if lxcpath != nil && len(lxcpath) == 1 {
		clxcpath := C.CString(lxcpath[0])
		defer C.free(unsafe.Pointer(clxcpath))

		size = int(C.list_all_containers(clxcpath, &cnames, nil))
	} else {

		size = int(C.list_all_containers(nil, &cnames, nil))
	}

	if size < 1 {
		return nil
	}
	return convertNArgs(cnames, size)
}

// Containers returns the defined and active containers on the system.
func Containers(lxcpath ...string) []Container {
	var containers []Container

	for _, v := range ContainerNames(lxcpath...) {
		container, err := NewContainer(v, lxcpath...)
		if err != nil {
			return nil
		}
		containers = append(containers, *container)
	}
	return containers
}

// DefinedContainerNames returns the names of the defined containers on the system.
func DefinedContainerNames(lxcpath ...string) []string {
	var size int
	var cnames **C.char

	if lxcpath != nil && len(lxcpath) == 1 {
		clxcpath := C.CString(lxcpath[0])
		defer C.free(unsafe.Pointer(clxcpath))

		size = int(C.list_defined_containers(clxcpath, &cnames, nil))
	} else {

		size = int(C.list_defined_containers(nil, &cnames, nil))
	}

	if size < 1 {
		return nil
	}
	return convertNArgs(cnames, size)
}

// DefinedContainers returns the defined containers on the system.
func DefinedContainers(lxcpath ...string) []Container {
	var containers []Container

	for _, v := range DefinedContainerNames(lxcpath...) {
		container, err := NewContainer(v, lxcpath...)
		if err != nil {
			return nil
		}
		containers = append(containers, *container)
	}
	return containers
}

// ActiveContainerNames returns the names of the active containers on the system.
func ActiveContainerNames(lxcpath ...string) []string {
	var size int
	var cnames **C.char

	if lxcpath != nil && len(lxcpath) == 1 {
		clxcpath := C.CString(lxcpath[0])
		defer C.free(unsafe.Pointer(clxcpath))

		size = int(C.list_active_containers(clxcpath, &cnames, nil))
	} else {

		size = int(C.list_active_containers(nil, &cnames, nil))
	}

	if size < 1 {
		return nil
	}
	return convertNArgs(cnames, size)
}

// ActiveContainers returns the active containers on the system.
func ActiveContainers(lxcpath ...string) []Container {
	var containers []Container

	for _, v := range ActiveContainerNames(lxcpath...) {
		container, err := NewContainer(v, lxcpath...)
		if err != nil {
			return nil
		}
		containers = append(containers, *container)
	}
	return containers
}
