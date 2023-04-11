package ldtorchestrator

import (
	"bufio"
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
	supportedLDTs   *LDTList
	otherLDTs       *LDTList
	repositories    []string
	os              string
	arch            string
}

func NewConfig(path string) *DiscoveryConfig {
	os, arch := getRuntimeInformation()
	return &DiscoveryConfig{
		repository_file: path,
		supportedLDTs:   NewLDTList(),
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
	newLDTs := github.FetchGithubReleases(c.repositories)
	for _, ldt := range newLDTs.LDTs {
		if ldt.Os == c.os && ldt.Arch == c.arch {
			c.supportedLDTs.LDTs = append(c.supportedLDTs.LDTs, ldt)
		} else {
			c.otherLDTs.LDTs = append(c.otherLDTs.LDTs, ldt)
		}
	}
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

func (c *DiscoveryConfig) updateRepositories() []string {
	var repositories []string
	file, err := os.Open(c.repository_file)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if !strings.HasPrefix(scanner.Text(), "#") {
			repositories = append(repositories, scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return repositories
}
