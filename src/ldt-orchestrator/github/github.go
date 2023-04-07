package github

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	. "longevity/src/ldt-orchestrator"

	"github.com/google/go-github/v51/github"
	"golang.org/x/oauth2"
)

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type GithubDiscoverer struct {
	client       *github.Client
	autenticated bool
}

func NewGithubDiscoverer(token string) *GithubDiscoverer {
	ctx := context.Background()
	val, present := os.LookupEnv(token)

	if !present {
		log.Println("Github token not found. Requests will be limited to 60 per hour.")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: val},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &GithubDiscoverer{client: github.NewClient(tc), autenticated: present}
}

func (gd *GithubDiscoverer) GetReleasesFromRepository(owner, repo string) []*github.RepositoryRelease {
	releases, _, err := gd.client.Repositories.ListReleases(context.Background(), owner, repo, nil)
	if err != nil {
		log.Fatal(err)
	}
	return releases
}

func (gd *GithubDiscoverer) filterLDTsFromReleases(releases []*github.RepositoryRelease) bool {
	ldt_list := NewLDTList()
	for _, release := range releases {
		for _, asset := range release.Assets {
			url := asset.GetBrowserDownloadURL()
			ldt := filterLDTInformationFromURL(url)

			ldt_list.Ldts = append(ldt_list.Ldts, ldt)
		}
	}
	return false
}

func filterLDTInformationFromURL(address string) LDT {
	u, _ := url.Parse(address)
	user := strings.Split(u.Path, "/")[1]

	version := strings.Split(u.Path, "/")[5]

	filename := strings.Split(u.Path, "/")[6]
	withoutSuffix := strings.Split(filename, ".")[0]

	ldtname, rest, _ := strings.Cut(withoutSuffix, "_")
	os, arch, _ := strings.Cut(rest, "_")

	fmt.Println(withoutSuffix)

	ldt := LDT{
		Name:    user + "/" + ldtname,
		Version: version,
		Os:      os,
		Arch:    arch,
		Url:     address,
	}

	return ldt

}
