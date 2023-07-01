package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	boo "longevity/src/bootstrapper"
	comms "longevity/src/communication"
	man "longevity/src/ldt-orchestrator/manager"
	mon "longevity/src/monitoring-dependency-manager"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
)

var repos string
var ldts string
var storage string
var socket string

type App struct {
	manager      *man.Manager
	monitor      *mon.Monitor
	bootstrapper *boo.Bootstrapper
}

func main() {
	initialize()
	defer func() {
		syscall.Unlink(socket)
	}()

	var monitor *mon.Monitor = mon.NewMonitor(ldts)
	app := &App{
		manager:      man.NewManager(repos, storage),
		monitor:      monitor,
		bootstrapper: boo.NewBootstrapper(monitor),
	}

	if err := app.monitor.DeserializeLDTs(); err != nil {
		panic(err)
	}

	if err := app.run(os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func initialize() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Main: Failed to load .env file")
	}
	socket = os.Getenv("SOCKET")
	repos = os.Getenv("META_REPOSITORY")
	storage = os.Getenv("ODM_DATA_DIRECTORY")
	if storage[len(storage)-1:] != "/" {
		storage = storage + "/"
	}
	ldts = storage + "ldt.list"
	parseFlags()
}

func parseFlags() {
	flag.StringVar(&repos, "repos", repos, "Path to the meta repositories file")
	flag.StringVar(&storage, "data-dir", storage, "Path to the ODM data directory")
	flag.Parse()
}

func flagHelp() {
	fmt.Printf("Usage: %s [OPTIONS]", os.Args[0])
	fmt.Printf("-repos \t Custom path to the repositories file")
	fmt.Printf("-data-dir \t Custom path the the ODM data directory")
}

func (app *App) run(out io.Writer) error {
	log.SetOutput(out)
	listener, err := net.Listen("unix", socket)
	if err != nil {
		log.Fatal("error listening on: ", err)
		return err
	}

	sigchannel := make(chan os.Signal, 1)
	signal.Notify(sigchannel, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go app.checkForShutdown(sigchannel)
	go app.monitor.DoKeepAlive()
	go app.monitor.RefreshLDTs()
	go app.monitor.Run(8080)
	go app.bootstrapper.Run(55443)

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

func (app *App) executeCommand(input string, in net.Conn) string {
	command := strings.Fields(input)

	switch command[0] {
	case "get":
		res := app.manager.GetAvailableLDTs()
		return res
	case "pull":
		if len(command) > 1 {
			ldt := app.manager.DownloadLDT(command[1])
			return ldt
		}
		return " "
	case "ps":
		res := app.monitor.ListLDTs()
		return res
	case "run":
		if len(command) > 1 {
			process, err := app.manager.RunLDT(command)
			if err != nil {
				panic(err)
			}
			app.monitor.Started <- process
			return process.Name
		}
		return " "
	case "start":
		if len(command) > 1 {
			process, err := app.manager.StartLDT(command, in)
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
	case "kill":
		if len(command) > 1 {
			pid, err := strconv.Atoi(command[1])
			if err != nil {
				panic(err)
			}
			res := app.manager.StopLDT(pid, false)
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

func (app *App) checkForShutdown(c chan os.Signal) {
	sig := <-c
	switch sig {
	case syscall.SIGINT, syscall.SIGTERM:
		log.Printf("Caught signal %s: shutting down.", sig)
		err := syscall.Unlink(socket)
		if err != nil {
			log.Fatal("error during unlinking: ", err)
		}
		app.shutdown()
		os.Exit(1)
	case syscall.SIGHUP:
		log.Printf("Caught signal %s: reloading.", sig)
		err := syscall.Unlink(socket)
		if err != nil {
			log.Fatal("error during unlinking: ", err)
		}
		app.shutdown()
		if err := app.run(os.Stdout); err != nil {
			log.Fatal(err)
		}
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
