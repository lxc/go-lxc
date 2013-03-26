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
		fmt.Printf("Create: %+v\n", z.Create("ubuntu"))
	} else {
		fmt.Printf("Starting rubik container...\n\n")
		fmt.Printf("Start: %+v\n", z.Start(false))
		fmt.Printf("State: %+v\n", z.State())
		fmt.Printf("Init PID: %+v\n", z.InitPID())
		fmt.Printf("Freeze: %+v\n", z.Freeze())
		fmt.Printf("State: %+v\n", z.State())
		fmt.Printf("Unfreeze: %+v\n", z.Unfreeze())
		fmt.Printf("State: %+v\n", z.State())
	}

	if z.Running() {
		fmt.Printf("Shutdown: %+v\n", z.Shutdown(10))
		fmt.Printf("State: %+v\n", z.State())
		fmt.Printf("Stop: %+v\n", z.Stop())
		fmt.Printf("State: %+v\n", z.State())
	}
}
