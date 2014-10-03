// Copyright © 2013, 2014, The Go-LXC Authors. All rights reserved.
// Use of this source code is governed by a LGPLv2.1
// license that can be found in the LICENSE file.

// +build linux

package lxc

type AttachOptions struct {
	// Stdinfd specifies the fd to read input from
	Stdinfd uintptr

	// Stdoutfd specifies the fd to write output to
	Stdoutfd uintptr

	// Stderrfd specifies the fd to write error output to
	Stderrfd uintptr

	// Env specifies the environment of the process.
	Env []string

	// Cwd specifies the working directory of the command.
	Cwd string

	// If ClearEnv is true the environment is cleared before
	// running the command.
	ClearEnv bool
}
