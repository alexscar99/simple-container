package main

// Taken from talk given by @lizrice: https://gist.github.com/lizrice/a5ef4d175fd0cd3491c7e8d716826d27

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// docker run <container> <cmd> <args>
// go run main.go run <cmd> <args>
func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("what???")
	}
}

func run() {
	/*

		Running "/proc/self/exe" is basically running a fork and exec

		fork and exec --> 'fork' starts a new process which is a copy of the
		one that calls it, while 'exec' replaces the current process image
		with another (different) one

	*/

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	/*

		SysProcAttr holds optional, operating system-specific attributes.
		Run passes it to os.StartProcess as the os.ProcAttr's Sys field.

		SysProcAttr is a struct that has a field of 'Cloneflags' that is
		of type 'uintptr' and it flags for clone calls (only on Linux).

		What we are trying to clone here are different namespaces. The first
		is UTS (UNIX Time Sharing) that deals with changing the hostname.

		The second is PID (process ID) which provides processes with an independent
		set of process IDs (PIDs) from other namespaces.

	*/

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}

	must(cmd.Run())
}

func child() {
	fmt.Printf("running %v as PID %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Use the separate root file system
	must(syscall.Chroot("/home/rootfs"))
	must(os.Chdir("/"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))
	must(cmd.Run())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
