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

func Run(distro string) {
	updateRepositories()

	c := make(chan *pkgcloud.Package, 100)
	wg := sync.WaitGroup{}
	for _, repo := range repositories {
		wg.Add(1)
		go GetPackagesFromRepo(repo, distro, c, &wg)
	}
	wg.Wait()
	close(c)

	for pkg := range c {
		packageList = append(packageList, *pkg)
	}
	fmt.Printf("Found %d packages\n", len(packageList))
}

func GetPackagesFromRepo(repo, distro string, c chan *pkgcloud.Package, wg *sync.WaitGroup) {
	client, _ := setup()
	client.ShowProgress(true)

	packages, err := client.All(repo)
	if err != nil {
		log.Fatal(err)
	}
	for _, pkg := range packages {
		c <- &pkg
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
