package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const (
	socketPath = "/tmp/orchestration-manager.sock"
)

func waitForAnswer(reader io.Reader) {
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer[:])
		if err != nil {
			return
		}
		val := string(buffer[0:n])
		fmt.Println(val)
		return
	}
}

func blockingWaitForAnswer(reader io.Reader) {
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer[:])
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
		log.Fatal(err)
	}
	defer connection.Close()

	var res string
	for _, arg := range os.Args[1:] {
		res += arg + " "
	}

	_, err = connection.Write([]byte(res))
	if err != nil {
		log.Fatal("write error:", err)
	}

	if os.Args[1] == "start" {
		blockingWaitForAnswer(connection)
	} else {
		waitForAnswer(connection)
	}
}
