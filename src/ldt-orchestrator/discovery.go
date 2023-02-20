package ldtorchestrator

import (
	"fmt"
	"log"

	"github.com/mlafeldt/pkgcloud"
)

var client *pkgcloud.Client
var err error
var packageList []pkgcloud.Package

func Run() {
	client, err = setup()
	client.ShowProgress(true)
	GetPackagesFromRepo("pixelboehm/longevity-digital-twins")

	for _, pkg := range packageList {
		fmt.Println(pkg.Name)
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

func clearCachedPackages() {
	packageList = nil
}

func setup() (*pkgcloud.Client, error) {
	client, err := pkgcloud.NewClient("")
	if err != nil {
		log.Fatal(err)
	}
	return client, err
}
