package types

import (
	"fmt"
	"strings"
	"sync"
	"text/tabwriter"
	"time"
)

type LDT struct {
	Name    string
	Vendor  string
	Version string
	Os      string
	Arch    string
	Url     string
	Hash    []byte
}

func NewLDT(name, vendor, version, os, arch, url string) *LDT {
	return &LDT{
		Name:    name,
		Vendor:  vendor,
		Version: version,
		Os:      os,
		Arch:    arch,
		Url:     url,
	}
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
	Pid              int
	Ldt              string
	Name             string
	Port             int
	Started          string
	Pairable         bool
	DeviceMacAddress string
}

func NewProcess(pid int, ldt string, name string, port int, deviceAddress string) *Process {
	var pairable bool = true
	if deviceAddress != "" {
		pairable = false
	}

	return &Process{
		Pid:              pid,
		Ldt:              ldt,
		Name:             name,
		Port:             port,
		Started:          time.Now().Format("2006-1-2 15:4:5"),
		Pairable:         pairable,
		DeviceMacAddress: deviceAddress,
	}
}

func (l *LDT) String() string {
	return fmt.Sprintf("%s \t %s \t %s \t %s \t %s \t %s \t %x", l.Name, l.Vendor, l.Version, l.Os, l.Arch, l.Url, l.Hash)
}

func (ll *LDTList) String() string {
	var result strings.Builder
	writer := tabwriter.NewWriter(&result, 0, 0, 3, ' ', 0)
	fmt.Fprintln(writer, "\tUser\tLDT\tVersion\tOS\tArch\tHash")
	for i, ldt := range ll.LDTs {
		fmt.Fprintf(writer, "%d\t%s\t%s\t%s\t%s\t%s\t%x\n", i, ldt.Vendor, ldt.Name, ldt.Version, ldt.Os, ldt.Arch, ldt.Hash[:6])
	}
	writer.Flush()
	return result.String()
}

func (p *Process) LdtType() string {
	return p.Ldt[strings.LastIndex(p.Ldt, "/")+1 : strings.LastIndex(p.Ldt, ":")]
}
