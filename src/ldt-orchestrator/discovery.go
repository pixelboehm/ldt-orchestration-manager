package ldtorchestrator

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/mlafeldt/pkgcloud"
)

var packageList []pkgcloud.Package
var repositories []string

func Run(name, pkg_type, dist string) {
	updateRepositories()

	wg := sync.WaitGroup{}
	for _, repo := range repositories {
		wg.Add(1)
		go GetPackagesFromRepo(repo, name, pkg_type, dist, &wg)
	}
	wg.Wait()
	log.Printf("Found %d packages\n", len(packageList))
	for _, pkg := range packageList {
		fmt.Println(pkg.PackageHtmlUrl)
	}
}

func GetPackagesFromRepo(repo, name, pkg_type, dist string, wg *sync.WaitGroup) {
	client, _ := setup()
	client.ShowProgress(true)

	packages, err := client.Search(repo, name, pkg_type, dist, 0)
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range packages {
		packageList = append(packageList, pkg)
	}
	wg.Done()
}

func updateRepositories() {
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
}

func clearCachedPackages() {
	packageList = nil
}

func clearCachedRepositories() {
	repositories = nil
}

func setup() (*pkgcloud.Client, error) {
	client, err := pkgcloud.NewClient("")
	if err != nil {
		log.Fatal(err)
	}
	return client, err
}
