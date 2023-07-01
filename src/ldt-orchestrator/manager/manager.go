package manager

import (
	"errors"
	"fmt"
	"io"
	"log"
	di "longevity/src/ldt-orchestrator/discovery"
	"longevity/src/ldt-orchestrator/unarchive"
	. "longevity/src/types"
	wot "longevity/src/wot-manager"
	"net"
	"os"
	"regexp"
)

type Manager struct {
	discovery *di.Discoverer
	storage   string
}

func NewManager(config, storage string) *Manager {
	manager := &Manager{
		discovery: di.NewDiscoverer(config),
		storage:   storage,
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
		log.Println("Manager: Failed to run LDT: ", err)
		return nil, err
	}

	log.Printf("Manager: Successfully started LDT with PID: %d\n", process.Pid)
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
		return nil, errors.New("Failed to prepare the execution")
	}

	process, err := start(ldt_path, ldt, ldt_name, port, device_IPv4, device_MAC, in)
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

func (manager *Manager) GetAvailableLDTs() string {
	manager.discovery.DiscoverLDTs()
	return manager.discovery.SupportedLDTs.String()
}

func (manager *Manager) GetURLFromLDTByID(id int) (string, error) {
	url, err := manager.discovery.GetUrlFromLDTByID(id)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (manager *Manager) GetURLFromLDTByName(user, ldt, version string) (string, error) {
	url, err := manager.discovery.GetURLFromLDTByName(user, ldt, version)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (manager *Manager) DownloadLDT(name string) string {
	manager.optionalScan()
	user, ldt_name, version := manager.SplitLDTInfos(name)
	url, err := manager.GetURLFromLDTByName(user, ldt_name, version)
	if err != nil {
		return err.Error()
	}

	ldtArchive := manager.downloadLDTArchive(url)
	defer os.Remove(ldtArchive)
	location := manager.storage + user + "/" + ldt_name + "/" + version
	ldt, err := unarchive.Untar(ldtArchive, location)
	if err != nil {
		log.Println("Manager: Failed to unpack LDT: ", err)
		return ""
	}

	log.Printf("Manager: Downloaded LDT %s/%s:%s\n", user, ldt_name, version)
	return ldt
}

func (manager *Manager) SplitLDTInfos(name string) (string, string, string) {
	reg, _ := regexp.Compile("[\\/\\:]+")
	result := reg.Split(name, -1)

	if result[2] != "latest" {
		result[2] = "v" + result[2]
	}
	return result[0], result[1], result[2]
}

func (manager *Manager) prepareExecution(ldt, name string) (string, string, int, string, string, error) {
	user, ldt_name, version := manager.SplitLDTInfos(ldt)
	var dest string = ""
	var known_ldt bool
	var dir string = ""
	var port int
	var err error
	var device_ipv4 string
	var device_mac string
	if name != "" {
		dest = manager.storage + name
		if _, err := os.Stat(dest); err == nil {
			known_ldt = true
			log.Println("Manager: Starting Known LDT: ", name)
			dir = dest
			wotm, err := wot.NewWoTmanager(dir)
			port = wotm.GetLdtPortFromDescription()
			if port == 0 {
				port = generateRandomPort()
			}
			device_ipv4 = wotm.GetDeviceIPv4AddressFromDescription()
			device_mac = wotm.GetDeviceMACAddressFromDescription()
			if err != nil {
				log.Println("Manager: Err: ", err)
			}
		}
	}
	if !known_ldt {
		if name == "" {
			name = GenerateRandomName()
		}
		for dest == "" {
			dest = manager.storage + name
			if _, err := os.Stat(dest); err == nil {
				dest = ""
				name = GenerateRandomName()
			}
		}

		log.Println("Starting unknown LDT: ", name)
		dir, err = createLdtSpecificDirectory(dest)
		if err != nil {
			return "", "", -1, "", "", err
		}
		manager.copyLdtDescription(ldt, dir)
		port = generateRandomPort()
	}
	src_exec := manager.storage + user + "/" + ldt_name + "/" + version + "/" + ldt_name
	dest_exec := dir + "/" + ldt_name
	if !known_ldt {
		err = symlinkLdtExecutable(src_exec, dest_exec)
		if err != nil {
			log.Println()
			return "", "", -1, "", "", errors.New(fmt.Sprint("Unable to symlink LDT", err))
		}
	}
	return dest_exec, name, port, device_ipv4, device_mac, nil
}

func (manager *Manager) copyLdtDescription(ldt, dest string) error {
	src_dir := manager.getLdtLocation(ldt)
	src_description := src_dir + "/" + "wotm/description.json"
	dest_description := dest + "/" + "wotm/description.json"
	os.MkdirAll(dest+"/"+"wotm", 0777)
	err := CopyFile(src_description, dest_description)
	if err != nil {
		return err
	}
	return nil
}

func (manager *Manager) optionalScan() {
	if len(manager.discovery.SupportedLDTs.LDTs) < 1 {
		manager.GetAvailableLDTs()
	}
}

func (manager *Manager) getLdtLocation(ldt string) string {
	user, ldt_name, version := manager.SplitLDTInfos(ldt)
	return manager.storage + "/" + user + "/" + ldt_name + "/" + version
}

func (manager *Manager) downloadLDTArchive(address string) string {
	name, err := download(address, manager.storage)

	if err != nil {
		log.Println("Manager: Failed to download LDT archive: ", err)
		return ""
	}
	return name
}

func createLdtSpecificDirectory(dir string) (string, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return "", err
		}
	}
	return dir, nil
}

func symlinkLdtExecutable(src, dest string) error {
	log.Println("Manager: Symlinking Executable")
	err := os.Symlink(src, dest)
	if err != nil {
		return err
	}
	return nil
}

func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	err = copyFileContents(src, dst)
	return
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
