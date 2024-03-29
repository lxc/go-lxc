// Copyright © 2013, 2014, The Go-LXC Authors. All rights reserved.
// Use of this source code is governed by a LGPLv2.1
// license that can be found in the LICENSE file.

//go:build linux && cgo
// +build linux,cgo

package main

import (
	"flag"
	"log"

	"github.com/lxc/go-lxc"
)

var (
	lxcpath string
	name    string
	backend lxc.BackendStore
)

func init() {
	flag.StringVar(&lxcpath, "lxcpath", lxc.DefaultConfigPath(), "Use specified container path")
	flag.StringVar(&name, "name", "rubik", "Name of the original container")
	flag.Var(&backend, "backend", "Backend type to use, possible values are [dir, zfs, btrfs, lvm, aufs, overlayfs, loopback, best]")
	flag.Parse()
}

func main() {
	c, err := lxc.NewContainer(name, lxcpath)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
	defer c.Release()

	if backend == 0 {
		log.Fatalf("ERROR: %s\n", lxc.ErrUnknownBackendStore)
	}

	log.Printf("Cloning the container using %s backend...\n", backend)
	err = c.Clone(name+"_"+backend.String(), lxc.CloneOptions{
		Backend: backend,
	})
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
}
