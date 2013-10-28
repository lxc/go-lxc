// +build linux

/*
 * type.go: Go bindings for lxc
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

package lxc

// #include <lxc/lxc.h>
// #include <lxc/lxccontainer.h>
import "C"

import (
	"fmt"
)

// Verbosity type
type Verbosity int

const (
	// Quiet makes some API calls not to write anything to stdout
	Quiet Verbosity = 1 << iota
	// Verbose makes some API calls write to stdout
	Verbose
)

// BackendStore type
type BackendStore int

const (
	// BtrFS backendstore type
	BtrFS BackendStore = iota
	// Directory backendstore type
	Directory
	// LVM backendstore type
	LVM
	// ZFS backendstore type
	ZFS
	// OverlayFS backendstore type
	OverlayFS
	// Loopback backendstore type
	Loopback
)

// BackendStore as string
func (t BackendStore) String() string {
	switch t {
	case Directory:
		return "dir"
	case ZFS:
		return "zfs"
	case BtrFS:
		return "btrfs"
	case LVM:
		return "lvm"
	case OverlayFS:
		return "overlayfs"
	case Loopback:
		return "loopback"
	}
	return "<INVALID>"
}

// State type
type State int

const (
	// STOPPED means container is not running
	STOPPED State = C.STOPPED
	// STARTING means container is starting
	STARTING State = C.STARTING
	// RUNNING means container is running
	RUNNING State = C.RUNNING
	// STOPPING means container is stopping
	STOPPING State = C.STOPPING
	// ABORTING means container is aborting
	ABORTING State = C.ABORTING
	// FREEZING means container is freezing
	FREEZING State = C.FREEZING
	// FROZEN means containe is frozen
	FROZEN State = C.FROZEN
	// THAWED means container is thawed
	THAWED State = C.THAWED
)

var stateMap = map[string]State{
	"STOPPED":  STOPPED,
	"STARTING": STARTING,
	"RUNNING":  RUNNING,
	"STOPPING": STOPPING,
	"ABORTING": ABORTING,
	"FREEZING": FREEZING,
	"FROZEN":   FROZEN,
	"THAWED":   THAWED,
}

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
	}
	return "<INVALID>"
}

// Taken from http://golang.org/doc/effective_go.html#constants

// ByteSize type
type ByteSize float64

const (
	_ = iota
	// KB - kilobyte
	KB ByteSize = 1 << (10 * iota)
	// MB - megabyte
	MB
	// GB - gigabyte
	GB
	// TB - terabyte
	TB
	// PB - petabyte
	PB
	// EB - exabyte
	EB
	// ZB - zettabyte
	ZB
	// YB - yottabyte
	YB
)

func (b ByteSize) String() string {
	switch {
	case b >= YB:
		return fmt.Sprintf("%.2fYB", b/YB)
	case b >= ZB:
		return fmt.Sprintf("%.2fZB", b/ZB)
	case b >= EB:
		return fmt.Sprintf("%.2fEB", b/EB)
	case b >= PB:
		return fmt.Sprintf("%.2fPB", b/PB)
	case b >= TB:
		return fmt.Sprintf("%.2fTB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.2fGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2fMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2fKB", b/KB)
	}
	return fmt.Sprintf("%.2fB", b)
}
