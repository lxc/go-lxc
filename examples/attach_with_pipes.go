// Copyright © 2013, 2014, The Go-LXC Authors. All rights reserved.
// Use of this source code is governed by a LGPLv2.1
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"io"
	"log"
	"os"
	"sync"

	"gopkg.in/lxc/go-lxc.v2"
)

var (
	lxcpath string
	name    string
	clear   bool
	wg      sync.WaitGroup
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
	defer lxc.PutContainer(c)

	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
	stderrReader, stderrWriter, err := os.Pipe()
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := io.Copy(os.Stdout, stdoutReader)
		if err != nil {
			log.Fatalf("ERROR: %s\n", err.Error())
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err = io.Copy(os.Stderr, stderrReader)
		if err != nil {
			log.Fatalf("ERROR: %s\n", err.Error())
		}
	}()

	if clear {
		log.Printf("AttachShellWithClearEnvironment\n")
		if err := c.AttachShellWithClearEnvironment(); err != nil {
			log.Fatalf("ERROR: %s\n", err.Error())
		}

		log.Printf("RunCommandWithClearEnvironment\n")
		if _, err := c.RunCommandWithClearEnvironment(os.Stdin.Fd(), stdoutWriter.Fd(), stderrWriter.Fd(), "uname", "-a"); err != nil {
			log.Fatalf("ERROR: %s\n", err.Error())
		}
	} else {
		log.Printf("AttachShell\n")
		if err := c.AttachShell(); err != nil {
			log.Fatalf("ERROR: %s\n", err.Error())
		}

		log.Printf("RunCommand\n")
		if _, err := c.RunCommand(os.Stdin.Fd(), stdoutWriter.Fd(), stderrWriter.Fd(), "uname", "-a"); err != nil {
			log.Fatalf("ERROR: %s\n", err.Error())
		}
	}

	if err := stdoutWriter.Close(); err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
	if err := stderrWriter.Close(); err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}

	wg.Wait()
}
