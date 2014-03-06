// Copyright © 2013, 2014, S.Çağlar Onur
// Use of this source code is governed by a LGPLv2.1
// license that can be found in the LICENSE file.
//
// Authors:
// S.Çağlar Onur <caglar@10ur.org>

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
