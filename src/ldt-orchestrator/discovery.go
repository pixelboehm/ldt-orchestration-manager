package ldtorchestrator

import (
	"bufio"
	"fmt"
	"log"
	"longevity/src/ldt-orchestrator/github"
	. "longevity/src/types"
	"net/url"
	"os"
	"strings"
)

type DiscoveryConfig struct {
	repository_file string
	ldtList         *LDTList
	repositories    []string
}

func NewConfig(path string) *DiscoveryConfig {
	return &DiscoveryConfig{
		repository_file: path,
		ldtList:         NewLDTList(),
		repositories:    make([]string, 0),
	}
}

func (c *DiscoveryConfig) DiscoverLDTs() {
	newLDTs := github.FetchGithubReleases(c.repositories)
	c.ldtList.LDTs = append(c.ldtList.LDTs, newLDTs.LDTs...)
	for _, ldt := range c.ldtList.LDTs {
		fmt.Printf("Name: %s \t OS: %s \t Architecture: %s\n", ldt.Name, ldt.Os, ldt.Arch)
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
