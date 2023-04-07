package types

import "sync"

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
