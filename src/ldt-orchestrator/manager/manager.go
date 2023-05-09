package manager

import (
	"log"
	di "longevity/src/ldt-orchestrator/discovery"
	"longevity/src/ldt-orchestrator/unarchive"
	. "longevity/src/types"
	"net"
)

type Manager struct {
	discovery *di.DiscoveryConfig
}

func NewManager(config, ldt_list_path string) *Manager {
	manager := &Manager{
		discovery: di.NewConfig(config),
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

func downloadLDTArchive(address string) string {
	name, err := download(address)

	if err != nil {
		log.Println("Manager: Failed to download LDT archive: ", err)
		return ""
	}
	log.Printf("Manager: Downloaded LDT Archive: %s\n", name)
	return name
}

func (manager *Manager) DownloadLDT(id int) string {
	manager.optionalScan()
	url, err := manager.GetURLFromLDTByID(id)

	if err != nil {
		return err.Error()
	}

	ldtArchive := downloadLDTArchive(url)
	ldt, err := unarchive.Untar(ldtArchive, "resources")
	if err != nil {
		log.Println("Manager: Failed to unpack LDT: ", err)
	}
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

func (manager *Manager) StopLDT(pid int, graceful bool) bool {
	success := stop(pid, graceful)
	return success

}

func (manager *Manager) optionalScan() {
	if len(manager.discovery.SupportedLDTs.LDTs) < 1 {
		manager.GetAvailableLDTs()
	}
}
