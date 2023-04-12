package github

import (
	"context"
	"log"
	"net/url"
	"os"
	"strings"

	. "longevity/src/types"

	"github.com/google/go-github/v51/github"
	"golang.org/x/oauth2"
)

type GithubClient struct {
	Client        *github.Client
	Authenticated bool
}

func NewGithubClient(token string) *GithubClient {
	ctx := context.Background()
	val, present := os.LookupEnv(token)

	if !present {
		log.Println("Github token not found. Requests will be limited to 60 per hour.")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: val},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &GithubClient{Client: github.NewClient(tc), Authenticated: present}
}

func FetchGithubReleases(repositories []string) *LDTList {
	gh := NewGithubClient("GH_ACCESS_TOKEN")
	ldt_list := NewLDTList()

	for _, repo := range repositories {
		owner, repo := parseRepository(repo)
		releases := gh.GetReleasesFromRepository(owner, repo)
		currentLDTs := gh.FilterLDTsFromReleases(releases)
		ldt_list.LDTs = append(ldt_list.LDTs, currentLDTs.LDTs...)
	}
	return ldt_list
}

func (gd *GithubClient) GetReleasesFromRepository(owner, repo string) []*github.RepositoryRelease {
	releases, _, err := gd.Client.Repositories.ListReleases(context.Background(), owner, repo, nil)
	if err != nil {
		log.Fatal(err)
	}
	return releases
}

func (gd *GithubClient) FilterLDTsFromReleases(releases []*github.RepositoryRelease) *LDTList {
	ldt_list := NewLDTList()
	for _, release := range releases {
		for _, asset := range release.Assets {
			url := asset.GetBrowserDownloadURL()
			if isArchive(url) {
				ldt := filterLDTInformationFromURL(url)
				ldt_list.LDTs = append(ldt_list.LDTs, ldt)
			}
		}
	}
	return ldt_list
}

func filterLDTInformationFromURL(address string) LDT {
	u, _ := url.Parse(address)
	user := strings.Split(u.Path, "/")[1]

	version := strings.Split(u.Path, "/")[5]

	filename := strings.Split(u.Path, "/")[6]
	withoutSuffix := strings.Split(filename, ".")[0]

	ldtname, rest, _ := strings.Cut(withoutSuffix, "_")
	os, arch, _ := strings.Cut(rest, "_")

	switch arch {
	case "x86_64":
		arch = "amd64"
	}

	ldt := LDT{
		Name:    user + "/" + ldtname,
		Version: version,
		Os:      strings.ToLower(os),
		Arch:    arch,
		Url:     address,
	}
	return ldt
}

func isArchive(file string) bool {
	return strings.HasSuffix(file, ".tar.gz")
}

func parseRepository(repo string) (string, string) {
	split := strings.Split(repo, "/")
	return split[3], split[4]
}
