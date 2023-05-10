package discovery

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"longevity/src/ldt-orchestrator/github"
	"longevity/src/types"
	. "longevity/src/types"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"sort"
	"strings"
)

type DiscoveryConfig struct {
	repository_source string
	SupportedLDTs     *LDTList
	otherLDTs         *LDTList
	repositories      []string
	os                string
	arch              string
}

func NewConfig(path string) *DiscoveryConfig {
	os, arch := getRuntimeInformation()
	return &DiscoveryConfig{
		repository_source: path,
		SupportedLDTs:     NewLDTList(),
		otherLDTs:         NewLDTList(),
		repositories:      make([]string, 0),
		os:                os,
		arch:              arch,
	}
}

func getRuntimeInformation() (string, string) {
	os := runtime.GOOS
	arch := runtime.GOARCH
	return os, arch
}

func (c *DiscoveryConfig) DiscoverLDTs() {
	_ = c.updateRepositories()
	newLDTs := github.FetchGithubReleases(c.repositories)
	var update_latest_tag_supported bool = false
	var update_latest_tag_other bool = false
	for _, ldt := range newLDTs.LDTs {
		if ldt.Os == c.os && ldt.Arch == c.arch {
			if !ldtAlreadyExists(&ldt, c.SupportedLDTs) {
				c.SupportedLDTs.LDTs = append(c.SupportedLDTs.LDTs, ldt)
				update_latest_tag_supported = true
			}
		} else {
			if !ldtAlreadyExists(&ldt, c.otherLDTs) {
				c.otherLDTs.LDTs = append(c.otherLDTs.LDTs, ldt)
				update_latest_tag_supported = true
			}
		}
	}
	if update_latest_tag_supported {
		updateLatestTag(&c.SupportedLDTs.LDTs)
	}
	if update_latest_tag_other {
		updateLatestTag(&c.otherLDTs.LDTs)
	}
}

func updateLatestTag(ldts *[]LDT) {
	sortLDTsByName(*ldts)

	var last_ldt_name string
	var current_latest_version string
	var latest_ldt_changed bool = false
	var latest_ldt *LDT

	for _, ldt := range *ldts {
		if ldt.Name != last_ldt_name {
			last_ldt_name = ldt.Name
			current_latest_version = ldt.Version
			latest_ldt = &ldt
			latest_ldt_changed = true
		} else if ldt.Name == last_ldt_name {
			latest_ldt_changed = false
			if ldt.Version > current_latest_version {
				current_latest_version = ldt.Version
				latest_ldt = &ldt
				latest_ldt_changed = true
			}
		}
		if latest_ldt_changed {
			latest_ldt.Version = "latest"
			*ldts = append([]types.LDT{*latest_ldt}, *ldts...)
		}
	}
}

func sortLDTsByName(ldts []LDT) {
	sort.Slice(ldts, func(i, j int) bool {
		return ldts[i].Name > ldts[j].Name
	})
}

func (c *DiscoveryConfig) GetUrlFromLDT(id int) (string, error) {
	if id >= len(c.SupportedLDTs.LDTs) {
		return "", errors.New("Failed to map ID to LDT")
	}
	return c.SupportedLDTs.LDTs[id].Url, nil
}

func (c *DiscoveryConfig) GetURLFromLDTByName(val []string) (string, error) {
	var ldt string = val[1]
	var tag string = val[2]
	for _, entry := range c.SupportedLDTs.LDTs {
		if entry.Name == ldt && entry.Version == tag {
			return entry.Url, nil
		}
	}
	return "", errors.New(fmt.Sprintf("Unable to find LDT: %s", val[1]))
}

func ldtAlreadyExists(ldt *LDT, ldt_list *LDTList) bool {
	for _, existingLDT := range ldt_list.LDTs {
		if string(ldt.Hash) == string(existingLDT.Hash) {
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

func (c *DiscoveryConfig) updateRepositories() error {
	var content io.Reader
	if isURL(c.repository_source) {
		resp, err := http.Get(c.repository_source)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		content = resp.Body
	} else {
		file, err := os.Open(c.repository_source)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		content = file
	}
	scanner := bufio.NewScanner(content)
	for scanner.Scan() {
		if !strings.HasPrefix(scanner.Text(), "#") {
			c.repositories = append(c.repositories, scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return nil
}

func isURL(input string) bool {
	u, err := url.Parse(input)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}
