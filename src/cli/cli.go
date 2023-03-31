package main

import (
	"log"
	"net"
	"os"
)

const (
	socketPath = "/tmp/orchestration-manager.sock"
)

func main() {
	c, err := net.Dial("unix", socketPath)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	_, err = c.Write([]byte(os.Args[1]))
	if err != nil {
		log.Fatal("write error:", err)
	}
}
