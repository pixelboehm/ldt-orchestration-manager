package ldtorchestrator

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type Process struct {
	Pid     int
	Name    string
	started time.Time
}

type ProcessList []Process

var runningProcesses ProcessList

func DownloadLDT(url string) (string, error) {
	file, err := os.Create("./resources/child_webserver")
	if err != nil {
		return "", err
	}
	defer file.Close()

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	if err := os.Chmod(file.Name(), 0755); err != nil {
		log.Fatalf("Could not set executable flag: %v", err)
	}

	log.Printf("Downloaded LDT: %s\n", file.Name())
	return file.Name(), nil
}

func StartLDT(name string) ProcessList {
	cmd := exec.Command("./" + name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		log.Fatal("Could not start LDT\n")
	}

	fmt.Printf("Started LDT with PID %d\n", cmd.Process.Pid)
	runningProcesses = append(runningProcesses, Process{
		Pid:     cmd.Process.Pid,
		Name:    name,
		started: time.Now(),
	})
	return runningProcesses
}

func StopLDT(pid int, graceful bool) {
	proc, err := os.FindProcess(pid)
	if err != nil {
		log.Fatal(err)
		return
	}
	if graceful == true {
		err = proc.Signal(os.Interrupt)
	} else {
		err = proc.Kill()
	}

	if err != nil {
		log.Fatalf("Unable to stop LDT \t graceful? %t", graceful)
	}
}
