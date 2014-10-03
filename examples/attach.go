// Copyright © 2013, 2014, The Go-LXC Authors. All rights reserved.
// Use of this source code is governed by a LGPLv2.1
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"os"

	"gopkg.in/lxc/go-lxc.v2"
)

var (
	lxcpath string
	name    string
	clear   bool
)

func init() {
	flag.StringVar(&lxcpath, "lxcpath", lxc.DefaultConfigPath(), "Use specified container path")
	flag.StringVar(&name, "name", "rubik", "Name of the original container")
	flag.BoolVar(&clear, "clear", false, "Attach with clear environment")
	flag.Parse()
}

func main() {
	c, err := lxc.NewContainer(name, lxcpath)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
	defer lxc.Release(c)

	if clear {
		log.Printf("AttachShellWithClearEnvironment\n")
		if err := c.AttachShellWithClearEnvironment(); err != nil {
			log.Fatalf("ERROR: %s\n", err.Error())
		}

		log.Printf("RunCommandWithClearEnvironment\n")
		_, err := c.RunCommand([]string{"uname", "-a"}, &lxc.AttachOptions{
			ClearEnv: true,
			Stdinfd:  os.Stdin.Fd(),
			Stdoutfd: os.Stdout.Fd(),
			Stderrfd: os.Stderr.Fd(),
		})
		if err != nil {
			log.Fatalf("ERROR: %s\n", err.Error())
		}

	} else {
		log.Printf("AttachShell\n")
		if err := c.AttachShell(); err != nil {
			log.Fatalf("ERROR: %s\n", err.Error())
		}

		log.Printf("RunCommand\n")
		_, err := c.RunCommand([]string{"uname", "-a"}, &lxc.AttachOptions{
			ClearEnv: false,
			Stdinfd:  os.Stdin.Fd(),
			Stdoutfd: os.Stdout.Fd(),
			Stderrfd: os.Stderr.Fd(),
		})
		if err != nil {
			log.Fatalf("ERROR: %s\n", err.Error())
		}
	}
}
