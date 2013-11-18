// Copyright © 2013, S.Çağlar Onur
// Use of this source code is governed by a LGPLv2.1
// license that can be found in the LICENSE file.
//
// Authors:
// S.Çağlar Onur <caglar@10ur.org>

// +build linux

package lxc_test

import (
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/caglar10ur/lxc"
)

const (
	ContainerType             = "busybox"
	ContainerName             = "rubik"
	SnapshotName              = "snap0"
	ContainerRestoreName      = "rubik-restore"
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
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
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
			z, err := lxc.NewContainer(strconv.Itoa(rand.Intn(10)))
			if err != nil {
				t.Errorf(err.Error())
			}
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
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if z.Defined() {
		t.Errorf("Defined_Negative failed...")
	}
}

func TestCreate(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if err := z.Create(ContainerType); err != nil {
		t.Errorf(err.Error())
	}
}

func TestClone(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if err := z.CloneToDirectory(CloneContainerName); err != nil {
		t.Errorf(err.Error())
	}

	if err := z.CloneToOverlayFS(CloneOverlayContainerName); err != nil {
		t.Errorf(err.Error())
	}
}

func TestCreateSnapshot(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if err := z.CreateSnapshot(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestRestoreSnapshot(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	snapshot := lxc.Snapshot{Name: SnapshotName}
	if err := z.RestoreSnapshot(snapshot, ContainerRestoreName); err != nil {
		t.Errorf(err.Error())
	}
}

func TestConcurrentCreate(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			z, err := lxc.NewContainer(strconv.Itoa(i))
			if err != nil {
				t.Errorf(err.Error())
			}
			defer lxc.PutContainer(z)

			// sleep for a while to simulate some dummy work
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(250)))

			if err := z.Create(ContainerType); err != nil {
				t.Errorf(err.Error())
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestSnapshots(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if _, err := z.Snapshots(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestConcurrentStart(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			z, err := lxc.NewContainer(strconv.Itoa(i))
			if err != nil {
				t.Errorf(err.Error())
			}
			defer lxc.PutContainer(z)

			z.SetDaemonize()
			if err := z.Start(); err != nil {
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

func TestContainerNames(t *testing.T) {
	if lxc.ContainerNames() == nil {
		t.Errorf("ContainerNames failed...")
	}
}

func TestActiveContainerNames(t *testing.T) {
	if lxc.ActiveContainerNames() == nil {
		t.Errorf("ContainerNames failed...")
	}
}

func TestContainers(t *testing.T) {
	if lxc.Containers() == nil {
		t.Errorf("Containers failed...")
	}
}

func TestActiveContainers(t *testing.T) {
	if lxc.ActiveContainers() == nil {
		t.Errorf("ActiveContainers failed...")
	}
}

func TestConfigFileName(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if z.ConfigFileName() != ConfigFileName {
		t.Errorf("ConfigFileName failed...")
	}
}

func TestDefined_Positive(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
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
			z, err := lxc.NewContainer(strconv.Itoa(rand.Intn(10)))
			if err != nil {
				t.Errorf(err.Error())
			}
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
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if z.InitPID() != -1 {
		t.Errorf("InitPID failed...")
	}
}

func TestStart(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	z.SetDaemonize()
	if err := z.Start(); err != nil {
		t.Errorf(err.Error())
	}

	z.Wait(lxc.RUNNING, 30)
	if z.State() != lxc.RUNNING {
		t.Errorf("Starting the container failed...")
	}
}

func TestMayControl(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if !z.MayControl() {
		t.Errorf("Controling the container failed...")
	}
}

func TestRunning(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if !z.Running() {
		t.Errorf("Checking the container failed...")
	}
}
func TestSetDaemonize(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	z.SetDaemonize()
	if !z.Daemonize() {
		t.Errorf("Daemonize failed...")
	}
}

func TestInitPID_Positive(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if z.InitPID() == -1 {
		t.Errorf("InitPID failed...")
	}
}

func TestName(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if z.Name() != ContainerName {
		t.Errorf("Name failed...")
	}
}

func TestFreeze(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
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
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
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
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if err := z.LoadConfigFile(ConfigFileName); err != nil {
		t.Errorf(err.Error())
	}
}

func TestSaveConfigFile(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if err := z.SaveConfigFile(ConfigFileName); err != nil {
		t.Errorf(err.Error())
	}
}

func TestConfigItem(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if z.ConfigItem("lxc.utsname")[0] != ContainerName {
		t.Errorf("ConfigItem failed...")
	}
}

func TestSetConfigItem(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if err := z.SetConfigItem("lxc.utsname", ContainerName); err != nil {
		t.Errorf(err.Error())
	}

	if z.ConfigItem("lxc.utsname")[0] != ContainerName {
		t.Errorf("ConfigItem failed...")
	}
}

func TestSetCgroupItem(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
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
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if err := z.ClearConfigItem("lxc.cap.drop"); err != nil {
		t.Errorf(err.Error())
	}
	if z.ConfigItem("lxc.cap.drop")[0] != "" {
		t.Errorf("ClearConfigItem failed...")
	}
}

func TestConfigKeys(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	keys := strings.Join(z.ConfigKeys("lxc.network.0"), " ")
	if !strings.Contains(keys, "mtu") {
		t.Errorf("Keys failed...")
	}
}

func TestInterfaces(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if _, err := z.Interfaces(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestMemoryUsage(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if _, err := z.MemoryUsage(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestSwapUsage(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if _, err := z.SwapUsage(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestMemoryLimit(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if _, err := z.MemoryLimit(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestSwapLimit(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if _, err := z.SwapLimit(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestSetMemoryLimit(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	oldMemLimit, _ := z.MemoryLimit()

	if err := z.SetMemoryLimit(oldMemLimit * 4); err != nil {
		t.Errorf(err.Error())
	}

	newMemLimit, _ := z.MemoryLimit()
	if newMemLimit != 4*oldMemLimit {
		t.Errorf("SetMemoryLimit failed")
	}
}

func TestSetSwapLimit(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	oldSwapLimit, _ := z.SwapLimit()

	if err := z.SetSwapLimit(oldSwapLimit / 4); err != nil {
		t.Errorf(err.Error())
	}

	newSwapLimit, _ := z.SwapLimit()
	if newSwapLimit != oldSwapLimit/4 {
		t.Errorf("SetSwapLimit failed")
	}
}

func TestAttachRunCommand(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	argsThree := []string{"/bin/sh", "-c", "/bin/ls -al > /dev/null"}
	if err := z.AttachRunCommand(argsThree...); err != nil {
		t.Errorf(err.Error())
	}
}

func TestConsoleGetFD(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if _, err := z.ConsoleGetFD(0); err != nil {
		t.Errorf(err.Error())
	}
}

func TestIPAddress(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if _, err := z.IPAddress("lo"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAddDeviceNode(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if err := z.AddDeviceNode("/dev/network_latency"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestRemoveDeviceNode(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if err := z.RemoveDeviceNode("/dev/network_latency"); err != nil {
		t.Errorf(err.Error())
	}
}

/*
func TestIPv4Addresses(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
    defer lxc.PutContainer(z)

    if _, err := z.IPv4Addresses(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestIPv6Addresses(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
    defer lxc.PutContainer(z)

    if _, err := z.IPv6Addresses(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestReboot(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
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
			z, err := lxc.NewContainer(strconv.Itoa(i))
			if err != nil {
				t.Errorf(err.Error())
			}
			defer lxc.PutContainer(z)

			if err := z.Shutdown(3); err != nil {
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
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if err := z.Shutdown(3); err != nil {
		t.Errorf(err.Error())
	}

	if z.Running() {
		t.Errorf("Shutting down the container failed...")
	}
}

func TestStop(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	z.SetDaemonize()
	if err := z.Start(); err != nil {
		t.Errorf(err.Error())
	}

	if err := z.Stop(); err != nil {
		t.Errorf(err.Error())
	}

	if z.Running() {
		t.Errorf("Stopping the container failed...")
	}
}

func TestDestroySnapshot(t *testing.T) {
	z, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	snapshot := lxc.Snapshot{Name: SnapshotName}
	if err := z.DestroySnapshot(snapshot); err != nil {
		t.Errorf(err.Error())
	}
}

func TestDestroy(t *testing.T) {
	z, err := lxc.NewContainer(CloneOverlayContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if err := z.Destroy(); err != nil {
		t.Errorf(err.Error())
	}

	z, err = lxc.NewContainer(CloneContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if err := z.Destroy(); err != nil {
		t.Errorf(err.Error())
	}

	z, err = lxc.NewContainer(ContainerRestoreName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(z)

	if err := z.Destroy(); err != nil {
		t.Errorf(err.Error())
	}

	z, err = lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
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
			z, err := lxc.NewContainer(strconv.Itoa(i))
			if err != nil {
				t.Errorf(err.Error())
			}
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
