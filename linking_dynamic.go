// Copyright © 2013, 2014, The Go-LXC Authors. All rights reserved.
// Use of this source code is governed by a LGPLv2.1
// license that can be found in the LICENSE file.

//go:build linux && cgo && !static_build
// +build linux,cgo,!static_build

package lxc

// #cgo CFLAGS: -std=gnu11 -Wvla -Werror
// #cgo pkg-config: lxc
import "C"
