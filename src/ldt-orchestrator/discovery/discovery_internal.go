package discovery

import (
	. "longevity/src/types"
	"net/url"
	"runtime"
	"sort"
	"strings"
)

func getRuntimeInformation() (string, string) {
	os := runtime.GOOS
	arch := runtime.GOARCH
	return os, arch
}

func updateLatestTag(ldts *[]LDT) {
	sortByName(*ldts)

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
			*ldts = append([]LDT{*latest_ldt}, *ldts...)
		}
	}
}

func sortByName(ldts []LDT) {
	sort.Slice(ldts, func(i, j int) bool {
		if ldts[i].Vendor != ldts[j].Vendor {
			return ldts[i].Vendor > ldts[j].Vendor
		} else if ldts[i].Name != ldts[j].Name {
			return ldts[i].Name > ldts[j].Name
		} else {
			return ldts[i].Version > ldts[j].Version
		}
	})
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

func isURL(input string) bool {
	u, err := url.Parse(input)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}
