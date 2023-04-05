package main

import (
	"bytes"
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

var repos string
var ldts string

func main() {
	flag.StringVar(&repos, "repos", "/etc/orchestration-manager/repositories.list", "Path to the repositories file")
	flag.StringVar(&ldts, "ldts", "/etc/orchestration-manager/ldt.list", "Path to store LDT status")
	flag.Parse()
	flag.Usage = help
	if err := runApp(os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func help() {
	fmt.Printf("Usage: %s [OPTIONS]", os.Args[0])
	fmt.Printf("--repos \t Custom path to the repositories file")
	fmt.Printf("--ldts \t Custom path to store LDT status")
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
		res := executeCommand(cmd)
		sendResultToSocket(in, res)
	}
}

func sendResultToSocket(out net.Conn, res string) {
	_, err := out.Write([]byte(res))
	if err != nil {
		log.Fatal("error writing to socket: ", err)
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

func executeCommand(command string) string {
	switch command {
	case "run":
		runManagingService()
		return "Running managing service"
	default:
		log.Println("Unkown command received: ", command)
		fallthrough
	case "help":
		result := cliHelp()
		return result.String()
	}
}

func runManagingService() {
	manager := &lo.Manager{}
	manager.Setup(repos, ldts)
	manager.Run()
}

func cliHelp() *bytes.Buffer {
	var buffer bytes.Buffer
	buffer.WriteString("Usage: cli [OPTIONS]\n")
	buffer.WriteString("help \t Show this help\n")
	buffer.WriteString("run \t Run the managing service\n")
	log.Println(buffer.String())
	return &buffer
}
