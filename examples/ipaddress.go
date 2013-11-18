/*
 * ipaddress.go
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

	log.Printf("IPAddress(\"lo\")\n")
	if addresses, err := c.IPAddress("lo"); err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	} else {
		for i, v := range addresses {
			log.Printf("%d) %s\n", i, v)
		}
	}

	log.Printf("IPAddresses()\n")
	if addresses, err := c.IPAddresses(); err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	} else {
		for i, v := range addresses {
			log.Printf("%d) %s\n", i, v)
		}
	}
	log.Printf("IPv4Addresses()\n")
	if addresses, err := c.IPv4Addresses(); err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	} else {
		for i, v := range addresses {
			log.Printf("%d) %s\n", i, v)
		}
	}
	log.Printf("IPv6Addresses()\n")
	if addresses, err := c.IPv6Addresses(); err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	} else {
		for i, v := range addresses {
			log.Printf("%d) %s\n", i, v)
		}
	}
}
