package manager

import (
	"errors"
	"fmt"
	"io"
	"log"
	wot "longevity/src/wot-manager"
	"os"
	"regexp"
)

func (manager *Manager) prepareExecution(ldt, name string) (string, string, int, string, string, error) {
	var dest string = ""
	var known_ldt bool
	var dir string = ""
	var port int
	var err error
	var device_ipv4 string
	var device_mac string
	_, ldt_name, _, _ := manager.splitLDTInfos(ldt)
	if name != "" {
		dest = manager.storage + name
		if _, err := os.Stat(dest); err == nil {
			known_ldt = true
			log.Println("<Manager>: Starting Known LDT: ", name)
			dir = dest
			wotm, err := wot.NewWoTmanager(dir)
			port = wotm.GetLdtPortFromDescription()
			if port == 0 {
				port = generateRandomPort()
			}
			device_ipv4 = wotm.GetDeviceIPv4AddressFromDescription()
			device_mac = wotm.GetDeviceMACAddressFromDescription()
			if err != nil {
				log.Println("<Manager>: Err: ", err)
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

		log.Println("<Manager>: Starting unknown LDT: ", name)
		dir, err = createLdtSpecificDirectory(dest)
		if err != nil {
			return "", "", -1, "", "", err
		}
		manager.copyLdtDescription(ldt, dir)
		port = generateRandomPort()
	}
	src_exec := manager.getLdtLocation(ldt) + "/" + ldt_name
	dest_exec := dir + "/" + ldt_name
	err = symlinkLdtExecutable(src_exec, dest_exec)
	if err != nil {
		log.Printf("Symlinking failed: %s\n", err)
		return "", "", -1, "", "", errors.New(fmt.Sprint("<Manager>: Unable to symlink LDT", err))
	}
	return dest_exec, name, port, device_ipv4, device_mac, nil
}

func (manager *Manager) copyLdtDescription(ldt, dest string) error {
	src_dir := manager.getLdtLocation(ldt)
	src_description := src_dir + "/" + "wotm/description.json"
	dest_description := dest + "/" + "wotm/description.json"
	os.MkdirAll(dest+"/"+"wotm", 0777)
	err := copyFile(src_description, dest_description)
	if err != nil {
		return err
	}
	return nil
}

func (manager *Manager) getLdtLocation(ldt string) string {
	user, ldt_name, version, _ := manager.splitLDTInfos(ldt)
	return manager.ldt_dir + user + "/" + ldt_name + "/" + version
}

func (manager *Manager) downloadLDTArchive(address string) (string, error) {
	name, err := download(address, manager.ldt_dir)

	if err != nil {
		log.Println("<Manager>: Failed to download LDT archive: ", err)
		return "", err
	}
	return name, nil
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
	if _, err := os.Lstat(dest); err == nil {
		os.Remove(dest)
	}
	err := os.Symlink(src, dest)
	if err != nil {
		return err
	}
	return nil
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

func copyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		return fmt.Errorf("<Manager>: CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("<Manager>: CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	err = copyFileContents(src, dst)
	return
}

func (manager *Manager) splitLDTInfos(name string) (string, string, string, error) {
	reg, _ := regexp.Compile("[\\/\\:]+")
	result := reg.Split(name, -1)

	if len(result) != 3 {
		return "", "", "", errors.New("Invalid LDT format specified")
	}

	if result[2] != "latest" {
		result[2] = "v" + result[2]
	}

	return result[0], result[1], result[2], nil
}

func (manager *Manager) getURLFromLDTByName(user, ldt, version string) (string, error) {
	url, err := manager.Discovery.GetURLFromLDTByName(user, ldt, version)
	if err != nil {
		return "", err
	}

	return url, nil
}
