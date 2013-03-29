/*
 * example.go
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
)

func main() {
	z := lxc.NewContainer("rubik")

	fmt.Printf("Container name: %s\n", lxc.GetVersion())
	fmt.Printf("Container name: %s\n", z.GetName())
	fmt.Printf("Config file: %+v\n", z.GetConfigFileName())
	fmt.Printf("Daemonize: %+v\n", z.GetDaemonize())
	fmt.Printf("Init PID: %+v\n", z.GetInitPID())
	fmt.Printf("Defined: %+v\n", z.Defined())
	fmt.Printf("Running: %+v\n", z.Running())
	fmt.Printf("State: %+v\n", z.GetState())
	z.SetDaemonize()
	fmt.Printf("Daemonize: %+v\n", z.GetDaemonize())

	if !z.Defined() {
		fmt.Printf("Creating container...\n")
		fmt.Printf("Create: %+v\n", z.Create("ubuntu", []string{"amd64", "quantal"}))
	} else {
		fmt.Printf("Starting container...\n\n")
		fmt.Printf("Start: %+v\n", z.Start(false, nil))
		fmt.Printf("State: %+v\n", z.GetState())
		fmt.Printf("Init PID: %+v\n", z.GetInitPID())
		fmt.Printf("Freeze: %+v\n", z.Freeze())
		fmt.Printf("State: %+v\n", z.GetState())
		fmt.Printf("Unfreeze: %+v\n", z.Unfreeze())
		fmt.Printf("State: %+v\n", z.GetState())
	}

	fmt.Printf("GetKeys: %s\n", z.GetKeys("lxc.network.0"))
	fmt.Printf("Wait 5 sec. (lxc.RUNNING): %+v\n", z.Wait(lxc.RUNNING, 5))

	if z.Running() {
		fmt.Printf("Shutting down container...\n\n")
		fmt.Printf("Shutdown: %+v\n", z.Shutdown(30))
		fmt.Printf("State: %+v\n", z.GetState())
		fmt.Printf("Stop: %+v\n", z.Stop())
		fmt.Printf("State: %+v\n", z.GetState())
		fmt.Printf("Destroy: %+v\n", z.Destroy())
	}
}
