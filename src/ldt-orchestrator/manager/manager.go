package manager

import (
	"errors"
	"fmt"
	"io"
	"log"
	di "longevity/src/ldt-orchestrator/discovery"
	"longevity/src/ldt-orchestrator/unarchive"
	. "longevity/src/types"
	"net"
	"os"
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

func (manager *Manager) GetURLFromLDTByName(user, ldt, version string) (string, error) {
	url, err := manager.discovery.GetURLFromLDTByName(user, ldt, version)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (manager *Manager) SplitLDTInfos(name string) (string, string, string) {
	reg, _ := regexp.Compile("[\\/\\:]+")
	result := reg.Split(name, -1)

	if result[2] != "latest" {
		result[2] = "v" + result[2]
	}
	return result[0], result[1], result[2]
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
	user, ldt_name, version := manager.SplitLDTInfos(name)
	url, err := manager.GetURLFromLDTByName(user, ldt_name, version)

	if err != nil {
		return err.Error()
	}

	ldtArchive := downloadLDTArchive(url)
	location := manager.storage + "/" + user + "/" + ldt_name + "/" + version
	ldt, err := unarchive.Untar(ldtArchive, location)
	if err != nil {
		log.Println("Manager: Failed to unpack LDT: ", err)
		return ""
	}

	log.Printf("Manager: Downloaded LDT %s/%s:%s\n", user, ldt_name, version)
	return ldt
}

func (manager *Manager) prepareExecution(ldt string) (string, string, error) {
	user, ldt_name, version := manager.SplitLDTInfos(ldt)
	random_name := GenerateRandomName()

	dest := manager.storage + "/" + random_name
	dir, err := createLdtSpecificDirectory(dest)
	if err != nil {
		return "", "", errors.New(fmt.Sprint("Could not create LDT specific directory", err))
	}

	manager.copyLdtDescription(ldt, dir)

	src_exec := manager.storage + "/" + user + "/" + ldt_name + "/" + version + "/" + ldt_name
	dest_exec := dir + "/" + ldt_name
	err = symlinkLdtExecutable(src_exec, dest_exec)
	if err != nil {
		log.Println()
		return "", "", errors.New(fmt.Sprint("Unable to symlink LDT", err))
	}
	return dest_exec, random_name, nil
}

func createLdtSpecificDirectory(dir string) (string, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return "", err
		}
	}
	return dir, nil
}

func (manager *Manager) copyLdtDescription(ldt, dest string) error {
	src_dir := manager.getLdtLocation(ldt)
	src_description := src_dir + "/" + "wotm/description.json"
	dest_description := dest + "/" + "wotm/description.json"
	os.MkdirAll(dest+"/"+"wotm", 0777)
	log.Printf("Source: %s", src_description)
	log.Printf("Dest: %s", dest_description)
	err := CopyFile(src_description, dest_description)
	if err != nil {
		return err
	}
	return nil
}

func symlinkLdtExecutable(src, dest string) error {
	log.Println("symlinking")
	log.Printf("Source: %s\n", src)
	log.Printf("Dest: %s\n", dest)
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
	if err = os.Link(src, dst); err == nil {
		return
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

func (manager *Manager) RunLDT(ldt string) (*Process, error) {
	ldt_path, random_name, err := manager.prepareExecution(ldt)
	if err != nil {
		return nil, errors.New("Failed to prepare the execution")
	}
	process, err := run(ldt_path, ldt, random_name)
	if err != nil {
		log.Println("Manager: Failed to run LDT: ", err)
		return nil, err
	}

	log.Printf("Manager: Successfully started LDT with PID: %d\n", process.Pid)
	return process, nil
}

func (manager *Manager) StartLDT(ldt string, in net.Conn) (*Process, error) {
	ldt_path, random_name, err := manager.prepareExecution(ldt)
	if err != nil {
		return nil, errors.New("Failed to prepare the execution")
	}
	process, err := start(ldt_path, ldt, random_name, in)
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

func (manager *Manager) getLdtLocation(ldt string) string {
	user, ldt_name, version := manager.SplitLDTInfos(ldt)
	return manager.storage + "/" + user + "/" + ldt_name + "/" + version
}
