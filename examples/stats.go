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
	"fmt"
	"github.com/caglar10ur/lxc"
)

var (
	name string
)

func init() {
	flag.StringVar(&name, "name", "rubik", "Name of the container")
	flag.Parse()
}

func main() {
	c := lxc.NewContainer(name)
	defer lxc.PutContainer(c)

	// mem
	memUsed, err := c.MemoryUsage()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	} else {
		fmt.Printf("MemoryUsage: %s\n", memUsed)
	}

	memLimit, err := c.MemoryLimit()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	} else {
		fmt.Printf("MemoryLimit: %s\n", memLimit)
	}

	// swap
	swapUsed, err := c.SwapUsage()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	} else {
		fmt.Printf("SwapUsage: %s\n", swapUsed)
	}

	swapLimit, err := c.SwapLimit()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	} else {
		fmt.Printf("SwapLimit: %s\n", swapLimit)
	}

	cpuTime, err := c.CPUTime()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
	fmt.Printf("cpuacct.usage: %s\n", cpuTime)

	cpuTimePerCPU, err := c.CPUTimePerCPU()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
	fmt.Printf("cpuacct.usageerrpercpu: %s\n", cpuTimePerCPU)

	cpuStats, err := c.CPUStats()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
	fmt.Printf("cpuacct.stat: %v\n", cpuStats)
}
