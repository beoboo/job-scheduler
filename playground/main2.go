package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

const ROOT_FILESYSTEM = "rootfs"

func main() {
	switch os.Args[1] {
	case "run":
		parent()
	case "child":
		child()
	default:
		panic("wat should I do")
	}
}

func parent() {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	if err := cmd.Run(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
}

func child() {
	// From the pivot_root man:
	// The following restrictions apply to new_root and put_old:
	//
	// - They must be directories.
	//
	// - new_root and put_old must not be on the same file system as the current
	//   root.
	//
	// - put_old must be underneath new_root, that is, adding a nonzero number of
	//   .. to the string pointed to by put_old must yield the same directory as
	//   new_root.
	//
	//  - No other file system may be mounted on put_old.
	newRoot := "/new"
	pivotRoot := ".pivot_root"
	pivotDir := filepath.Join(newRoot, pivotRoot)

	fmt.Printf("Mounting %s\n", newRoot)
	must(syscall.Mount(newRoot, newRoot, "bind", syscall.MS_BIND|syscall.MS_REC, ""))
	fmt.Printf("Creating pivot root on %s\n", pivotDir)
	must(os.MkdirAll(pivotDir, 0777))
	fmt.Printf("Pivoting root from %s to %s\n", newRoot, pivotDir)
	must(syscall.PivotRoot(newRoot, pivotDir))
	fmt.Printf("Changing dir to /\n")
	must(os.Chdir("/"))
	fmt.Printf("5")

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	err := syscall.Sethostname([]byte("child"))
	if err != nil {
		return
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
}

func must(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
