package database

import (
	"database/sql"
	"fmt"
	"log"
	. "longevity/src/model"
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

func AddDeviceToDatabase(d *Device) {
	db, err := sql.Open("sqlite3", "./longevity.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	statement, _ := db.Prepare("INSERT INTO devices (name, macAddress, twin, version) VALUES (?, ?, ?, ?)")
	statement.Exec(d.Name, d.MacAddress, d.Twin, d.Version)
	log.Printf("Inserted device %s!\n", d.Name)
}

func setup() {
	file, err := os.Create("./sqlite.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sqlite.db created")
}
