// Copyright © 2013, 2014, The Go-LXC Authors. All rights reserved.
// Use of this source code is governed by a LGPLv2.1
// license that can be found in the LICENSE file.

//go:build linux && cgo
// +build linux,cgo

package main

import (
	"flag"
	"log"
	"time"

	"github.com/lxc/go-lxc"
)

var (
	lxcpath string
	name    string
)

func init() {
	flag.StringVar(&lxcpath, "lxcpath", lxc.DefaultConfigPath(), "Use specified container path")
	flag.StringVar(&name, "name", "rubik", "Name of the container")
	flag.Parse()
}

func main() {
	c, err := lxc.NewContainer(name, lxcpath)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
	defer c.Release()

	log.Printf("Freezing the container...\n")
	if err := c.Freeze(); err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
	c.Wait(lxc.FROZEN, 10*time.Second)
}
