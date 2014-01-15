/*
 * clone.go
 *
 * Copyright © 2013, S.Çağlar Onur
 *
 * Authors:
 * S.Çağlar Onur <caglar@10ur.org>
 *
 * This library is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2, as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License along
 * with this program; if not, write to the Free Software Foundation, Inc.,
 * 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
 */

package main

import (
	"flag"
	"log"

	"github.com/caglar10ur/lxc"
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
	defer lxc.PutContainer(c)

	directoryClone := name + "Directory"
	overlayClone := name + "Overlayfs"
	btrfsClone := name + "Btrfs"

	log.Printf("Cloning the container using Directory backend...\n")
	if err := c.Clone(directoryClone); err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}

	log.Printf("Cloning the container using Overlayfs backend...\n")
	if err := c.CloneUsing(overlayClone, lxc.Overlayfs, lxc.CloneSnapshot); err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}

	log.Printf("Cloning the container using Btrfs backend...\n")
	if err := c.CloneUsing(btrfsClone, lxc.Btrfs, 0); err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
}
