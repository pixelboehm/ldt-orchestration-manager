package types

import (
	"fmt"
	"log"
	wotm "longevity/src/wot-manager"
	"strings"
	"sync"
	"text/tabwriter"
	"time"
)

type LDT struct {
	Name    string
	User    string
	Version string
	Os      string
	Arch    string
	Url     string
	Hash    []byte
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

type Process struct {
	Pid     int
	Ldt     string
	Name    string
	Desc    wotm.WoTDescription
	Started time.Time
}

func NewProcess(pid int, ldt string, name string) *Process {
	wotm, err := wotm.NewWoTmanager(ldt)
	if err != nil {
		log.Fatal(err)
	}
	desc, err := wotm.FetchWoTDescription()
	if err != nil {
		log.Fatal(err)
	}

	return &Process{
		Pid:     pid,
		Ldt:     ldt,
		Name:    name,
		Desc:    desc,
		Started: time.Now(),
	}
}

func (l *LDT) String() string {
	return fmt.Sprintf("%s \t %s \t %s \t %s \t %s \t %s \t %x", l.Name, l.User, l.Version, l.Os, l.Arch, l.Url, l.Hash)
}

func (ll *LDTList) String() string {
	var result strings.Builder
	writer := tabwriter.NewWriter(&result, 0, 0, 3, ' ', 0)
	fmt.Fprintln(writer, "\tUser\tLDT\tVersion\tOS\tArch\tHash")
	for i, ldt := range ll.LDTs {
		fmt.Fprintf(writer, "%d\t%s\t%s\t%s\t%s\t%s\t%x\n", i, ldt.Name, ldt.User, ldt.Version, ldt.Os, ldt.Arch, ldt.Hash[:6])
	}
	writer.Flush()
	return result.String()
}
