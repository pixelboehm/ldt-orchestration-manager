package ldtorchestrator

import (
	"bufio"
	"log"
	"os"

	"github.com/mlafeldt/pkgcloud"
)

var client *pkgcloud.Client
var err error
var packageList []pkgcloud.Package
var repositories []string

func Run() {
	client, err = setup()
	client.ShowProgress(true)
	updateRepositories()

	for _, repo := range repositories {
		GetPackagesFromRepo(repo)
	}
}

func GetPackagesFromRepo(repo string) {
	if client == nil {
		log.Fatal("Client is not initialized")
	}
	packages, err := client.All(repo)
	if err != nil {
		log.Fatal(err)
	}
	for _, pkg := range packages {
		packageList = append(packageList, pkg)
	}
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
