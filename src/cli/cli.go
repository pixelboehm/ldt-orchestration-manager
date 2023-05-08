package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	socketPath = "/tmp/orchestration-manager.sock"
)

func waitForAnswer(connection net.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := connection.Read(buffer[:])
		if err != nil {
			return
		}
		val := string(buffer[0:n])
		fmt.Println(val)
		return
	}
}

func blockingWaitForAnswer(connection net.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := connection.Read(buffer[:])
		if err != nil {
			return
		}
		val := string(buffer[0:n])
		fmt.Println(val)
	}
}

func main() {
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
		go checkForShutdown(connection)
		blockingWaitForAnswer(connection)
	} else {
		waitForAnswer(connection)
	}
}

func checkForShutdown(connection net.Conn) {
	ticker := time.NewTicker(10 * time.Second)
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-channel
	log.Printf("Caught signal, shutting down")
	_, err := connection.Write([]byte("shutdown"))
	if err != nil {
		log.Fatal(err)
	}
	<-ticker.C
	os.Exit(0)
}
