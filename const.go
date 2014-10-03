// Copyright © 2013, 2014, The Go-LXC Authors. All rights reserved.
// Use of this source code is governed by a LGPLv2.1
// license that can be found in the LICENSE file.

// +build linux

package lxc

const (
	// WaitForever timeout
	WaitForever int = iota - 1
	// DontWait timeout
	DontWait
)

const (
	isDefined = 1 << iota
	isNotDefined
	isRunning
	isNotRunning
)
