// Copyright © 2013, 2014, The Go-LXC Authors. All rights reserved.
// Use of this source code is governed by a LGPLv2.1
// license that can be found in the LICENSE file.

// +build linux

package lxc

import (
	"os"
)

// AttachOptions type is used for defining various attach options
type AttachOptions struct {

	// Specify the namespaces to attach to, as OR'ed list of clone flags (syscall.CLONE_NEWNS | syscall.CLONE_NEWUTS ...)
	Namespaces int

	// Specify the architecture which the kernel should appear to be running as to the command executed.
	Arch Personality

	// Cwd specifies the working directory of the command.
	Cwd string

	// Uid specifies the user id to run as.
	Uid int

	// Gid specifies the group id to run as.
	Gid int

	// If ClearEnv is true the environment is cleared before running the command.
	ClearEnv bool

	// Env specifies the environment of the process.
	Env []string

	// EnvToKeep specifies the environment of the process when ClearEnv is true.
	EnvToKeep []string

	// Stdinfd specifies the fd to read input from.
	StdinFd uintptr

	// Stdoutfd specifies the fd to write output to.
	StdoutFd uintptr

	// Stderrfd specifies the fd to write error output to.
	StderrFd uintptr
}

// DefaultAttachOptions is a convenient set of options to be used
var DefaultAttachOptions = &AttachOptions{
	Namespaces: -1,
	Arch:       -1,
	Cwd:        "/",
	Uid:        -1,
	Gid:        -1,
	ClearEnv:   false,
	Env:        nil,
	EnvToKeep:  nil,
	StdinFd:    os.Stdin.Fd(),
	StdoutFd:   os.Stdout.Fd(),
	StderrFd:   os.Stderr.Fd(),
}
