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
		fmt.Println(string(buffer[0:n]))
		return
	}
}

func main() {
	connection, err := net.Dial("unix", socketPath)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	_, err = connection.Write([]byte(os.Args[1]))
	if err != nil {
		log.Fatal("write error:", err)
	}
	waitForAnswer(connection)
}
