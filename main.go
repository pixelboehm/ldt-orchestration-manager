package main

import (
	rest "longevity/src/rest"
)

func main() {
	router := rest.Setup()
	rest.Start(router)
}
