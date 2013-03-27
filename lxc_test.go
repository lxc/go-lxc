/*
 * lxc_test.go: Go bindings for lxc
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

package lxc

import (
	"fmt"
	"testing"
)

func TestAll(t *testing.T) {
	z := NewContainer("rubik")

	fmt.Printf("Config file:%+v\n", z.ConfigFileName())
	fmt.Printf("Daemonize: %+v\n", z.Daemonize())
	fmt.Printf("Init PID: %+v\n", z.InitPID())
	fmt.Printf("Defined: %+v\n", z.Defined())
	fmt.Printf("Running: %+v\n", z.Running())
	fmt.Printf("State: %+v\n", z.State())
	z.SetDaemonize()
	fmt.Printf("Daemonize: %+v\n", z.Daemonize())

	if !z.Defined() {
		fmt.Printf("Creating rubik container...\n")
		fmt.Printf("Create: %+v\n", z.Create("ubuntu", []string{"amd64", "quantal"}))
	} else {
		fmt.Printf("Starting rubik container...\n\n")
		fmt.Printf("Start: %+v\n", z.Start(false, nil))
		fmt.Printf("State: %+v\n", z.State())
		fmt.Printf("Init PID: %+v\n", z.InitPID())
		fmt.Printf("Freeze: %+v\n", z.Freeze())
		fmt.Printf("State: %+v\n", z.State())
		fmt.Printf("Unfreeze: %+v\n", z.Unfreeze())
		fmt.Printf("State: %+v\n", z.State())
	}

	if z.Running() {
		fmt.Printf("Shutdown: %+v\n", z.Shutdown(30))
		fmt.Printf("State: %+v\n", z.State())
		fmt.Printf("Stop: %+v\n", z.Stop())
		fmt.Printf("State: %+v\n", z.State())
	}
	fmt.Printf("Destroy: %+v\n", z.Destroy())
}
