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
	"os"
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
	ContainerName             = "lorem"
	SnapshotName              = "snap0"
	ContainerRestoreName      = "ipsum"
	ContainerCloneName        = "consectetur"
	ContainerCloneOverlayName = "adipiscing"
	Distro                    = "ubuntu"
	Release                   = "saucy"
	Arch                      = "amd64"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func unprivileged() bool {
	if os.Geteuid() != 0 {
		return true
	}
	return false
}

func TestVersion(t *testing.T) {
	t.Logf("LXC version: %s", lxc.Version())
}

func TestDefaultConfigPath(t *testing.T) {
	if lxc.DefaultConfigPath() == "" {
		t.Errorf("DefaultConfigPath failed...")
	}
}

func TestSetConfigPath(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	currentPath := c.ConfigPath()
	if err := c.SetConfigPath("/tmp"); err != nil {
		t.Errorf(err.Error())
	}
	newPath := c.ConfigPath()

	if currentPath == newPath {
		t.Errorf("SetConfigPath failed...")
	}
}

func TestGetContainer(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	lxc.GetContainer(c)
	lxc.PutContainer(c)
}

func TestConcurrentDefined_Negative(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i <= 100; i++ {
		wg.Add(1)
		go func() {
			c, err := lxc.NewContainer(strconv.Itoa(rand.Intn(10)))
			if err != nil {
				t.Errorf(err.Error())
			}
			defer lxc.PutContainer(c)

			// sleep for a while to simulate some dummy work
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(250)))

			if c.Defined() {
				t.Errorf("Defined_Negative failed...")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestDefined_Negative(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if c.Defined() {
		t.Errorf("Defined_Negative failed...")
	}
}

func TestExecute(t *testing.T) {
	if unprivileged() {
		t.Skip("skipping test in unprivileged mode.")
	}

	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if _, err := c.Execute("/bin/true"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestSetVerbosity(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	c.SetVerbosity(lxc.Quiet)
}

func TestCreate(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if unprivileged() {
		if err := c.CreateAsUser(Distro, Release, Arch); err != nil {
			t.Errorf("ERROR: %s\n", err.Error())
		}
	} else {
		if err := c.Create(ContainerType); err != nil {
			t.Errorf(err.Error())
		}
	}
}

func TestClone(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if err := c.Clone(ContainerCloneName); err != nil {
		t.Errorf(err.Error())
	}

	if !unprivileged() {
		if err := c.CloneUsing(ContainerCloneOverlayName, lxc.Overlayfs, lxc.CloneSnapshot|lxc.CloneKeepName|lxc.CloneKeepMACAddr); err != nil {
			t.Errorf(err.Error())
		}
	}
}

func TestCreateSnapshot(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if _, err := c.CreateSnapshot(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestRestoreSnapshot(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	snapshot := lxc.Snapshot{Name: SnapshotName}
	if err := c.RestoreSnapshot(snapshot, ContainerRestoreName); err != nil {
		t.Errorf(err.Error())
	}
}

func TestConcurrentCreate(t *testing.T) {
	if unprivileged() {
		t.Skip("skipping test in unprivileged mode.")
	}

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			c, err := lxc.NewContainer(strconv.Itoa(i))
			if err != nil {
				t.Errorf(err.Error())
			}
			defer lxc.PutContainer(c)

			// sleep for a while to simulate some dummy work
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(250)))

			if err := c.Create(ContainerType); err != nil {
				t.Errorf(err.Error())
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestSnapshots(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if _, err := c.Snapshots(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestConcurrentStart(t *testing.T) {
	if unprivileged() {
		t.Skip("skipping test in unprivileged mode.")
	}

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			c, err := lxc.NewContainer(strconv.Itoa(i))
			if err != nil {
				t.Errorf(err.Error())
			}
			defer lxc.PutContainer(c)

			if err := c.Start(); err != nil {
				t.Errorf(err.Error())
			}

			c.Wait(lxc.RUNNING, 30)
			if !c.Running() {
				t.Errorf("Starting the container failed...")
			}

			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestConfigFileName(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if c.ConfigFileName() == "" {
		t.Errorf("ConfigFileName failed...")
	}
}

func TestDefined_Positive(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if !c.Defined() {
		t.Errorf("Defined_Positive failed...")
	}
}

func TestConcurrentDefined_Positive(t *testing.T) {
	if unprivileged() {
		t.Skip("skipping test in unprivileged mode.")
	}

	var wg sync.WaitGroup

	for i := 0; i <= 100; i++ {
		wg.Add(1)
		go func() {
			c, err := lxc.NewContainer(strconv.Itoa(rand.Intn(10)))
			if err != nil {
				t.Errorf(err.Error())
			}
			defer lxc.PutContainer(c)

			// sleep for a while to simulate some dummy work
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(250)))

			if !c.Defined() {
				t.Errorf("Defined_Positive failed...")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestInitPID_Negative(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if c.InitPID() != -1 {
		t.Errorf("InitPID failed...")
	}
}

func TestStart(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if err := c.Start(); err != nil {
		t.Errorf(err.Error())
	}

	c.Wait(lxc.RUNNING, 30)
	if c.State() != lxc.RUNNING {
		t.Errorf("Starting the container failed...")
	}
}

func TestControllable(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if !c.Controllable() {
		t.Errorf("Controling the container failed...")
	}
}

func TestContainerNames(t *testing.T) {
	if lxc.ContainerNames() == nil {
		t.Errorf("ContainerNames failed...")
	}
}

func TestDefinedContainerNames(t *testing.T) {
	if lxc.DefinedContainerNames() == nil {
		t.Errorf("DefinedContainerNames failed...")
	}
}

func TestActiveContainerNames(t *testing.T) {
	if lxc.ActiveContainerNames() == nil {
		t.Errorf("ActiveContainerNames failed...")
	}
}

func TestContainers(t *testing.T) {
	if lxc.Containers() == nil {
		t.Errorf("Containers failed...")
	}
}

func TestDefinedContainers(t *testing.T) {
	if lxc.DefinedContainers() == nil {
		t.Errorf("DefinedContainers failed...")
	}
}

func TestActiveContainers(t *testing.T) {
	if lxc.ActiveContainers() == nil {
		t.Errorf("ActiveContainers failed...")
	}
}

func TestRunning(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if !c.Running() {
		t.Errorf("Checking the container failed...")
	}
}

func TestWantDaemonize(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if err := c.WantDaemonize(false); err != nil || c.Daemonize() {
		t.Errorf("WantDaemonize failed...")
	}
}

func TestWantCloseAllFds(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if err := c.WantCloseAllFds(true); err != nil {
		t.Errorf("WantCloseAllFds failed...")
	}
}

func TestSetLogLevel(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if err := c.SetLogLevel(lxc.WARN); err != nil || c.LogLevel() != lxc.WARN {
		t.Errorf("SetLogLevel( failed...")
	}
}

func TestSetLogFile(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if err := c.SetLogFile("/tmp/" + ContainerName); err != nil || c.LogFile() != "/tmp/"+ContainerName {
		t.Errorf("SetLogFile failed...")
	}
}

func TestInitPID_Positive(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if c.InitPID() == -1 {
		t.Errorf("InitPID failed...")
	}
}

func TestName(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if c.Name() != ContainerName {
		t.Errorf("Name failed...")
	}
}

func TestFreeze(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if err := c.Freeze(); err != nil {
		t.Errorf(err.Error())
	}

	c.Wait(lxc.FROZEN, 30)
	if c.State() != lxc.FROZEN {
		t.Errorf("Freezing the container failed...")
	}
}

func TestUnfreeze(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if err := c.Unfreeze(); err != nil {
		t.Errorf(err.Error())
	}

	c.Wait(lxc.RUNNING, 30)
	if c.State() != lxc.RUNNING {
		t.Errorf("Unfreezing the container failed...")
	}
}

func TestLoadConfigFile(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if err := c.LoadConfigFile(c.ConfigFileName()); err != nil {
		t.Errorf(err.Error())
	}
}

func TestSaveConfigFile(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if err := c.SaveConfigFile(c.ConfigFileName()); err != nil {
		t.Errorf(err.Error())
	}
}

func TestConfigItem(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if c.ConfigItem("lxc.utsname")[0] != ContainerName {
		t.Errorf("ConfigItem failed...")
	}
}

func TestSetConfigItem(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if err := c.SetConfigItem("lxc.utsname", ContainerName); err != nil {
		t.Errorf(err.Error())
	}

	if c.ConfigItem("lxc.utsname")[0] != ContainerName {
		t.Errorf("ConfigItem failed...")
	}
}

func TestSetCgroupItem(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	maxMem := c.CgroupItem("memory.max_usage_in_bytes")[0]
	currentMem := c.CgroupItem("memory.limit_in_bytes")[0]
	if err := c.SetCgroupItem("memory.limit_in_bytes", maxMem); err != nil {
		t.Errorf(err.Error())
	}
	newMem := c.CgroupItem("memory.limit_in_bytes")[0]

	if newMem == currentMem {
		t.Errorf("SetCgroupItem failed...")
	}
}

func TestClearConfigItem(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if err := c.ClearConfigItem("lxc.cap.drop"); err != nil {
		t.Errorf(err.Error())
	}
	if c.ConfigItem("lxc.cap.drop")[0] != "" {
		t.Errorf("ClearConfigItem failed...")
	}
}

func TestConfigKeys(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	keys := strings.Join(c.ConfigKeys("lxc.network.0"), " ")
	if !strings.Contains(keys, "mtu") {
		t.Errorf("Keys failed...")
	}
}

func TestInterfaces(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if _, err := c.Interfaces(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestMemoryUsage(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if _, err := c.MemoryUsage(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestKernelMemoryUsage(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if _, err := c.KernelMemoryUsage(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestMemorySwapUsage(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if _, err := c.MemorySwapUsage(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestBlkioUsage(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if _, err := c.BlkioUsage(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestMemoryLimit(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if _, err := c.MemoryLimit(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestKernelMemoryLimit(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if _, err := c.KernelMemoryLimit(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestMemorySwapLimit(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if _, err := c.MemorySwapLimit(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestSetMemoryLimit(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	oldMemLimit, _ := c.MemoryLimit()

	if err := c.SetMemoryLimit(oldMemLimit * 4); err != nil {
		t.Errorf(err.Error())
	}

	newMemLimit, _ := c.MemoryLimit()
	if newMemLimit != oldMemLimit*4 {
		t.Errorf("SetMemoryLimit failed")
	}
}

func TestSetKernelMemoryLimit(t *testing.T) {
	t.Skip("skipping test")

	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	oldMemLimit, _ := c.KernelMemoryLimit()
	if err := c.SetKernelMemoryLimit(oldMemLimit * 4); err != nil {
		t.Errorf(err.Error())
	}

	newMemLimit, _ := c.KernelMemoryLimit()
	if newMemLimit != oldMemLimit*4 {
		t.Errorf("SetKernelMemoryLimit failed")
	}
}

func TestSetMemorySwapLimit(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	oldMemorySwapLimit, _ := c.MemorySwapLimit()

	if err := c.SetMemorySwapLimit(oldMemorySwapLimit / 4); err != nil {
		t.Errorf(err.Error())
	}

	newMemorySwapLimit, _ := c.MemorySwapLimit()
	if newMemorySwapLimit != oldMemorySwapLimit/4 {
		t.Errorf("SetSwapLimit failed")
	}
}

func TestCPUTime(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if _, err := c.CPUTime(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestCPUTimePerCPU(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if _, err := c.CPUTimePerCPU(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestCPUStats(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if _, err := c.CPUStats(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestRunCommand(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	argsThree := []string{"/bin/sh", "-c", "/bin/ls -al > /dev/null"}
	if err := c.RunCommand(os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(), argsThree...); err != nil {
		t.Errorf(err.Error())
	}
}

func TestRunCommandWithClearEnvironment(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	argsThree := []string{"/bin/sh", "-c", "/bin/ls -al > /dev/null"}
	if err := c.RunCommandWithClearEnvironment(os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(), argsThree...); err != nil {
		t.Errorf(err.Error())
	}
}

func TestConsoleGetFD(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if _, err := c.ConsoleGetFD(0); err != nil {
		t.Errorf(err.Error())
	}
}

func TestIPAddress(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if _, err := c.IPAddress("lo"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAddDeviceNode(t *testing.T) {
	if unprivileged() {
		t.Skip("skipping test in unprivileged mode.")
	}

	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if err := c.AddDeviceNode("/dev/network_latency"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestRemoveDeviceNode(t *testing.T) {
	if unprivileged() {
		t.Skip("skipping test in unprivileged mode.")
	}

	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if err := c.RemoveDeviceNode("/dev/network_latency"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestIPv4Addresses(t *testing.T) {
	t.Skip("skipping test")

	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if _, err := c.IPv4Addresses(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestIPv6Addresses(t *testing.T) {
	t.Skip("skipping test")

	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if _, err := c.IPv6Addresses(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestReboot(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if err := c.Reboot(); err != nil {
		t.Errorf("Rebooting the container failed...")
	}
	c.Wait(lxc.RUNNING, 30)
}

func TestConcurrentShutdown(t *testing.T) {
	if unprivileged() {
		t.Skip("skipping test in unprivileged mode.")
	}

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			c, err := lxc.NewContainer(strconv.Itoa(i))
			if err != nil {
				t.Errorf(err.Error())
			}
			defer lxc.PutContainer(c)

			if err := c.Shutdown(3); err != nil {
				t.Errorf(err.Error())
			}

			c.Wait(lxc.STOPPED, 30)
			if c.Running() {
				t.Errorf("Shutting down the container failed...")
			}

			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestShutdown(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if err := c.Shutdown(3); err != nil {
		t.Errorf(err.Error())
	}

	c.Wait(lxc.STOPPED, 30)
	if c.Running() {
		t.Errorf("Shutting down the container failed...")
	}
}

func TestStop(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if err := c.Start(); err != nil {
		t.Errorf(err.Error())
	}

	if err := c.Stop(); err != nil {
		t.Errorf(err.Error())
	}

	c.Wait(lxc.STOPPED, 30)
	if c.Running() {
		t.Errorf("Stopping the container failed...")
	}
}

func TestDestroySnapshot(t *testing.T) {
	c, err := lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	snapshot := lxc.Snapshot{Name: SnapshotName}
	if err := c.DestroySnapshot(snapshot); err != nil {
		t.Errorf(err.Error())
	}
}

func TestDestroy(t *testing.T) {
	if !unprivileged() {
		c, err := lxc.NewContainer(ContainerCloneOverlayName)
		if err != nil {
			t.Errorf(err.Error())
		}
		defer lxc.PutContainer(c)

		if err := c.Destroy(); err != nil {
			t.Errorf(err.Error())
		}
	}

	c, err := lxc.NewContainer(ContainerCloneName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if err := c.Destroy(); err != nil {
		t.Errorf(err.Error())
	}

	c, err = lxc.NewContainer(ContainerRestoreName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if err := c.Destroy(); err != nil {
		t.Errorf(err.Error())
	}

	c, err = lxc.NewContainer(ContainerName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lxc.PutContainer(c)

	if err := c.Destroy(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestConcurrentDestroy(t *testing.T) {
	if unprivileged() {
		t.Skip("skipping test in unprivileged mode.")
	}

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			c, err := lxc.NewContainer(strconv.Itoa(i))
			if err != nil {
				t.Errorf(err.Error())
			}
			defer lxc.PutContainer(c)

			// sleep for a while to simulate some dummy work
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(250)))

			if err := c.Destroy(); err != nil {
				t.Errorf(err.Error())
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
