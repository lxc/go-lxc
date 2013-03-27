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
	//  "strings"
	"testing"
)

func TestCreate(t *testing.T) {
	z := NewContainer("rubik")
	z.SetDaemonize()
	fmt.Printf("Creating container...\n")
	fmt.Printf("Create: %+v\n", z.Create("ubuntu", []string{"amd64", "quantal"}))

	if !z.Defined() {
		t.Errorf("Creating a container failed...")
	}
}

func TestStart(t *testing.T) {
	z := NewContainer("rubik")
	z.SetDaemonize()
	fmt.Printf("Starting container...\n")
	fmt.Printf("Start: %+v\n", z.Start(false, nil))

	if !z.Running() {
		t.Errorf("Starting a container failed...")
	}
}

func TestFreeze(t *testing.T) {
	z := NewContainer("rubik")
	z.SetDaemonize()
	fmt.Printf("Freezing container...\n")
	fmt.Printf("Freeze: %+v\n", z.Freeze())

	if z.GetState() != "FROZEN" {
		t.Errorf("Freezing a container failed...")
	}
}

func TestUnfreeze(t *testing.T) {
	z := NewContainer("rubik")
	z.SetDaemonize()
	fmt.Printf("Unfreezing container...\n")
	fmt.Printf("Unfreeze: %+v\n", z.Unfreeze())

	if z.GetState() != "RUNNING" {
		t.Errorf("Unfreezing a container failed...")
	}
}

/*
func TestAll(t *testing.T) {
    z := NewContainer("rubik")

    fmt.Printf("Container name: %s\n", z.GetName())
    fmt.Printf("Config file: %+v\n", z.GetConfigFileName())
    fmt.Printf("Load Config File: %+v\n", z.LoadConfigFile("/var/lib/lxc/rubik/config"))
    fmt.Printf("Save Config File: %+v\n", z.SaveConfigFile("config"))
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

    utsname_key := "lxc.utsname"
    utsname_value := z.GetConfigItem(utsname_key)[0]
    fmt.Printf("GetConfigItem: %s\n", utsname_value)
    fmt.Printf("SetConfigItem: %+v\n", z.SetConfigItem(utsname_key, "kibur"))
    fmt.Printf("GetConfigItem: %s\n", z.GetConfigItem(utsname_key))
    fmt.Printf("SetConfigItem: %+v\n", z.SetConfigItem(utsname_key, utsname_value))

    fmt.Printf("GetConfigItem: %s\n", z.GetConfigItem("lxc.arch"))
    fmt.Printf("GetConfigItem: %s\n", z.GetConfigItem("lxc.mount"))

    caps_key := "lxc.cap.drop"
    caps_value := strings.Join(z.GetConfigItem(caps_key), " ")
    fmt.Printf("GetConfigItem: %s\n", caps_value)
    fmt.Printf("ClearConfigItem: %+v\n", z.ClearConfigItem(caps_key))
    fmt.Printf("GetConfigItem: %s\n", z.GetConfigItem(caps_key))
    fmt.Printf("SetConfigItem: %+v\n", z.SetConfigItem(caps_key, caps_value))
    fmt.Printf("GetConfigItem: %s\n", z.GetConfigItem(caps_key))

    fmt.Printf("GetKeys: %s\n", z.GetKeys("lxc.network.0"))
    fmt.Printf("Wait 5 sec. (RUNNING): %+v\n", z.Wait(RUNNING, 5))

    if z.Running() {
        fmt.Printf("Shutting down container...\n\n")
        fmt.Printf("Shutdown: %+v\n", z.Shutdown(30))
        fmt.Printf("State: %+v\n", z.GetState())
        fmt.Printf("Stop: %+v\n", z.Stop())
        fmt.Printf("State: %+v\n", z.GetState())
        fmt.Printf("Destroy: %+v\n", z.Destroy())
    }
}
*/

func TestShutdown(t *testing.T) {
	z := NewContainer("rubik")
	fmt.Printf("Shutting down a container...\n")
	fmt.Printf("Shut down: %+v\n", z.Shutdown(30))
	if z.Running() {
		t.Errorf("Shutting down a container failed...")
	}
}

func TestDestroy(t *testing.T) {
	z := NewContainer("rubik")
	fmt.Printf("Destroying container...\n")
	fmt.Printf("Destroy: %+v\n", z.Destroy())
	if z.Defined() {
		t.Errorf("Destroying a container failed...")
	}
}
