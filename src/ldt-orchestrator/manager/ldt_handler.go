package manager

import (
	"io"
	"log"
	. "longevity/src/types"
	"net"
	"os"
	"os/exec"
	"syscall"
)

func run(ldt string) (*Process, error) {
	makeExecutable(ldt)

	cmd := exec.Command("./" + ldt)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	process := NewProcess(cmd.Process.Pid, ldt)
	return process, nil
}

func start(ldt string, in *net.Conn) (*Process, *io.Writer, error) {
	makeExecutable(ldt)

	cmd := exec.Command("./" + ldt)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	cmd.Stdout = *in
	cmd.Stderr = *in
	cmd.Stdin = *in

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}
	process := NewProcess(cmd.Process.Pid, ldt)

	return process, &cmd.Stdout, nil
}

func stop(pid int, graceful bool) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		log.Printf("Failed to find process with PID %d\n", pid)
	}
	if graceful {
		if err = proc.Signal(os.Interrupt); err != nil {
			log.Printf("Failed to stop LDT %d gracefully\n", pid)
			return false
		}
	} else {
		if err = proc.Kill(); err != nil {
			log.Printf("Failed to kill LDT %d\n", pid)
			return false
		}
	}
	return true

}

func makeExecutable(ldt string) {
	if _, err := os.Stat(ldt); os.IsNotExist(err) {
		return
	}
	file, err := os.Open(ldt)
	if err != nil {
		log.Fatal("Failed to open LDT: ", err)
	}

	if err := os.Chmod(file.Name(), 0755); err != nil {
		log.Fatal("Failed to set executable Flag: ", err)
	}
}
