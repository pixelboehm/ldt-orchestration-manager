package discovery

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"longevity/src/ldt-orchestrator/github"
	. "longevity/src/types"
	"net/http"
	"os"
	"strings"
)

type Discoverer struct {
	repository_source string
	SupportedLDTs     *LDTList
	otherLDTs         *LDTList
	repositories      []string
	os                string
	arch              string
}

func NewDiscoverer(path string) *Discoverer {
	os, arch := getRuntimeInformation()
	return &Discoverer{
		repository_source: path,
		SupportedLDTs:     NewLDTList(),
		otherLDTs:         NewLDTList(),
		repositories:      make([]string, 0),
		os:                os,
		arch:              arch,
	}
}

func (discoverer *Discoverer) DiscoverLDTs() {
	_ = discoverer.updateRepositories()
	newLDTs := github.FetchGithubReleases(discoverer.repositories)
	var update_latest_tag_supported bool = false
	var update_latest_tag_other bool = false
	for _, ldt := range newLDTs.LDTs {
		if ldt.Os == discoverer.os && ldt.Arch == discoverer.arch {
			if !ldtAlreadyExists(&ldt, discoverer.SupportedLDTs) {
				discoverer.SupportedLDTs.LDTs = append(discoverer.SupportedLDTs.LDTs, ldt)
				update_latest_tag_supported = true
			}
		} else {
			if !ldtAlreadyExists(&ldt, discoverer.otherLDTs) {
				discoverer.otherLDTs.LDTs = append(discoverer.otherLDTs.LDTs, ldt)
				update_latest_tag_supported = true
			}
		}
	}
	if update_latest_tag_supported {
		updateLatestTag(&discoverer.SupportedLDTs.LDTs)
	}
	if update_latest_tag_other {
		updateLatestTag(&discoverer.otherLDTs.LDTs)
	}
}

func (discoverer *Discoverer) GetUrlFromLDTByID(id int) (string, error) {
	if id >= len(discoverer.SupportedLDTs.LDTs) {
		return "", errors.New("Failed to map ID to LDT")
	}
	return discoverer.SupportedLDTs.LDTs[id].Url, nil
}

func (discoverer *Discoverer) GetURLFromLDTByName(user, ldt, tag string) (string, error) {
	for _, entry := range discoverer.SupportedLDTs.LDTs {
		if entry.Vendor == user && entry.Name == ldt && entry.Version == tag {
			return entry.Url, nil
		}
	}
	return "", errors.New(fmt.Sprintf("Unable to find LDT %s", ldt))
}

func (c *Discoverer) updateRepositories() error {
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
