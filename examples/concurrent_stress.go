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
	"flag"
	"runtime"
	"strconv"
	"sync"

	logger "github.com/caglar10ur/gologger"
	"github.com/lxc/go-lxc"
)

var (
	lxcpath       string
	iteration     int
	count         int
	template      string
	debug         bool
	startstop     bool
	createdestroy bool
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.StringVar(&lxcpath, "lxcpath", lxc.DefaultConfigPath(), "Use specified container path")
	flag.StringVar(&template, "template", "busybox", "Template to use")
	flag.IntVar(&count, "count", 10, "Number of operations to run concurrently")
	flag.IntVar(&iteration, "iteration", 1, "Number times to run the test")
	flag.BoolVar(&debug, "debug", false, "Flag to control debug output")
	flag.BoolVar(&startstop, "startstop", false, "Flag to execute Start and Stop")
	flag.BoolVar(&createdestroy, "createdestroy", false, "Flag to execute Create and Destroy")
	flag.Parse()
}

func main() {
	log := logger.New(nil)
	if debug {
		log.SetLogLevel(logger.Debug)
	}

	log.Debugf("Using %d GOMAXPROCS", runtime.NumCPU())

	var wg sync.WaitGroup

	for i := 0; i < iteration; i++ {
		log.Debugf("-- ITERATION %d --", i+1)
		for _, mode := range []string{"CREATE", "START", "STOP", "DESTROY"} {
			log.Debugf("\t-- %s --", mode)
			for j := 0; j < count; j++ {
				wg.Add(1)
				go func(i int, mode string) {
					c, err := lxc.NewContainer(strconv.Itoa(i), lxcpath)
					if err != nil {
						log.Fatalf("ERROR: %s\n", err.Error())
					}
					defer lxc.PutContainer(c)

					if mode == "CREATE" && startstop == false {
						log.Debugf("\t\tCreating the container (%d)...\n", i)
						if err := c.Create(template); err != nil {
							log.Errorf("\t\t\tERROR: %s\n", err.Error())
						}
					} else if mode == "START" && createdestroy == false {
						log.Debugf("\t\tStarting the container (%d)...\n", i)
						if err := c.Start(); err != nil {
							log.Errorf("\t\t\tERROR: %s\n", err.Error())
						}
					} else if mode == "STOP" && createdestroy == false {
						log.Debugf("\t\tStoping the container (%d)...\n", i)
						if err := c.Stop(); err != nil {
							log.Errorf("\t\t\tERROR: %s\n", err.Error())
						}
					} else if mode == "DESTROY" && startstop == false {
						log.Debugf("\t\tDestroying the container (%d)...\n", i)
						if err := c.Destroy(); err != nil {
							log.Errorf("\t\t\tERROR: %s\n", err.Error())
						}
					}
					wg.Done()
				}(j, mode)
			}
			wg.Wait()
		}
	}
}
