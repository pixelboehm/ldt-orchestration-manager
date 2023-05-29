package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	socketPath = "/tmp/orchestration-manager.sock"
)

func main() {
	process := make(chan int, 1)

	connection, err := net.Dial("unix", socketPath)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	var command string = prepareExecutionCommand(os.Args)
	_, err = connection.Write([]byte(command))
	if err != nil {
		log.Fatal("write error:", err)
	}

	if os.Args[1] == "start" {
		go checkForShutdown(connection, process)
		waitForAnswer(connection, process, true)
	} else {
		waitForAnswer(connection, nil, false)
	}
}

func waitForAnswer(connection net.Conn, process chan int, blocking bool) {
	var pid int = 0
	if blocking {
		go checkIfAttachedProcessIsStillRunning(&pid)
	}

	buffer := make([]byte, 4096)
	for {
		n, err := connection.Read(buffer[:])
		if err != nil {
			return
		}
		val := string(buffer[0:n])
		val = strings.Trim(val, "\n")
		if blocking {
			if pid == 0 {
				pid, err = strconv.Atoi(val)
				if err == nil {
					process <- pid
				}
				continue
			}
		}
		fmt.Println(val)
		if !blocking {
			return
		}
	}
}

func checkIfAttachedProcessIsStillRunning(pid *int) {
	for {
		ticker := time.NewTicker(100 * time.Millisecond)
		if *pid == 0 {
			continue
		}
		process, err := os.FindProcess(*pid)
		if err != nil {
			log.Println(err)
			return
		}
		err = process.Signal(syscall.Signal(0))
		if err != nil {
			os.Exit(0)
		}
		<-ticker.C
	}
}

func prepareExecutionCommand(args []string) string {
	var command string
	for _, arg := range os.Args[1:] {
		command += arg + " "
	}
	return command
}

func checkForShutdown(connection net.Conn, process chan int) {
	ticker := time.NewTicker(1 * time.Second)
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-channel
	log.Printf("Caught signal, shutting down")
	pid := <-process
	log.Printf("Shutdown Process: %d\n", pid)

	proc, err := os.FindProcess(pid)
	if err != nil {
		log.Printf("Failed to find process with PID %d\n", pid)
	}
	if err = proc.Signal(os.Interrupt); err != nil {
		log.Println("Failed to stop LDT gracefully")
	}

	<-ticker.C
	os.Exit(0)
}
