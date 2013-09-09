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
	"fmt"
	"github.com/caglar10ur/lxc"
)

var (
	name string
)

func init() {
	flag.StringVar(&name, "name", "rubik", "Name of the original container")
	flag.Parse()
}

func main() {
	c := lxc.NewContainer(name)
	defer lxc.PutContainer(c)

	directoryClone := name + "Directory"
	overlayClone := name + "OverlayFS"

	fmt.Printf("Cloning the container using Directory backend...\n")
	if err := c.CloneToDirectory(directoryClone); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}

	fmt.Printf("Cloning the container using OverlayFS backend...\n")
	if err := c.CloneToOverlayFS(overlayClone); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}
