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

func prepareCommand(ldt, name string) (*exec.Cmd, string) {
	makeExecutable(ldt)

	if name == "" {
		name = generateRandomName()
	}

	cmd := exec.Command(ldt, name)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	return cmd, name
}

func run(ldt_full, ldt string) (*Process, error) {
	cmd, name := prepareCommand(ldt_full, "")
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go waitOnProcess(cmd)

	process := NewProcess(cmd.Process.Pid, ldt, name)
	return process, nil
}

func start(ldt_full, ldt string, in net.Conn) (*Process, error) {
	cmd, name := prepareCommand(ldt_full, "")
	cmd.Stdout = in
	cmd.Stderr = in
	cmd.Stdin = in

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go waitOnProcess(cmd)

	process := NewProcess(cmd.Process.Pid, ldt, name)

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

func generateRandomName() string {
	rand.Seed(time.Now().UnixNano())
	return adjectives[rand.Intn(len(adjectives))] + "_" + dogs[rand.Intn((len(dogs)))]
}

func waitOnProcess(cmd *exec.Cmd) {
	err := cmd.Wait()
	if err != nil {
		log.Fatalln(err)
	}
}
