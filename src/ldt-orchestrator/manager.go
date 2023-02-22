package ldtorchestrator

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/cavaliergopher/grab/v3"
)

func DownloadPackage(url string) (string, error) {
	// create client
	client := grab.NewClient()
	req, _ := grab.NewRequest(".", url)

	// start download
	log.Printf("Downloading %v...\n", req.URL())
	resp := client.Do(req)
	log.Printf("  %v\n", resp.HTTPResponse.Status)

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		return resp.Filename, err
	}

	fmt.Printf("Download saved to ./%v \n", resp.Filename)

	if err := os.Chmod(resp.Filename, 0777); err != nil {
		log.Fatal(err)
	}
	return resp.Filename, nil
}

func StartPackageDetached(pkg string) {
	cmd := exec.Command("./" + pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Start()

	fmt.Printf("Started child process with PID %d\n", cmd.Process.Pid)
}
