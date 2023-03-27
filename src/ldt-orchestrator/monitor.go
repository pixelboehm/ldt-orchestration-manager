package ldtorchestrator

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	ldt_list = "src/ldt-orchestrator/ldt.list"
)

type Monitor struct {
	started   chan Process
	stopped   chan int
	processes []Process
}

func NewMonitor() *Monitor {
	return &Monitor{
		started: make(chan Process),
		stopped: make(chan int),
	}
}

func (m *Monitor) RefreshLDTs() {
	for {
		select {
		case started := <-m.started:
			m.RegisterLDT(started)
		case stopped := <-m.stopped:
			m.RemoveLDT(stopped)
		default:
		}
	}
}

func (m *Monitor) RegisterLDT(ldt Process) {
	m.processes = append(m.processes, ldt)
	log.Printf("New LDT %s with PID %d started at %s\n", ldt.Name, ldt.Pid, ldt.started.Format("02-01-2006 15:04:05"))
}

func (m *Monitor) RemoveLDT(pid int) {
	for i, ldt := range m.processes {
		if ldt.Pid == pid {
			m.processes = append(m.processes[:i], m.processes[i+1:]...)
		}
	}
	log.Printf("Removing LDT with PID %d\n", pid)
}

func (m *Monitor) SerializeLDTs() error {
	filename := "src/ldt-orchestrator/ldt.list"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("Could not create file: %s", filename)
		return err
	}
	defer file.Close()

	template := "%s\t%d\t%s\n"
	writer := bufio.NewWriter(file)
	for _, ldt := range m.processes {
		res := fmt.Sprintf(template, ldt.Name, ldt.Pid, ldt.started.Format("02-01-2006 15:04:05"))
		writer.WriteString(res)
	}

	writer.Flush()
	return nil
}

func (m *Monitor) DeserializeLDTs() error {
	if checkFileExists(ldt_list) {
		file, err := os.Open(ldt_list)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var name string
			var pid int
			var day string
			var hour string
			fmt.Sscanf(scanner.Text(), "%s\t%d\t%s%s", &name, &pid, &day, &hour)

			time, err := time.Parse("02-01-2006 15:04:05", day+" "+hour)
			if err != nil {
				log.Fatal(err)
				return err
			}

			m.processes = append(m.processes, Process{Pid: pid, Name: name, started: time})
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}
