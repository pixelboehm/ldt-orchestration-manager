package main

import (
	"bytes"
	"errors"
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
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

const (
	socket = "/tmp/orchestration-manager.sock"
)

var repos string
var ldts string
var storage string
var env string = ".env"

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
	var manager *man.Manager = man.NewManager(repos, storage)
	app := &App{
		manager:      manager,
		monitor:      monitor,
		bootstrapper: boo.NewBootstrapper(monitor, manager),
	}

	if err := app.monitor.DeserializeLDTs(); err != nil {
		panic(err)
	}

	if err := app.run(os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func initialize() {
	parseFlags()
	if err := godotenv.Load(env, "./src/ldt-orchestrator/github/github.env"); err != nil {
		log.Fatal("Main: Failed to load .env files")
	}
	if repos == "" {
		repos = os.Getenv("META_REPOSITORY")
	}
	if storage == "" {
		storage = os.Getenv("ODM_DATA_DIRECTORY")
	}
	if storage[len(storage)-1:] != "/" {
		storage = storage + "/"
	}
	ldts = storage + "ldt.list"
}

func parseFlags() {
	flag.StringVar(&repos, "repos", repos, "Meta repositories file")
	flag.StringVar(&storage, "data-dir", storage, "ODM data directory")
	flag.StringVar(&env, "env", env, ".env variable")
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
	go app.monitor.DoKeepAlive(5)
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
		app.executeCommand(<-commands, in)
	}
}

func (app *App) executeCommand(input string, in net.Conn) {
	command := strings.Fields(input)

	switch command[0] {
	case "get":
		result := app.manager.GetAvailableLDTs()
		comms.SendResultToSocket(in, result)
	case "kill":
		if len(command) > 1 {
			name := command[1]
			pid, err := app.monitor.GetPidViaLdtName(name)
			if err != nil {
				result := fmt.Sprintf("LDT %s does not exist\n", name)
				comms.SendResultToSocket(in, result)
				return
			}
			result := app.manager.StopLDT(pid, name, false)
			comms.SendResultToSocket(in, result)
			return
		}
		comms.SendResultToSocket(in, " ")
	case "pull":
		if len(command) > 1 {
			ldt_name := command[1]
			ldt, err := app.manager.DownloadLDT(ldt_name)
			if err != nil {
				result := fmt.Sprintf("Failed to Download LDT %s: %v\n", ldt_name, err)
				comms.SendResultToSocket(in, result)
				return
			}
			comms.SendResultToSocket(in, ldt)
			return
		}
		comms.SendResultToSocket(in, " ")
	case "ps":
		result := app.monitor.ListLDTs()
		comms.SendResultToSocket(in, result)
	case "run":
		if len(command) > 1 {
			if err := app.manager.CheckIfLdtFormatIsValid(command[1]); err != nil {
				comms.SendResultToSocket(in, err.Error())
				return
			}
			var ldt_name string = command[1]
			if !app.manager.LDTExists(ldt_name) {
				ticker := time.NewTicker(2 * time.Second)
				_ = app.manager.GetAvailableLDTs()
				_, err := app.manager.DownloadLDT(ldt_name)
				if err != nil {
					result := fmt.Sprintf("Failed to Download LDT: %s\n", ldt_name)
					comms.SendResultToSocket(in, result)
					return
				}
				<-ticker.C
			}
			process, err := app.manager.RunLDT(command)
			if err != nil {
				panic(err)
			}
			app.monitor.Started <- process
			result := process.Name
			comms.SendResultToSocket(in, result)
			return
		}
		comms.SendResultToSocket(in, " ")
	case "show":
		if len(command) > 1 {
			path := storage + command[1] + "/wotm/description.json"
			if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
				comms.SendResultToSocket(in, "LDT does not exist")
				return
			}
			comms.SendResultToSocket(in, path)
			return
		}
		comms.SendResultToSocket(in, " ")
	case "start":
		if len(command) > 1 {
			var ldt_name string = command[1]
			if !app.manager.LDTExists(ldt_name) {
				ticker := time.NewTicker(2 * time.Second)
				_ = app.manager.GetAvailableLDTs()
				_, err := app.manager.DownloadLDT(ldt_name)
				if err != nil {
					result := fmt.Sprintf("Failed to Download LDT: %s\n", ldt_name)
					comms.SendResultToSocket(in, result)
					return
				}
				<-ticker.C
			}
			process, err := app.manager.StartLDT(command, in)
			if err != nil {
				panic(err)
			}
			app.monitor.Started <- process
			result := fmt.Sprint(process.Pid)
			comms.SendResultToSocket(in, result)
			return
		}
		comms.SendResultToSocket(in, " ")
	case "stop":
		if len(command) > 1 {
			name := command[1]
			pid, err := app.monitor.GetPidViaLdtName(name)
			if err != nil {
				result := fmt.Sprintf("LDT %s does not exist\n", name)
				comms.SendResultToSocket(in, result)
				return
			}
			result := app.manager.StopLDT(pid, name, true)
			comms.SendResultToSocket(in, result)
			return
		}
		comms.SendResultToSocket(in, " ")
	case "rm":
		if len(command) > 1 {
			ldt := command[1]
			path := storage + command[1]
			if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
				comms.SendResultToSocket(in, "LDT does not exist")
				return
			}
			var result string
			if err := os.RemoveAll(path); err != nil {
				result = fmt.Sprintf("Failed to remove LDT: %s\n", ldt)
				comms.SendResultToSocket(in, result)
				return
			}
			result = fmt.Sprintf("Successfully removed LDT: %s\n", ldt)
			comms.SendResultToSocket(in, result)
			return
		}
		comms.SendResultToSocket(in, " ")
	default:
		log.Println("Unkown command received: ", command)
		fallthrough
	case "help":
		result := cliHelp().String()
		comms.SendResultToSocket(in, result)
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
	buffer.WriteString("help\t\t\t\t\t\tShow this help\n")
	buffer.WriteString("get\t\t\t\t\t\tget all available LDTs\n")
	buffer.WriteString("kill\t<ldt-identifier>\t\t\tungraceful stop specified LDT\n")
	buffer.WriteString("pull\t<ldt-name>\t\t\t\tdownload specified LDT\n")
	buffer.WriteString("ps\t\t\t\t\t\tlist all running LDTs\n")
	buffer.WriteString("run\t<ldt-name>\t[ldt-identifier]\tdetached start LDT\n")
	buffer.WriteString("show\t<ldt-identifier>\t\t\tdisplays Web-of-Things description\n")
	buffer.WriteString("start\t<ldt-name>\t[ldt-identifier]\tattached start LDT\n")
	buffer.WriteString("stop\t<ldt-identifier>\t\t\tgraceful stop specified LDT \n")
	buffer.WriteString("rm\t<ldt-identifier>\t\t\tremoves LDT cache \n")
	return &buffer
}

func (app *App) shutdown() {
	if err := app.monitor.SerializeLDTs(); err != nil {
		panic(err)
	}
}
