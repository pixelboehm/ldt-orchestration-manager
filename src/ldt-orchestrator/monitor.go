package ldtorchestrator

import (
	"log"
)

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
			m.processes = append(m.processes, ldt)
			log.Printf("New LDT %s with PID %d started at %s\n", ldt.Name, ldt.Pid, ldt.started.Format("02-01-2006 15:04:05"))
		default:
		}
	}
}
