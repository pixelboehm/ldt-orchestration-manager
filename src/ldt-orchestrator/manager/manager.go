package manager

import (
	"errors"
	"fmt"
	"log"
	di "longevity/src/ldt-orchestrator/discovery"
	"longevity/src/ldt-orchestrator/unarchive"
	. "longevity/src/types"
	"net"
	"os"
	"strings"
)

type Manager struct {
	Discovery *di.Discoverer
	storage   string
	ldt_dir   string
}

func NewManager(config, storage string) *Manager {
	ldt_dir := storage + "LDTs"
	os.Mkdir(ldt_dir, 0777)
	manager := &Manager{
		Discovery: di.NewDiscoverer(config),
		storage:   storage,
		ldt_dir:   ldt_dir + "/",
	}
	return manager
}

func (manager *Manager) RunLDT(args []string) (*Process, error) {
	var ldt string = args[1]
	var ldt_path string
	var err error
	var ldt_name string
	var port int
	var device_IPv4 string
	var device_MAC string

	if len(args) > 2 {
		ldt_name = args[2]
	}
	ldt_path, ldt_name, port, device_IPv4, device_MAC, err = manager.prepareExecution(ldt, ldt_name)
	if err != nil {
		return nil, errors.New("Failed to prepare the execution")
	}

	process, err := run(ldt_path, ldt, ldt_name, port, device_IPv4, device_MAC)
	if err != nil {
		log.Println("<Manager>: Failed to run LDT: ", err)
		return nil, err
	}

	log.Printf("<Manager>: Successfully started LDT with PID: %d\n", process.Pid)
	return process, nil
}

func (manager *Manager) StartLDT(args []string, in net.Conn) (*Process, error) {
	var ldt string = args[1]
	var ldt_path string
	var err error
	var ldt_name string
	var port int
	var device_IPv4 string
	var device_MAC string

	if len(args) > 2 {
		ldt_name = args[2]
	}
	ldt_path, ldt_name, port, device_IPv4, device_MAC, err = manager.prepareExecution(ldt, ldt_name)
	if err != nil {
		return nil, errors.New("<Manager>: Failed to prepare the execution")
	}

	process, err := start(ldt_path, ldt, ldt_name, port, device_IPv4, device_MAC, in)
	if err != nil {
		log.Println("<Manager>: Failed to start LDT: ", err)
		return nil, err
	}
	log.Printf("<Manager>: Successfully started LDT with PID: %d\n", process.Pid)
	return process, nil
}

func (manager *Manager) StopLDT(pid int, name string, graceful bool) string {
	var result string
	success := stop(pid, graceful)
	if success {
		result = fmt.Sprintf("<Manager>: Successfully stopped LDT %s\n", name)
	} else {
		result = fmt.Sprintf("<Manager>: Failed to stop LDT with PID %s\n", name)
	}
	return result
}

func (manager *Manager) GetAvailableLDTs() string {
	manager.Discovery.DiscoverLDTs()
	return manager.Discovery.SupportedLDTs.String()
}

func (manager *Manager) GetURLFromLDTByID(id int) (string, error) {
	url, err := manager.Discovery.GetUrlFromLDTByID(id)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (manager *Manager) DownloadLDT(name string) (string, error) {
	manager.OptionalScan()

	if err := manager.CheckIfLdtFormatIsValid(name); err != nil {
		return "", err
	}
	user, ldt_name, version, _ := manager.splitLDTInfos(name)

	url, err := manager.getURLFromLDTByName(user, ldt_name, version)
	if err != nil {
		return "", err
	}

	ldtArchive, err := manager.downloadLDTArchive(url)
	if err != nil {
		return "", err
	}
	defer os.Remove(ldtArchive)
	location := manager.getLdtLocation(name)
	ldt, err := unarchive.Untar(ldtArchive, location)
	if err != nil {
		log.Println("<Manager>: Failed to unpack LDT: ", err)
		return "", err
	}

	if strings.HasPrefix(version, "v") {
		version = version[1:]
	}

	log.Printf("<Manager>: Downloaded LDT %s/%s:%s\n", user, ldt_name, version)
	return ldt, nil
}

func (manager *Manager) CheckIfLdtFormatIsValid(name string) error {
	_, _, _, err := manager.splitLDTInfos(name)
	if err != nil {
		return err
	}

	return nil
}

func (manager *Manager) OptionalScan() {
	manager.GetAvailableLDTs()
}

func (manager *Manager) LDTExists(ldt string) bool {
	ldt_path := manager.getLdtLocation(ldt)
	if _, err := os.Stat(ldt_path); err != nil {
		return false
	}
	return true
}
