package types

import (
	"fmt"
	"sync"
)

type LDT struct {
	Name    string
	Version string
	Os      string
	Arch    string
	Url     string
}

type LDTList struct {
	LDTs []LDT
	Lock sync.Mutex
}

func NewLDTList() *LDTList {
	return &LDTList{
		LDTs: nil,
		Lock: sync.Mutex{},
	}
}

func (l *LDT) String() string {
	return fmt.Sprintf("%s \t %s \t %s \t %s \t %s", l.Name, l.Version, l.Os, l.Arch, l.Url)
}

func (ll *LDTList) String() string {
	var result string = "Name \t Version \t OS \t Arch \t URL\n"
	for _, ldt := range ll.LDTs {
		result += fmt.Sprintf("%s \t %s \t %s \t %s \t %s\n",
			ldt.Name, ldt.Version, ldt.Os, ldt.Arch, ldt.Url)
	}

	return result
}
