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
	"strings"
	"testing"
)

const (
	CONTAINER_NAME   = "rubik"
	CONFIG_FILE_PATH = "/var/lib/lxc"
	CONFIG_FILE_NAME = "/var/lib/lxc/rubik/config"
)

func TestGetVersion(t *testing.T) {

	if GetVersion() == "" {
		t.Errorf("GetVersion failed...")
	}
}

func TestGetDefaultConfigPath(t *testing.T) {
	if GetDefaultConfigPath() != CONFIG_FILE_PATH {
		t.Errorf("GetDefaultConfigPath failed...")
	}
}

func TestGetSetConfigPath(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	current_path := z.GetConfigPath()
	z.SetConfigPath("/tmp")
	new_path := z.GetConfigPath()

	if current_path == new_path {
		t.Errorf("GetSetConfigPath failed...")
	}
}

func TestDefined_Negative(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	if z.Defined() {
		t.Errorf("Defined_Negative failed...")
	}
}

func TestCreate(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	t.Logf("Creating the container...\n")
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

	t.Logf("Starting the container...\n")
	z.SetDaemonize()
	z.Start(false, nil)

	z.Wait(RUNNING, 30)
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

	t.Logf("Freezing the container...\n")
	z.Freeze()

	z.Wait(FROZEN, 30)
	if z.GetState() != FROZEN {
		t.Errorf("Freezing the container failed...")
	}
}

func TestUnfreeze(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	t.Logf("Unfreezing the container...\n")
	z.Unfreeze()

	z.Wait(RUNNING, 30)
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

func TestGetSetCgroupItem(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	max_mem := z.GetCgroupItem("memory.max_usage_in_bytes")[0]
	current_mem := z.GetCgroupItem("memory.limit_in_bytes")[0]
	z.SetCgroupItem("memory.limit_in_bytes", max_mem)
	new_mem := z.GetCgroupItem("memory.limit_in_bytes")[0]

	if new_mem == current_mem {
		t.Errorf("GetSetCgroupItem failed...")
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

func TestGetMemoryUsageInBytes(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	mem_used, _ := z.GetMemoryUsageInBytes()
	swap_used, _ := z.GetSwapUsageInBytes()
	mem_limit, _ := z.GetMemoryLimitInBytes()
	swap_limit, _ := z.GetSwapLimitInBytes()

	t.Logf("Mem usage: %0.0f\n", mem_used)
	t.Logf("Mem usage: %s\n", mem_used)
	t.Logf("Swap usage: %0.0f\n", swap_used)
	t.Logf("Swap usage: %s\n", swap_used)
	t.Logf("Mem limit: %0.0f\n", mem_limit)
	t.Logf("Mem limit: %s\n", mem_limit)
	t.Logf("Swap limit: %0.0f\n", swap_limit)
	t.Logf("Swap limit: %s\n", swap_limit)
}

func TestShutdown(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	t.Logf("Shutting down the container...\n")
	z.Shutdown(30)

	if z.Running() {
		t.Errorf("Shutting down the container failed...")
	}
}

func TestStop(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	t.Logf("Stopping the container...\n")
	z.Stop()
	if z.Running() {
		t.Errorf("Stopping the container failed...")
	}
}

func TestDestroy(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)

	t.Logf("Destroying the container...\n")
	z.Destroy()

	if z.Defined() {
		t.Errorf("Destroying the container failed...")
	}
}
