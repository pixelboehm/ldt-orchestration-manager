package ldtorchestrator

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync"
)

type LDTList struct {
	Ldts []LDT
	lock sync.Mutex
}

type LDT struct {
	Name    string
	Version string
	Os      string
	Arch    string
	Url     string
}

type DiscoveryConfig struct {
	repository_file string
	ldtList         *LDTList
	repositories    []string
}

func NewLDTList() *LDTList {
	return &LDTList{
		Ldts: nil,
		lock: sync.Mutex{},
	}
}

func NewConfig(path string) *DiscoveryConfig {
	return &DiscoveryConfig{
		repository_file: path,
		ldtList:         NewLDTList(),
		repositories:    make([]string, 0),
	}
}

// func (c *DiscoveryConfig) FetchGithubReleases() {
// 	c.repositories = c.updateRepositories()
// 	for _, repo := range c.repositories {
// 		if isGithubRepository(repo) {
// 			owner, repo := parseRepository(repo)
// 			_, err := github.GetReleasesFromRepository(owner, repo)
// 			if err != nil {
// 				log.Println(err)
// 			}
// 			// for _, release := range releases {
// 			// 	c.ldtList.lock.Lock()
// 			// 	c.filterLDTs(release)
// 			// 	c.ldtList.lock.Unlock()
// 			// }
// 		}
// 	}
// }

func parseRepository(repo string) (string, string) {
	split := strings.Split(repo, "/")
	return split[3], split[4]
}

func isGithubRepository(repo string) bool {
	return strings.HasPrefix(repo, "https://github.com")
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
