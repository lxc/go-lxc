/*
 * limit.go
 *
 * Copyright © 2013, 2014, S.Çağlar Onur
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

	"gopkg.in/lxc/go-lxc.v1"
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
	defer lxc.PutContainer(c)

	memLimit, err := c.MemoryLimit()
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
	memorySwapLimit, err := c.MemorySwapLimit()
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}

	if err := c.SetMemoryLimit(memLimit / 4); err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
	if err := c.SetMemorySwapLimit(memorySwapLimit / 4); err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
}
