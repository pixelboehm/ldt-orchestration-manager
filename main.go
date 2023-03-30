package main

import (
	"flag"
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

var config string

func main() {
	flag.StringVar(&config, "config", "/etc/orchestration-manager/repositories.list", "Path to the repositories file")
	flag.Parse()
	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS]", os.Args[0])
		fmt.Printf("--config \t Custom path to the repositories file")
	}
	if err := runApp(os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func runApp(out io.Writer) error {
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
		cmd := getCommandFromSocket(in)
		executeCommand(cmd)
	}
}

func getCommandFromSocket(in net.Conn) string {
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
		if err := runApp(os.Stdout); err != nil {
			log.Fatal(err)
		}
	}
}

func executeCommand(command string) {
	switch command {
	case "run":
		runManagingService()
	default:
		fmt.Println("Dont know what to do")
	}
}

func runManagingService() {
	manager := &lo.Manager{}
	manager.Setup(config)
	manager.Run()
}
