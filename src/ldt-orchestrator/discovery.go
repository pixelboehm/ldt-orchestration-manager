package ldtorchestrator

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/google/go-github/v38/github"
	"github.com/pixelboehm/pkgcloud"
	"golang.org/x/oauth2"
)

type PackageList struct {
	packages []string
	lock     sync.Mutex
}

type DiscoveryConfig struct {
	repository_file string
	packageList     *PackageList
	repositories    []string
}

func NewPackageList() *PackageList {
	return &PackageList{
		packages: nil,
		lock:     sync.Mutex{},
	}
}

func NewConfig(path string) *DiscoveryConfig {
	return &DiscoveryConfig{
		repository_file: path,
		packageList:     NewPackageList(),
		repositories:    make([]string, 0),
	}
}

func (c *DiscoveryConfig) GetLDTs(name, pkg_type, dist string) {
	c.repositories = c.updateRepositories()
	wg := sync.WaitGroup{}

	client, _ := setupPackagecloudClient()
	client.ShowProgress(false)

	for _, repo := range c.repositories {
		wg.Add(1)
		go FetchPackageProperties(client, repo, name, pkg_type, dist, c.packageList, &wg)
	}
	wg.Wait()
	log.Printf("Found %d packages\n", len(c.packageList.packages))
}

func FetchPackageProperties(client *pkgcloud.Client, repo, name, pkg_type, dist string, packageList *PackageList, wg *sync.WaitGroup) {
	packages, err := client.Search(repo, name, pkg_type, dist, 0)
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range packages {
		packageList.lock.Lock()
		packageList.packages = append(packageList.packages, pkg.Name)
		packageList.lock.Unlock()
	}
	wg.Done()
}

func FetchLDTProperties() error {
	client, err := setupGithubClient(os.Getenv("ACCESS_TOKEN"))
	if err != nil {
		return err
	}
	releases, _, err := client.Repositories.ListReleases(context.Background(), "pixelboehm", "ldt", nil)
	if err != nil {
		fmt.Println("Error getting releases:", err)
		return err
	}

	for _, release := range releases {
		fmt.Println(*release.Name)
		fmt.Printf("%s\n", *release.Assets[1].BrowserDownloadURL)
	}
	return nil
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

func clearCachedPackages(packageList *PackageList) {
	packageList.packages = nil
}

func clearCachedRepositories(repositories *[]string) {
	*repositories = nil
}

func setupGithubClient(github_token string) (*github.Client, error) {
	token := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: github_token},
	)
	oauthClient := oauth2.NewClient(context.Background(), token)

	client := github.NewClient(oauthClient)
	return client, nil
}

func setupPackagecloudClient() (*pkgcloud.Client, error) {
	client, err := pkgcloud.NewClient("")
	if err != nil {
		log.Fatal(err)
	}
	return client, err
}
