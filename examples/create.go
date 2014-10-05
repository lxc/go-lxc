// Copyright © 2013, 2014, The Go-LXC Authors. All rights reserved.
// Use of this source code is governed by a LGPLv2.1
// license that can be found in the LICENSE file.

// +build linux

package main

import (
	"flag"
	"log"
	"os"

	"gopkg.in/lxc/go-lxc.v2"
)

var (
	lxcpath  string
	template string
	distro   string
	release  string
	arch     string
	name     string
	verbose  bool
)

func init() {
	flag.StringVar(&lxcpath, "lxcpath", lxc.DefaultConfigPath(), "Use specified container path")
	flag.StringVar(&template, "template", "ubuntu", "Template to use")
	flag.StringVar(&distro, "distro", "ubuntu", "Template to use")
	flag.StringVar(&release, "release", "trusty", "Template to use")
	flag.StringVar(&arch, "arch", "amd64", "Template to use")
	flag.StringVar(&name, "name", "rubik", "Name of the container")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")
	flag.Parse()
}

func main() {
	c, err := lxc.NewContainer(name, lxcpath)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
	defer lxc.Release(c)

	log.Printf("Creating container...\n")
	if verbose {
		c.SetVerbosity(lxc.Verbose)
	}

	options := lxc.BusyboxTemplateOptions
	if os.Geteuid() != 0 {
		options = lxc.DownloadTemplateOptions
		options.Release = release
		options.Distro = distro
		options.Arch = arch
	} else {
		options = lxc.UbuntuTemplateOptions
		options.Template = template
		options.Release = release
		options.Arch = arch
	}

	if err := c.Create(options); err != nil {
		log.Printf("ERROR: %s\n", err.Error())
	}
}
