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
	lxcpath    string
	template   string
	distro     string
	release    string
	arch       string
	name       string
	verbose    bool
	flush      bool
	validation bool
	fssize     string
	bdevtype   string
)

func init() {
	flag.StringVar(&lxcpath, "lxcpath", lxc.DefaultConfigPath(), "Use specified container path")
	flag.StringVar(&template, "template", "download", "Template to use")
	flag.StringVar(&distro, "distro", "ubuntu", "Template to use")
	flag.StringVar(&release, "release", "trusty", "Template to use")
	flag.StringVar(&arch, "arch", "amd64", "Template to use")
	flag.StringVar(&name, "name", "rubik", "Name of the container")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")
	flag.BoolVar(&flush, "flush", false, "Flush the cache")
	flag.BoolVar(&validation, "validation", false, "GPG validation")

	flag.StringVar(&bdevtype, "bdev", "dir", "backing store type")
	flag.StringVar(&fssize, "fssize", "", "backing store size")
	// TODO support more flags for zfs, lvm, or rbd

	flag.Parse()
}

func main() {
	c, err := lxc.NewContainer(name, lxcpath)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
	defer c.Release()

	log.Printf("Creating container...\n")
	if verbose {
		c.SetVerbosity(lxc.Verbose)
	}

	var backend lxc.BackendStore
	if err := (&backend).Set(bdevtype); err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}

	var bdevSize lxc.ByteSize
	if fssize != "" {
		var err error
		bdevSize, err = lxc.ParseBytes(fssize)
		if err != nil {
			log.Fatalf("ERROR: %s\n", err.Error())
		}
	}

	options := lxc.TemplateOptions{
		Template:             template,
		Distro:               distro,
		Release:              release,
		Arch:                 arch,
		FlushCache:           flush,
		DisableGPGValidation: validation,
		Backend:              backend,
		BackendSpecs: &lxc.BackendStoreSpecs{
			FSSize: uint64(bdevSize),
		},
	}

	c.SetLogFile("log")
	c.SetLogLevel(lxc.DEBUG)

	if err := c.Create(options); err != nil {
		log.Printf("ERROR: %s\n", err.Error())
	}
}
