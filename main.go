package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	comms "longevity/src/communication"
	man "longevity/src/ldt-orchestrator/manager"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

const (
	socketpath = "/tmp/orchestration-manager.sock"
)

var repos string
var ldts string

type App struct {
	manager *man.Manager
}

func main() {
	parseFlags()
	app := &App{
		manager: man.NewManager(repos, ldts),
	}

	if err := app.run(os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func parseFlags() {
	flag.StringVar(&repos, "repos", "/etc/orchestration-manager/repositories.list", "Path to the repositories file")
	flag.StringVar(&ldts, "ldts", "/etc/orchestration-manager/ldt.list", "Path to store LDT status")
	flag.Parse()
	flag.Usage = flagHelp
}

func flagHelp() {
	fmt.Printf("Usage: %s [OPTIONS]", os.Args[0])
	fmt.Printf("--repos \t Custom path to the repositories file")
	fmt.Printf("--ldts \t Custom path to store LDT status")
}

func (app *App) run(out io.Writer) error {
	log.SetOutput(out)
	listener, err := net.Listen("unix", socketpath)
	if err != nil {
		log.Fatal("error listening on: ", err)
		return err
	}

	sigchannel := make(chan os.Signal, 1)
	signal.Notify(sigchannel, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go app.checkForShutdown(sigchannel)

	for {
		in, err := listener.Accept()
		if err != nil {
			log.Fatal("error accepting connection: ", err)
			return err
		}
		cmd := comms.GetCommandFromSocket(in)
		res := app.executeCommand(cmd)
		comms.SendResultToSocket(in, res)
	}
}

func (app *App) checkForShutdown(c chan os.Signal) {
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
		if err := app.run(os.Stdout); err != nil {
			log.Fatal(err)
		}
	}
}

func (app *App) executeCommand(input string) string {
	command := strings.Fields(input)
	switch command[0] {
	case "get":
		res := app.manager.GetAvailableLDTs()
		return res
	case "pull":
		id, err := strconv.Atoi(command[1])
		if err != nil {
			log.Fatal(err)
		}
		ldt := app.manager.DownloadLDT(id)
		return ldt
	case "run":
		process, _ := app.manager.StartLDT(command[1])
		return process.Name
	default:
		log.Println("Unkown command received: ", command)
		fallthrough

	case "help":
		result := cliHelp()
		return result.String()
	}
}

func cliHelp() *bytes.Buffer {
	var buffer bytes.Buffer
	buffer.WriteString("Usage: cli [OPTIONS]\n")
	buffer.WriteString("help \t Show this help\n")
	buffer.WriteString("run \t Run the managing service\n")
	log.Println(buffer.String())
	return &buffer
}
