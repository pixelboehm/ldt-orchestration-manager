package ldtorchestrator

import (
	"bufio"
	"log"
	"os"
	"sync"

	"github.com/pixelboehm/pkgcloud"
)

type PackageList struct {
	packages []pkgcloud.Package
	lock     sync.Mutex
}

func GetLDTs(name, pkg_type, dist string) {
	packageList := &PackageList{
		packages: nil,
		lock:     sync.Mutex{},
	}
	repositories := updateRepositories()
	wg := sync.WaitGroup{}

	client, _ := setup()
	client.ShowProgress(false)

	for _, repo := range repositories {
		wg.Add(1)
		go FetchPackageProperties(client, repo, name, pkg_type, dist, packageList, &wg)
	}
	wg.Wait()
	log.Printf("Found %d packages\n", len(packageList.packages))
}

func FetchPackageProperties(client *pkgcloud.Client, repo, name, pkg_type, dist string, packageList *PackageList, wg *sync.WaitGroup) {
	packages, err := client.Search(repo, name, pkg_type, dist, 0)
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range packages {
		packageList.lock.Lock()
		packageList.packages = append(packageList.packages, pkg)
		packageList.lock.Unlock()
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

func clearCachedPackages(packageList *PackageList) {
	packageList.packages = nil
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
