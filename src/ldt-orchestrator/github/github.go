package github

import (
	"context"
	"log"
	"os"

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
		log.Println("<Disovery>: Github token not found. Requests will be limited to 60 per hour.")
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
		panic(err)
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
				finalizeLDT(ldt)
				ldt_list.LDTs = append(ldt_list.LDTs, *ldt)
			}
		}
	}
	return ldt_list
}
