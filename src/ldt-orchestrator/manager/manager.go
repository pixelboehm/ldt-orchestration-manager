package manager

import (
	"fmt"
	"log"
	di "longevity/src/ldt-orchestrator/discovery"
	"longevity/src/ldt-orchestrator/unarchive"
	. "longevity/src/types"
	"net"
	"regexp"
)

type Manager struct {
	discovery *di.DiscoveryConfig
	storage   string
}

func NewManager(config, storage string) *Manager {
	manager := &Manager{
		discovery: di.NewConfig(config),
		storage:   storage,
	}
	return manager
}

func (manager *Manager) GetAvailableLDTs() string {
	manager.discovery.DiscoverLDTs()
	return manager.discovery.SupportedLDTs.String()
}

func (manager *Manager) GetURLFromLDTByID(id int) (string, error) {
	url, err := manager.discovery.GetUrlFromLDT(id)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (manager *Manager) GetURLFromLDTByName(name []string) (string, error) {
	url, err := manager.discovery.GetURLFromLDTByName(name)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (manager *Manager) SplitLDTInfos(name string) []string {
	reg, _ := regexp.Compile("[\\/\\:]+")
	result := reg.Split(name, -1)

	if result[2] != "latest" {
		result[2] = "v" + result[2]
	}
	return result
}

func downloadLDTArchive(address string) string {
	name, err := download(address)

	if err != nil {
		log.Println("Manager: Failed to download LDT archive: ", err)
		return ""
	}
	return name
}

func (manager *Manager) DownloadLDT(name string) string {
	manager.optionalScan()
	infos := manager.SplitLDTInfos(name)
	url, err := manager.GetURLFromLDTByName(infos)

	if err != nil {
		return err.Error()
	}

	ldtArchive := downloadLDTArchive(url)
	location := manager.storage + "/" + infos[0] + "/" + infos[1] + "/" + infos[2]
	ldt, err := unarchive.Untar(ldtArchive, location)
	if err != nil {
		log.Println("Manager: Failed to unpack LDT: ", err)
	}

	log.Printf("Manager: Downloaded LDT %s/%s:%s\n", infos[0], infos[1], infos[2])
	return ldt
}

func (manager *Manager) RunLDT(ldt string) (*Process, error) {
	process, err := run(ldt)
	if err != nil {
		log.Println("Manager: Failed to run LDT: ", err)
		return nil, err
	}

	log.Printf("Manager: Successfully started LDT with PID: %d\n", process.Pid)
	return process, nil
}

func (manager *Manager) StartLDT(ldt string, in net.Conn) (*Process, error) {
	process, err := start(ldt, in)
	if err != nil {
		log.Println("Manager: Failed to start LDT: ", err)
		return nil, err
	}

	log.Printf("Manager: Successfully started LDT with PID: %d\n", process.Pid)
	return process, nil
}

func (manager *Manager) StopLDT(pid int, graceful bool) string {
	var result string
	success := stop(pid, graceful)
	if success {
		result = fmt.Sprintf("Successfully stopped LDT with PID %d\n", pid)
	} else {
		result = fmt.Sprintf("Failed to stop LDT with PID %d\n", pid)
	}
	return result
}

func (manager *Manager) optionalScan() {
	if len(manager.discovery.SupportedLDTs.LDTs) < 1 {
		manager.GetAvailableLDTs()
	}
}
