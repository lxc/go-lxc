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
	"strings"
	"testing"
)

const (
	CONTAINER_NAME   = "rubik"
	CONFIG_FILE_NAME = "/var/lib/lxc/rubik/config"
)

func TestDefined_Negative(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	if z.Defined() {
		t.Errorf("Defined_Negative failed...")
	}
}

func TestCreate(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	fmt.Printf("Creating the container...\n")
	z.Create("ubuntu", []string{"amd64", "quantal"})

	if !z.Defined() {
		t.Errorf("Creating the container failed...")
	}
}

func TestGetConfigFileName(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	if z.GetConfigFileName() != CONFIG_FILE_NAME {
		t.Errorf("GetConfigFileName failed...")
	}
}

func TestDefined_Positive(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	if !z.Defined() {
		t.Errorf("Defined failed...")
	}
}

func TestInitPID_Negative(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	if z.GetInitPID() != -1 {
		t.Errorf("GetInitPID failed...")
	}
}

func TestStart(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	fmt.Printf("Starting the container...\n")
	z.SetDaemonize()
	z.Start(false, nil)

	z.Wait(RUNNING, 5)
	if !z.Running() {
		t.Errorf("Starting the container failed...")
	}
}

func TestSetDaemonize(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	z.SetDaemonize()
	if !z.GetDaemonize() {
		t.Errorf("GetDaemonize failed...")
	}
}

func TestInitPID_Positive(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	if z.GetInitPID() == -1 {
		t.Errorf("GetInitPID failed...")
	}
}

func TestGetName(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	if z.GetName() != CONTAINER_NAME {
		t.Errorf("GetName failed...")
	}
}

func TestFreeze(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	fmt.Printf("Freezing the container...\n")
	z.Freeze()

	z.Wait(FROZEN, 5)
	if z.GetState() != FROZEN {
		t.Errorf("Freezing the container failed...")
	}
}

func TestUnfreeze(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	fmt.Printf("Unfreezing the container...\n")
	z.Unfreeze()

	z.Wait(RUNNING, 5)
	if z.GetState() != RUNNING {
		t.Errorf("Unfreezing the container failed...")
	}
}

func TestLoadConfigFile(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	if !z.LoadConfigFile(CONFIG_FILE_NAME) {
		t.Errorf("LoadConfigFile failed...")
	}
}

func TestSaveConfigFile(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	if !z.SaveConfigFile(CONFIG_FILE_NAME) {
		t.Errorf("LoadConfigFile failed...")
	}
}

func TestGetConfigItem(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	if z.GetConfigItem("lxc.utsname")[0] != CONTAINER_NAME {
		t.Errorf("GetConfigItem failed...")
	}
}

func TestSetConfigItem(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	z.SetConfigItem("lxc.utsname", CONTAINER_NAME)
	if z.GetConfigItem("lxc.utsname")[0] != CONTAINER_NAME {
		t.Errorf("GetConfigItem failed...")
	}
}

func TestClearConfigItem(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	z.ClearConfigItem("lxc.cap.drop")
	if z.GetConfigItem("lxc.cap.drop")[0] != "" {
		t.Errorf("ClearConfigItem failed...")
	}
}

func TestGetKeys(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	keys := strings.Join(z.GetKeys("lxc.network.0"), " ")
	if !strings.Contains(keys, "mtu") {
		t.Errorf("GetKeys failed...")
	}
}

func TestGetNumberOfNetworkInterfaces(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	if z.GetNumberOfNetworkInterfaces() != 1 {
		t.Errorf("GetNumberOfNetworkInterfaces failed...")
	}
}

func TestShutdown(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	fmt.Printf("Shutting down the container...\n")
	z.Shutdown(30)

	if z.Running() {
		t.Errorf("Shutting down the container failed...")
	}
}

func TestStop(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	fmt.Printf("Stopping the container...\n")
	z.Stop()
	if z.Running() {
		t.Errorf("Stopping the container failed...")
	}
}

func TestDestroy(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	fmt.Printf("Destroying the container...\n")
	z.Destroy()

	if z.Defined() {
		t.Errorf("Destroying the container failed...")
	}
}
