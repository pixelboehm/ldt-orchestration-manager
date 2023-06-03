package manager

import (
	"fmt"
	"log"
	. "longevity/src/types"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"
)

var adjectives = []string{"joyful", "confident", "radiant", "brave", "compassionate", "creative", "enthusiastic", "energetic", "gracious", "generous", "honest", "kind", "lively", "passionate", "resourceful", "strong", "vibrant", "wise", "witty", "zealous"}

var dogs = []string{"affenpinscher", "australian_cattle_dog", "basset_hound", "bearded_collie", "bernese_mountain_dog", "border_collie", "boxer", "bulldog", "cavalier_king_charles_spaniel", "chihuahua", "dachshund", "english_cocker_spaniel", "german_shepherd_dog", "golden_retriever", "jack_russell_terrier", "labrador_retriever", "poodle", "pug", "siberian_husky", "west_highland_white_terrier"}

func prepareCommand(ldt_exec, name string, port int, device_address string) (*exec.Cmd, string) {
	makeExecutable(ldt_exec)

	cmd := exec.Command(ldt_exec, name, fmt.Sprint(port), device_address)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	return cmd, name
}

func findOpenPort() int {
	var port int

	for {
		port = generateRandomPort()
		if portIsAvailable(port) {
			break
		}
	}
	return port
}

func run(ldt_full, ldt, random_name string, port int, device_IPv4, device_MAC string) (*Process, error) {
	cmd, name := prepareCommand(ldt_full, random_name, port, device_IPv4)
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go waitOnProcess(cmd)

	process := NewProcess(cmd.Process.Pid, ldt, name, port, device_MAC)
	return process, nil
}

func start(ldt_full, ldt, random_name string, port int, device_IPv4, device_MAC string, in net.Conn) (*Process, error) {
	cmd, name := prepareCommand(ldt_full, random_name, port, device_IPv4)

	cmd.Stdout = in
	cmd.Stderr = in
	cmd.Stdin = in

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go waitOnProcess(cmd)
	process := NewProcess(cmd.Process.Pid, ldt, name, port, device_MAC)

	return process, nil
}

func stop(pid int, graceful bool) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		log.Printf("Failed to find process with PID %d\n", pid)
	}
	if graceful {
		if err = proc.Signal(os.Interrupt); err != nil {
			log.Printf("Failed to stop LDT %d gracefully\n", pid)
			return false
		}
	} else {
		if err = proc.Kill(); err != nil {
			log.Printf("Failed to kill LDT %d\n", pid)
			return false
		}
	}
	return true

}

func makeExecutable(ldt string) {
	if _, err := os.Stat(ldt); os.IsNotExist(err) {
		return
	}
	file, err := os.Open(ldt)
	if err != nil {
		panic(fmt.Sprintf("Failed to open LDT: %v\n", err))
	}
	if err := os.Chmod(file.Name(), 0755); err != nil {
		panic(fmt.Sprintf("Failed to set executable Flag: %v\n", err))
	}
}

func GenerateRandomName() string {
	rand.Seed(time.Now().UnixNano())
	return adjectives[rand.Intn(len(adjectives))] + "_" + dogs[rand.Intn((len(dogs)))]
}

func generateRandomPort() int {
	rand.Seed(time.Now().UnixNano())
	var min int = 30000
	var max int = 50000
	return rand.Intn(max-min) + min
}

func portIsAvailable(port int) bool {
	checker, err := net.Listen("tcp", ":"+fmt.Sprint(port))
	if err != nil {
		return false
	}
	_ = checker.Close()
	return true
}

func waitOnProcess(cmd *exec.Cmd) {
	err := cmd.Wait()
	if err != nil {
		log.Fatalln(err)
	}
}
