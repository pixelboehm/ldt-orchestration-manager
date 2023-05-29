package monitor

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"time"
)

func formatJSON(data json.RawMessage) string {
	formatted, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("Error formatting JSON: %v", err)
		return string(data)
	}
	return string(formatted)
}

func convertTime(started string) string {
	currentTime := time.Now().Format("2006-1-2 15:4:5")
	newCurrentTime, err := time.Parse("2006-1-2 15:4:5", currentTime)
	if err != nil {
		log.Println("Monitor: Failed to parse time")
		return "Unknown"
	}
	startTime, err := time.Parse("2006-1-2 15:4:5", started)
	if err != nil {
		log.Println("Monitor: Failed to parse time")
		return "Unknown"
	}
	uptime := newCurrentTime.Sub(startTime)
	return fmt.Sprint(uptime)
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}

func ldtIsRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		log.Println(err)
		return false
	}
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return false
	}
	return true
}

func getIPAddress() (string, error) {
	hostname, _ := os.Hostname()

	ipAddr, err := net.ResolveIPAddr("ip4", hostname)
	if err != nil {
		return "", errors.New(fmt.Sprint("Monitor: Failed wo obtain Host-IP Address"))
	}
	return ipAddr.IP.String(), nil
}
