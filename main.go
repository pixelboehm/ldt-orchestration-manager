package main

import (
	"fmt"
	"io"
	"log"
	lo "longevity/src/ldt-orchestrator"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	socketpath = "/tmp/orchestration-manager.sock"
)

func main() {
	if err := run(os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func run(out io.Writer) error {
	log.SetOutput(out)
	listener, err := net.Listen("unix", socketpath)
	if err != nil {
		log.Fatal("error listening on: ", err)
		return err
	}

	sigchannel := make(chan os.Signal, 1)
	signal.Notify(sigchannel, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go checkForShutdown(sigchannel)

	for {
		in, err := listener.Accept()
		if err != nil {
			log.Fatal("error accepting connection: ", err)
			return err
		}
		cmd := getCommand(in)
		executeCommand(cmd)
	}
}

func runMonitor() {
	manager := &lo.Manager{Monitor: lo.NewMonitor()}
	manager.Run()
}

func getCommand(in net.Conn) string {
	for {
		buf := make([]byte, 512)
		nr, err := in.Read(buf)
		if err != nil {
			return "help"
		}

		data := buf[0:nr]
		return string(data)
	}
}

func executeCommand(command string) {
	switch command {
	case "run":
		runMonitor()
	default:
		fmt.Println("Dont know what to do")
	}
}

func checkForShutdown(c chan os.Signal) {
	sig := <-c
	switch sig {
	case syscall.SIGINT, syscall.SIGTERM:
		log.Printf("Caught signal %s: shutting down.", sig)
		err := syscall.Unlink(socketpath)
		if err != nil {
			log.Fatal("error during unlinking: ", err)
		}
		os.Exit(1)
	case syscall.SIGHUP:
		log.Printf("Caught signal %s: reloading.", sig)
		err := syscall.Unlink(socketpath)
		if err != nil {
			log.Fatal("error during unlinking: ", err)
		}
		if err := run(os.Stdout); err != nil {
			log.Fatal(err)
		}
	}
}
