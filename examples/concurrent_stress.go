/*
 * concurrent_stress.go
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
	"fmt"
	"github.com/caglar10ur/lxc"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var wg sync.WaitGroup

	nOfContainers := 10

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			name := strconv.Itoa(rand.Intn(nOfContainers))

			c := lxc.NewContainer(name)
			defer lxc.PutContainer(c)

			// sleep for a while to simulate some dummy work
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))

			if c.Defined() {
				if !c.Running() {
					c.SetDaemonize()
					fmt.Printf("Starting the container (%s)...\n", name)
					if err := c.Start(false); err != nil {
						fmt.Printf("ERROR: %s\n", err.Error())
					}
				} else {
					fmt.Printf("Stopping the container (%s)...\n", name)
					if err := c.Stop(); err != nil {
						fmt.Printf("ERROR: %s\n", err.Error())
					}
				}
			} else {
				fmt.Printf("Creating the container (%s)...\n", name)
				if err := c.Create("ubuntu", "amd64", "quantal"); err != nil {
					fmt.Printf("ERROR: %s\n", err.Error())
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	for i := 0; i < nOfContainers; i++ {
		name := strconv.Itoa(i)

		c := lxc.NewContainer(name)
		defer lxc.PutContainer(c)

		c.Stop()

		fmt.Printf("Destroying the container (%s)...\n", name)
		if err := c.Destroy(); err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
		}
	}
}
