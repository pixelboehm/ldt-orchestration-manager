package manager

import (
	"log"
	di "longevity/src/ldt-orchestrator/discovery"
	mo "longevity/src/ldt-orchestrator/monitor"
	"longevity/src/ldt-orchestrator/unarchive"
	. "longevity/src/types"
)

type Manager struct {
	monitor   *mo.Monitor
	discovery *di.DiscoveryConfig
}

func NewManager(config, ldt_list_path string) *Manager {
	manager := &Manager{
		monitor:   mo.NewMonitor(ldt_list_path),
		discovery: di.NewConfig(config),
	}

	if err := manager.monitor.DeserializeLDTs(); err != nil {
		log.Fatal(err)
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
		log.Println("Failed to download LDT archive: ", err)
		return ""
	}
	log.Printf("Downloaded LDT Archive: %s\n", name)
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
		log.Println("Failed to unpack LDT: ", err)
	}
	return ldt
}

func (manager *Manager) StartLDT(ldt string) (*Process, error) {
	process, err := start(ldt)
	if err != nil {
		log.Fatal("Failed to start LDT: ", err)
		return nil, err
	}

	log.Printf("Successfully started LDT with PID: %d\n", process.Pid)
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

func (manager *Manager) shutdown() {
	if err := manager.monitor.SerializeLDTs(); err != nil {
		log.Fatal(err)
	}
}
