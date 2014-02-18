/*
 * stats.go
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
	defer lxc.PutContainer(c)

	// mem
	memUsed, err := c.MemoryUsage()
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	} else {
		log.Printf("MemoryUsage: %s\n", memUsed)
	}

	memLimit, err := c.MemoryLimit()
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	} else {
		log.Printf("MemoryLimit: %s\n", memLimit)
	}

	// kmem
	kmemUsed, err := c.KernelMemoryUsage()
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	} else {
		log.Printf("KernelMemoryUsage: %s\n", kmemUsed)
	}

	kmemLimit, err := c.KernelMemoryLimit()
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	} else {
		log.Printf("KernelMemoryLimit: %s\n", kmemLimit)
	}

	// swap
	swapUsed, err := c.MemorySwapUsage()
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	} else {
		log.Printf("MemorySwapUsage: %s\n", swapUsed)
	}

	swapLimit, err := c.MemorySwapLimit()
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	} else {
		log.Printf("MemorySwapLimit: %s\n", swapLimit)
	}

	// blkio
	blkioUsage, err := c.BlkioUsage()
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	} else {
		log.Printf("BlkioUsage: %s\n", blkioUsage)
	}

	cpuTime, err := c.CPUTime()
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
	log.Printf("cpuacct.usage: %s\n", cpuTime)

	cpuTimePerCPU, err := c.CPUTimePerCPU()
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
	log.Printf("cpuacct.usageerrpercpu: %v\n", cpuTimePerCPU)

	cpuStats, err := c.CPUStats()
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
	log.Printf("cpuacct.stat: %v\n", cpuStats)

	interfaceStats, err := c.InterfaceStats()
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
	log.Printf("InterfaceStats: %v\n", interfaceStats)
}
