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

func waitForAnswer(connection net.Conn) {
	buffer := make([]byte, 2048)
	for {
		n, err := connection.Read(buffer[:])
		if err != nil {
			return
		}
		if n > 1 {
			val := string(buffer[0:n])
			fmt.Println(strings.Trim(val, "\n"))
		}
		return
	}
}

func blockingWaitForAnswer(connection net.Conn, process chan int) {
	buffer := make([]byte, 2048)
	var pid int = 0
	go checkIfAttachedProcessIsStillRunning(&pid)
	for {
		n, err := connection.Read(buffer[:])
		if err != nil {
			return
		}
		val := string(buffer[0:n])
		if pid == 0 {
			pid, err = strconv.Atoi(val)
			if err == nil {
				process <- pid
			}
		}
		fmt.Print(val)
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

func main() {
	process := make(chan int, 1)

	connection, err := net.Dial("unix", socketPath)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	var command string
	for _, arg := range os.Args[1:] {
		command += arg + " "
	}

	_, err = connection.Write([]byte(command))
	if err != nil {
		log.Fatal("write error:", err)
	}

	if os.Args[1] == "start" {
		go checkForShutdown(connection, process)
		blockingWaitForAnswer(connection, process)
	} else {
		waitForAnswer(connection)
	}
}

func checkForShutdown(connection net.Conn, process chan int) {
	ticker := time.NewTicker(1 * time.Second)
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-channel
	log.Printf("Caught signal, shutting down")
	pid := <-process
	log.Printf("Shutdown Process: %d\n", pid)
	if err := syscall.Kill(pid, syscall.SIGINT); err != nil {
		log.Fatal(err)
	}
	<-ticker.C
	os.Exit(0)
}
