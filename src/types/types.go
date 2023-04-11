package types

import (
	"fmt"
	"strings"
	"sync"
	"text/tabwriter"
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
	var result strings.Builder
	writer := tabwriter.NewWriter(&result, 0, 0, 3, ' ', 0)
	fmt.Fprintln(writer, "Name\tVersion\tOS\tArch")
	for _, ldt := range ll.LDTs {
		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", ldt.Name, ldt.Version, ldt.Os, ldt.Arch)
	}
	writer.Flush()

	return result.String()
}
