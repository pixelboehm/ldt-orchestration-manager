package ldtorchestrator

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func DownloadPackage(url string) (string, error) {
	file, err := os.Create("./resources/child_webserver")
	if err != nil {
		return "", err
	}
	defer file.Close()

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	if err := os.Chmod(file.Name(), 0755); err != nil {
		log.Fatalf("Could not set executable flag: %v", err)
	}

	log.Printf("Downloaded LDT: %s\n", file.Name())
	return file.Name(), nil
}

func StartPackageDetached(pkg string) {
	cmd := exec.Command("./" + pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Start()

	fmt.Printf("Started child process with PID %d\n", cmd.Process.Pid)
}
