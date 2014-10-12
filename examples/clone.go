// Copyright © 2013, 2014, The Go-LXC Authors. All rights reserved.
// Use of this source code is governed by a LGPLv2.1
// license that can be found in the LICENSE file.

// +build linux,cgo

package main

import (
	"flag"
	"log"

	"gopkg.in/lxc/go-lxc.v2"
)

var (
	lxcpath string
	name    string
)

func init() {
	flag.StringVar(&lxcpath, "lxcpath", lxc.DefaultConfigPath(), "Use specified container path")
	flag.StringVar(&name, "name", "rubik", "Name of the original container")
	flag.Parse()
}

func main() {
	c, err := lxc.NewContainer(name, lxcpath)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}

	directoryClone := name + "Directory"
	overlayClone := name + "Overlayfs"
	btrfsClone := name + "Btrfs"
	aufsClone := name + "Aufs"

	log.Printf("Cloning the container using Directory backend...\n")
	if err := c.Clone(directoryClone); err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}

	log.Printf("Cloning the container using Overlayfs backend...\n")
	if err := c.CloneUsing(overlayClone, lxc.Overlayfs, lxc.CloneSnapshot); err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}

	log.Printf("Cloning the container using Aufs backend...\n")
	if err := c.CloneUsing(aufsClone, lxc.Aufs, lxc.CloneSnapshot); err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}

	log.Printf("Cloning the container using Btrfs backend...\n")
	if err := c.CloneUsing(btrfsClone, lxc.Btrfs, 0); err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
}
