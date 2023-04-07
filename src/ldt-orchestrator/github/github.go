package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

func GetReleasesFromRepository(owner, repo string) ([]Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var releases []Release
	err = json.NewDecoder(resp.Body).Decode(&releases)
	if err != nil {
		fmt.Println("Error decoding releases: ", err)
		return nil, err
	}
	return releases, nil
}

func filterLDTs(release Release) bool {
	return false
}

func filterURL(address string) (string, string, string, string) {
	u, _ := url.Parse(address)
	user := strings.Split(u.Path, "/")[1]

	version := strings.Split(u.Path, "/")[5]

	filename := strings.Split(u.Path, "/")[6]
	withoutSuffix := strings.Split(filename, ".")[0]

	ldtname, rest, _ := strings.Cut(withoutSuffix, "_")
	os, arch, _ := strings.Cut(rest, "_")

	fmt.Println(withoutSuffix)

	return user + "/" + ldtname, version, os, arch

}
