/*
 * lxc_test.go: Go bindings for lxc
 *
 * Copyright © 2013, S.Çağlar Onur
 *
 * Authors:
 * S.Çağlar Onur <caglar@10ur.org>
 *
 * This library is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 2.1 of the License, or (at your option) any later version.

 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.

 * You should have received a copy of the GNU Lesser General Public
 * License along with this library; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301  USA
 */
package lxc_test

import (
	"github.com/caglar10ur/lxc"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	ContainerName             = "rubik"
	CloneContainerName        = "O"
	CloneOverlayContainerName = "O_o"
	ConfigFilePath            = "/var/lib/lxc"
	ConfigFileName            = "/var/lib/lxc/rubik/config"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func TestVersion(t *testing.T) {
	t.Logf("LXC version: %s", lxc.Version())
}

func TestDefaultConfigPath(t *testing.T) {
	if lxc.DefaultConfigPath() != ConfigFilePath {
		t.Errorf("DefaultConfigPath failed...")
	}
}

func TestSetConfigPath(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	currentPath := z.ConfigPath()
	if err := z.SetConfigPath("/tmp"); err != nil {
		t.Errorf(err.Error())
	}
	newPath := z.ConfigPath()

	if currentPath == newPath {
		t.Errorf("SetConfigPath failed...")
	}
}

func TestConcurrentDefined_Negative(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i <= 100; i++ {
		wg.Add(1)
		go func() {
			z := lxc.NewContainer(strconv.Itoa(rand.Intn(10)))
			defer lxc.PutContainer(z)

			// sleep for a while to simulate some dummy work
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(250)))

			if z.Defined() {
				t.Errorf("Defined_Negative failed...")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestDefined_Negative(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	if z.Defined() {
		t.Errorf("Defined_Negative failed...")
	}
}

func TestCreate(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	if err := z.Create("ubuntu", "amd64", "quantal"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestClone(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	if err := z.CloneToDirectory(CloneContainerName); err != nil {
		t.Errorf(err.Error())
	}

	if err := z.CloneToOverlayFS(CloneOverlayContainerName); err != nil {
		t.Errorf(err.Error())
	}
}

func TestConcurrentCreate(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			z := lxc.NewContainer(strconv.Itoa(i))
			defer lxc.PutContainer(z)

			// sleep for a while to simulate some dummy work
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(250)))

			if err := z.Create("ubuntu", "amd64", "quantal"); err != nil {
				t.Errorf(err.Error())
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestContainerNames(t *testing.T) {
	t.Logf("Containers: %+v\n", lxc.ContainerNames())
}

func TestActiveContainerNames(t *testing.T) {
	t.Logf("Active Containers: %+v\n", lxc.ActiveContainerNames())
}

func TestContainers(t *testing.T) {
	for _, v := range lxc.Containers() {
		t.Logf("%s: %s", v.Name(), v.State())
	}
}

func TestActiveContainers(t *testing.T) {
	for _, v := range lxc.ActiveContainers() {
		t.Logf("%s: %s", v.Name(), v.State())
	}
}

func TestConcurrentStart(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			z := lxc.NewContainer(strconv.Itoa(i))
			defer lxc.PutContainer(z)

			z.SetDaemonize()
			if err := z.Start(false); err != nil {
				t.Errorf(err.Error())
			}
			z.Wait(lxc.RUNNING, 30)
			if !z.Running() {
				t.Errorf("Starting the container failed...")
			}

			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestConfigFileName(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)
	if z.ConfigFileName() != ConfigFileName {
		t.Errorf("ConfigFileName failed...")
	}
}

func TestDefined_Positive(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	if !z.Defined() {
		t.Errorf("Defined_Positive failed...")
	}
}

func TestConcurrentDefined_Positive(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i <= 100; i++ {
		wg.Add(1)
		go func() {
			z := lxc.NewContainer(strconv.Itoa(rand.Intn(10)))
			defer lxc.PutContainer(z)

			// sleep for a while to simulate some dummy work
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(250)))

			if !z.Defined() {
				t.Errorf("Defined_Positive failed...")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestInitPID_Negative(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	if z.InitPID() != -1 {
		t.Errorf("InitPID failed...")
	}
}

func TestStart(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	z.SetDaemonize()
	z.Start(false)

	z.Wait(lxc.RUNNING, 30)
}

func TestMayControl(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	if !z.MayControl() {
		t.Errorf("Controling the container failed...")
	}
}

func TestRunning(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	if !z.Running() {
		t.Errorf("Checking the container failed...")
	}
}
func TestSetDaemonize(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	z.SetDaemonize()
	if !z.Daemonize() {
		t.Errorf("Daemonize failed...")
	}
}

func TestInitPID_Positive(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	if z.InitPID() == -1 {
		t.Errorf("InitPID failed...")
	}
}

func TestName(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	if z.Name() != ContainerName {
		t.Errorf("Name failed...")
	}
}

func TestFreeze(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	if err := z.Freeze(); err != nil {
		t.Errorf(err.Error())
	}

	z.Wait(lxc.FROZEN, 30)
	if z.State() != lxc.FROZEN {
		t.Errorf("Freezing the container failed...")
	}
}

func TestUnfreeze(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	if err := z.Unfreeze(); err != nil {
		t.Errorf(err.Error())
	}

	z.Wait(lxc.RUNNING, 30)
	if z.State() != lxc.RUNNING {
		t.Errorf("Unfreezing the container failed...")
	}
}

func TestLoadConfigFile(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	if err := z.LoadConfigFile(ConfigFileName); err != nil {
		t.Errorf(err.Error())
	}
}

func TestSaveConfigFile(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	if err := z.SaveConfigFile(ConfigFileName); err != nil {
		t.Errorf(err.Error())
	}
}

func TestConfigItem(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	if z.ConfigItem("lxc.utsname")[0] != ContainerName {
		t.Errorf("ConfigItem failed...")
	}
}

func TestSetConfigItem(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	if err := z.SetConfigItem("lxc.utsname", ContainerName); err != nil {
		t.Errorf(err.Error())
	}
	if z.ConfigItem("lxc.utsname")[0] != ContainerName {
		t.Errorf("ConfigItem failed...")
	}
}

func TestSetCgroupItem(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	maxMem := z.CgroupItem("memory.max_usage_in_bytes")[0]
	currentMem := z.CgroupItem("memory.limit_in_bytes")[0]
	if err := z.SetCgroupItem("memory.limit_in_bytes", maxMem); err != nil {
		t.Errorf(err.Error())
	}
	newMem := z.CgroupItem("memory.limit_in_bytes")[0]

	if newMem == currentMem {
		t.Errorf("SetCgroupItem failed...")
	}
}

func TestClearConfigItem(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	if err := z.ClearConfigItem("lxc.cap.drop"); err != nil {
		t.Errorf(err.Error())
	}
	if z.ConfigItem("lxc.cap.drop")[0] != "" {
		t.Errorf("ClearConfigItem failed...")
	}
}

func TestKeys(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	keys := strings.Join(z.Keys("lxc.network.0"), " ")
	if !strings.Contains(keys, "mtu") {
		t.Errorf("Keys failed...")
	}
}

func TestInterfaces(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	if interfaces, err := z.Interfaces(); err != nil {
		t.Errorf(err.Error())
	} else {
		for i, v := range interfaces {
			t.Logf("%d) %s\n", i, v)
		}
	}
}

func TestMemoryUsage(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	memUsed, _ := z.MemoryUsage()
	swapUsed, _ := z.SwapUsage()
	memLimit, _ := z.MemoryLimit()
	swapLimit, _ := z.SwapLimit()

	t.Logf("Mem usage: %0.0f\n", memUsed)
	t.Logf("Mem usage: %s\n", memUsed)
	t.Logf("Swap usage: %0.0f\n", swapUsed)
	t.Logf("Swap usage: %s\n", swapUsed)
	t.Logf("Mem limit: %0.0f\n", memLimit)
	t.Logf("Mem limit: %s\n", memLimit)
	t.Logf("Swap limit: %0.0f\n", swapLimit)
	t.Logf("Swap limit: %s\n", swapLimit)

}

/*
func TestReboot(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	t.Logf("Rebooting the container...\n")
	z.Reboot()

	if z.Running() {
		t.Errorf("Rebooting the container failed...")
	}
}
*/

func TestConcurrentShutdown(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			z := lxc.NewContainer(strconv.Itoa(i))
			defer lxc.PutContainer(z)

			if err := z.Shutdown(30); err != nil {
				t.Errorf(err.Error())
			}
			if z.Running() {
				t.Errorf("Shutting down the container failed...")
			}

			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestShutdown(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	if err := z.Shutdown(30); err != nil {
		t.Errorf(err.Error())
	}

	if z.Running() {
		t.Errorf("Shutting down the container failed...")
	}
}

func TestStop(t *testing.T) {
	z := lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	z.SetDaemonize()
	if err := z.Start(false); err != nil {
		t.Errorf(err.Error())
	}

	if err := z.Stop(); err != nil {
		t.Errorf(err.Error())
	}

	if z.Running() {
		t.Errorf("Stopping the container failed...")
	}
}

func TestDestroy(t *testing.T) {
	z := lxc.NewContainer(CloneOverlayContainerName)
	defer lxc.PutContainer(z)

	if err := z.Destroy(); err != nil {
		t.Errorf(err.Error())
	}

	z = lxc.NewContainer(CloneContainerName)
	defer lxc.PutContainer(z)

	if err := z.Destroy(); err != nil {
		t.Errorf(err.Error())
	}

	z = lxc.NewContainer(ContainerName)
	defer lxc.PutContainer(z)

	if err := z.Destroy(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestConcurrentDestroy(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			z := lxc.NewContainer(strconv.Itoa(i))
			defer lxc.PutContainer(z)

			// sleep for a while to simulate some dummy work
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(250)))

			if err := z.Destroy(); err != nil {
				t.Errorf(err.Error())
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
