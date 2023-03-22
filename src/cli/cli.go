package main

import (
	"log"
	"net"
	"os"
)

func main() {
	c, err := net.Dial("unix", "/tmp/orchestration-manager.sock")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	_, err = c.Write([]byte(os.Args[1]))
	if err != nil {
		log.Fatal("write error:", err)
	}
}
