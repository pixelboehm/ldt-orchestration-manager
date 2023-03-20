package ldtorchestrator

import (
	"log"
)

type Monitor struct {
	ldts chan Process
}

func (m *Monitor) RefreshLDTs() {
	for {
		select {
		case ldt := <-m.ldts:
			log.Printf("New LDT %s with PID %d started at %s\n", ldt.Name, ldt.Pid, ldt.started)
		default:
		}
	}
}
