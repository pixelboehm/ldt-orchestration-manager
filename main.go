package main

import (
	"fmt"
	lo "longevity/src/ldt-orchestrator"
	"os"
	"runtime"
	"time"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	default:
		panic("Don't know what to do")
	}
}

func run() {
	manager := &lo.Manager{RunningProcesses: make([]lo.Process, 0)}
	manager.Run()
}

func timer() func() {
	name := callerName(1)
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func callerName(skip int) string {
	const unknown = "unknown"
	pcs := make([]uintptr, 1)
	n := runtime.Callers(skip+2, pcs)
	if n < 1 {
		return unknown
	}
	frame, _ := runtime.CallersFrames(pcs).Next()
	if frame.Function == "" {
		return unknown
	}
	return frame.Function
}
