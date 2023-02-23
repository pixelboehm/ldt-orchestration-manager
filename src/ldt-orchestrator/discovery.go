package ldtorchestrator

import (
	"bufio"
	"log"
	"os"
	"sync"

	"github.com/mlafeldt/pkgcloud"
)

func GetPackages(name, pkg_type, dist string) {
	var packageList []pkgcloud.Package
	repositories := updateRepositories()

	wg := sync.WaitGroup{}
	for _, repo := range repositories {
		wg.Add(1)
		go FetchPackageProperties(repo, name, pkg_type, dist, &packageList, &wg)
	}
	wg.Wait()
	log.Printf("Found %d packages\n", len(packageList))
}

func FetchPackageProperties(repo, name, pkg_type, dist string, packageList *[]pkgcloud.Package, wg *sync.WaitGroup) {
	client, _ := setup()
	client.ShowProgress(true)

	packages, err := client.Search(repo, name, pkg_type, dist, 0)
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range packages {
		*packageList = append(*packageList, pkg)
	}
	wg.Done()
}

func updateRepositories() []string {
	var repositories []string
	file, err := os.Open("src/ldt-orchestrator/repositories.list")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		repositories = append(repositories, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return repositories
}

func clearCachedPackages(packageList *[]pkgcloud.Package) {
	*packageList = nil
}

func clearCachedRepositories(repositories *[]string) {
	*repositories = nil
}

func setup() (*pkgcloud.Client, error) {
	client, err := pkgcloud.NewClient("")
	if err != nil {
		log.Fatal(err)
	}
	return client, err
}
