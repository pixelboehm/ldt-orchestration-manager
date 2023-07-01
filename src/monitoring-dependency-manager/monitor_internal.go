package monitoring_dependency_manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	wotm "longevity/src/wot-manager"
	"net"
	"os"
	"syscall"
	"time"
)

func loadDescription(ldt string) string {
	const base string = "/usr/local/etc/orchestration-manager/"
	var desc_path string = base + ldt
	wotm, err := wotm.NewWoTmanager(desc_path)
	if err != nil {
		log.Fatal(err)
	}
	wotm_desc, err := wotm.FetchWoTDescription()
	if err != nil {
		log.Println("Monitor: Failed to fetch WoT Description: ", err)
	}
	desc, err := json.MarshalIndent(wotm_desc, "", "   ")
	if err != nil {
		log.Fatal(err)
	}
	return string(desc)
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
