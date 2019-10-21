// Copyright © 2013, 2014, The Go-LXC Authors. All rights reserved.
// Use of this source code is governed by a LGPLv2.1
// license that can be found in the LICENSE file.

// +build >go1.10,linux,cgo

package lxc

// #cgo CFLAGS: -fvisibility=hidden
import "C"
