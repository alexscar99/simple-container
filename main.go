package main

// Taken from talk given by @lizrice: https://gist.github.com/lizrice/a5ef4d175fd0cd3491c7e8d716826d27

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	default:
		panic("what???")
	}
}

// Sets up our command and then we are calling Run to actually run it
func run() {
	fmt.Printf("running %v\n", os.Args[2:])

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	/*
		SysProcAttr holds optional, operating system-specific attributes.
		Run passes it to os.StartProcess as the os.ProcAttr's Sys field.

		SysProcAttr is a struct that has a field of 'Cloneflags' that is
		of type 'uintptr' and it flags for clone calls (only on Linux).

		We're only going to pass in one flag and it is 'syscall.CLONE_NEWUTS'.
		This will clone a new UTS (Unix Time Sharing), which is the
		namespace that specifically isolates the hostname.
	*/

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}

	must(cmd.Run())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
