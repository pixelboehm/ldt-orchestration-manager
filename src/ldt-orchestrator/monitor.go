package ldtorchestrator

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

// todo: seperate channel for started and stopped LDT
type Monitor struct {
	ldts      chan Process
	processes []Process
}

func NewMonitor() *Monitor {
	return &Monitor{ldts: make(chan Process)}
}

func (m *Monitor) RefreshLDTs() {
	for {
		select {
		case ldt := <-m.ldts:
			m.RegisterLDT(ldt)
		default:
		}
	}
}

func (m *Monitor) RegisterLDT(ldt Process) {
	m.processes = append(m.processes, ldt)
	if err := m.SerializeLDTs(); err != nil {
		log.Fatal(err)
	}
	log.Printf("New LDT %s with PID %d started at %s\n", ldt.Name, ldt.Pid, ldt.started.Format("02-01-2006 15:04:05"))
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
	filename := "src/ldt-orchestrator/ldt.list"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Could not open file: %s", filename)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var name string
		var pid int
		var started string
		fmt.Sscanf(scanner.Text(), "%s\t%d\t%s", &name, &pid, &started)

		time, err := time.Parse("02-01-2006 15:04:05", started)
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
	return nil
}
