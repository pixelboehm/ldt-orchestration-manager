package main

import (
	"fmt"
	"log"
	lo "longevity/src/ldt-orchestrator"
	"runtime"
	"time"
)

func main() {
	defer timer()()
	var name, pkg_type, dist string
	lo.GetPackages(name, pkg_type, dist)
	pkg, err := lo.DownloadPackage("http://localhost:8081/getPackage")
	if err != nil {
		log.Fatal(err)
	}
	lo.StartPackageDetached(pkg)
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
