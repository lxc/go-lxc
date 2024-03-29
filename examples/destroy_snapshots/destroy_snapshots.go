// Copyright © 2013, 2014, The Go-LXC Authors. All rights reserved.
// Use of this source code is governed by a LGPLv2.1
// license that can be found in the LICENSE file.

//go:build linux && cgo
// +build linux,cgo

package main

import (
	"log"

	"github.com/lxc/go-lxc"
)

func main() {
	c := lxc.Containers()
	for i := range c {
		log.Printf("%s\n", c[i].Name())
		l, err := c[i].Snapshots()
		if err != nil {
			log.Fatalf("ERROR: %s\n", err.Error())
		}

		for _, s := range l {
			log.Printf("Destroying Snaphot: %s\n", s.Name)
			if err := c[i].DestroySnapshot(s); err != nil {
				log.Fatalf("ERROR: %s\n", err.Error())
			}
		}
		c[i].Release()
	}
}
