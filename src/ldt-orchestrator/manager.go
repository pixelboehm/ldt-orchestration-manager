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

type Manager struct {
	monitor   *Monitor
	discovery *DiscoveryConfig
}

func (manager *Manager) Setup(config string) {
	if err := manager.monitor.DeserializeLDTs(); err != nil {
		log.Fatal(err)
	}
	manager.monitor = NewMonitor()
	manager.discovery = NewConfig(config)
}

func (manager *Manager) Run() {
	go manager.monitor.RefreshLDTs()
	ticker := time.NewTicker(10 * time.Second)

	var name, pkg_type, dist string
	manager.discovery.GetLDTs(name, pkg_type, dist)
	ldt, err := manager.DownloadLDT("http://localhost:8081/getPackage")
	if err != nil {
		log.Fatal(err)
	}
	process, err := manager.StartLDT(ldt)
	if err != nil {
		log.Fatal(err)
	}
	manager.monitor.started <- *process
	<-ticker.C

	if err := manager.StopLDT(process.Pid, true); err != nil {
		log.Fatal(err)
	}
	manager.monitor.stopped <- process.Pid
	log.Printf("Stopped LDT %s with PID %d\n", process.Name, process.Pid)
	manager.shutdown()
}

func (manager *Manager) DownloadLDT(url string) (string, error) {
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

func (manager *Manager) StartLDT(name string) (*Process, error) {
	cmd := exec.Command("./" + name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		log.Fatal("Could not start LDT\n")
		return nil, err
	}

	fmt.Printf("Started LDT with PID %d\n", cmd.Process.Pid)
	return &Process{
		Pid:     cmd.Process.Pid,
		Name:    name,
		started: time.Now(),
	}, nil

}

func (manager *Manager) StopLDT(pid int, graceful bool) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if graceful == true {
		err = proc.Signal(os.Interrupt)
	} else {
		err = proc.Kill()
	}

	if err != nil {
		log.Fatalf("Unable to stop LDT \t graceful? %t", graceful)
		return err
	}
	return nil
}

func (manager *Manager) shutdown() {
	if err := manager.monitor.SerializeLDTs(); err != nil {
		log.Fatal(err)
	}
}
