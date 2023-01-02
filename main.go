package main

import (
	"fmt"
	"longevity/src/database"
)

func main() {
	fmt.Print("main\n")
	database.Start()
}
