/*
 * concurrent_start.go
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
	"runtime"
	"strconv"
	"sync"
	"github.com/caglar10ur/lxc"
)

var (
	lxcpath string
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&lxcpath, "lxcpath", lxc.DefaultConfigPath(), "Use specified container path")
	flag.Parse()
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			c, err := lxc.NewContainer(strconv.Itoa(i), lxcpath)
			if err != nil {
				log.Fatalf("ERROR: %s\n", err.Error())
			}
			defer lxc.PutContainer(c)

			c.SetDaemonize()
			log.Printf("Starting the container (%d)...\n", i)
			if err := c.Start(false); err != nil {
				log.Fatalf("ERROR: %s\n", err.Error())
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
