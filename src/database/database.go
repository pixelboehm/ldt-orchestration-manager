package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func Start() {
	setup()
	sqliteDatabase, err := sql.Open("sqlite3", "./longevity.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	createTable()
	defer sqliteDatabase.Close()
}

func createTable() {
	db, err := sql.Open("sqlite3", "./longevity.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS devices (id INTEGER PRIMARY KEY, name VARCHAR(64), macAddress VARCHAR(17), twin VARCHAR(64), version string)")

	if err != nil {
		log.Fatal(err)
	} else {
		statement.Exec()
		fmt.Println("Created Table devices")
	}
}

func AddDeviceToDatabase(name string, macAddress string, twin string, version string) {
	db, err := sql.Open("sqlite3", "./longevity.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	statement, _ := db.Prepare("INSERT INTO devices (name, macAddress, twin, version) VALUES (?, ?, ?, ?)")
	statement.Exec(name, macAddress, twin, version)
	log.Printf("Inserted device %s!\n", name)
}

func setup() {
	file, err := os.Create("./sqlite.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sqlite.db created")
}
