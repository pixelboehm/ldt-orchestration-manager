package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	comms "longevity/src/communication"
	man "longevity/src/ldt-orchestrator/manager"
	mon "longevity/src/ldt-orchestrator/monitor"
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
	monitor *mon.Monitor
}

func main() {
	defer func() {
		syscall.Unlink(socketpath)
	}()
	parseFlags()
	app := &App{
		manager: man.NewManager(repos, ldts),
		monitor: mon.NewMonitor(ldts),
	}

	if err := app.monitor.DeserializeLDTs(); err != nil {
		panic(err)
	}

	if err := app.run(os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func parseFlags() {
	flag.StringVar(&repos, "repos", "https://raw.githubusercontent.com/pixelboehm/meta-ldt/main/repositories.list", "Path to the repositories file")
	flag.StringVar(&ldts, "ldts", "/usr/local/etc/orchestration-manager/ldt.list", "Path to store LDT status")
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
	go app.monitor.DoKeepAlive()
	go app.monitor.RefreshLDTs()

	commands := make(chan string)
	for {
		in, err := listener.Accept()
		if err != nil {
			log.Fatal("error accepting connection: ", err)
			return err
		}
		go comms.GetCommandFromSocket(in, commands)
		result := app.executeCommand(<-commands, in)
		comms.SendResultToSocket(in, result)
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
		app.shutdown()
		os.Exit(1)
	case syscall.SIGHUP:
		log.Printf("Caught signal %s: reloading.", sig)
		err := syscall.Unlink(socketpath)
		if err != nil {
			log.Fatal("error during unlinking: ", err)
		}
		app.shutdown()
		if err := app.run(os.Stdout); err != nil {
			log.Fatal(err)
		}
	}
}

func (app *App) executeCommand(input string, in net.Conn) string {
	command := strings.Fields(input)

	switch command[0] {
	case "get":
		res := app.manager.GetAvailableLDTs()
		return res
	case "pull":
		if len(command) > 1 {
			id, err := strconv.Atoi(command[1])
			if err != nil {
				panic(err)
			}
			ldt := app.manager.DownloadLDT(id)
			return ldt
		}
		return " "
	case "ps":
		res := app.monitor.ListLDTs()
		return res
	case "run":
		if len(command) > 1 {
			process, err := app.manager.RunLDT(command[1])
			if err != nil {
				panic(err)
			}
			app.monitor.Started <- process
			return process.Name
		}
		return " "
	case "start":
		if len(command) > 1 {
			process, err := app.manager.StartLDT(command[1], in)
			if err != nil {
				panic(err)
			}
			app.monitor.Started <- process
			return fmt.Sprint(process.Pid)
		}
		return " "
	case "stop":
		if len(command) > 1 {
			pid, err := strconv.Atoi(command[1])
			if err != nil {
				panic(err)
			}
			res := app.manager.StopLDT(pid, true)
			return res
		}
		return " "
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

func (app *App) shutdown() {
	if err := app.monitor.SerializeLDTs(); err != nil {
		panic(err)
	}
}
