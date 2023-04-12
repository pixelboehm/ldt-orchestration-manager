package discovery

import (
	"bufio"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"longevity/src/ldt-orchestrator/github"
	. "longevity/src/types"
	"net/url"
	"os"
	"runtime"
	"strings"
)

type DiscoveryConfig struct {
	repository_file string
	SupportedLDTs   *LDTList
	otherLDTs       *LDTList
	repositories    []string
	os              string
	arch            string
}

func NewConfig(path string) *DiscoveryConfig {
	os, arch := getRuntimeInformation()
	return &DiscoveryConfig{
		repository_file: path,
		SupportedLDTs:   NewLDTList(),
		otherLDTs:       NewLDTList(),
		repositories:    make([]string, 0),
		os:              os,
		arch:            arch,
	}
}

func getRuntimeInformation() (string, string) {
	os := runtime.GOOS
	arch := runtime.GOARCH
	return os, arch
}

func (c *DiscoveryConfig) DiscoverLDTs() {
	c.updateRepositories()
	newLDTs := github.FetchGithubReleases(c.repositories)
	for _, ldt := range newLDTs.LDTs {
		hash := string(createHash(&ldt))
		if ldt.Os == c.os && ldt.Arch == c.arch {
			if !hashAlreadyExists(hash, c.SupportedLDTs.LDTs) {
				c.SupportedLDTs.LDTs = append(c.SupportedLDTs.LDTs, ldt)
			}
		} else {
			if !hashAlreadyExists(hash, c.otherLDTs.LDTs) {
				c.otherLDTs.LDTs = append(c.otherLDTs.LDTs, ldt)
			}
		}
	}
}

func (c *DiscoveryConfig) GetUrlFromLDT(id int) (string, error) {
	if id >= len(c.SupportedLDTs.LDTs) {
		return "", errors.New("Failed to map ID to LDT")
	}
	return c.SupportedLDTs.LDTs[id].Url, nil

}

func createHash(l *LDT) []byte {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", l)))
	return h.Sum(nil)
}

func hashAlreadyExists(hash string, ldt []LDT) bool {
	for _, ldt := range ldt {
		if hash == string(createHash(&ldt)) {
			return true
		}
	}
	return false
}

func isGithubRepository(repo string) bool {
	stuff, _ := url.Parse(repo)
	if stuff.Host != "" {
		return strings.Contains(stuff.Host, "github.com")
	} else if strings.HasPrefix(stuff.Path, "www.github.com") {
		return true
	} else if strings.HasPrefix(stuff.Path, "github.com") {
		return true
	} else {
		return false
	}
}

func (c *DiscoveryConfig) updateRepositories() {
	file, err := os.Open(c.repository_file)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if !strings.HasPrefix(scanner.Text(), "#") {
			c.repositories = append(c.repositories, scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
