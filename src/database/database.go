package database

import (
	"database/sql"
	"errors"
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

func UpdateDevice(macAddress string, d *Device) {
	if !checkIfDeviceExists(macAddress) {
		log.Printf("Device with macAddress %s does not exist", macAddress)
		return
	}

	err := checkMatchingMacAdress(macAddress, d)
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("sqlite3", "./longevity.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	statement, _ := db.Prepare("update devices set name=?, twin=?, version=? where macAddress=?")
	statement.Exec(d.Name, d.Twin, d.Version, macAddress)
	log.Println("Successfully updated the device in database!")
}

func RemoveDevice(macAddress string) {
	if !checkIfDeviceExists(macAddress) {
		log.Printf("Device with macAddress %s does not exist", macAddress)
		return
	}
	db, err := sql.Open("sqlite3", "./longevity.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	statement, _ := db.Prepare("delete from devices where macAddress=?")
	statement.Exec(macAddress)
	log.Println("Successfully deleted the device in database!")
}

func ReadTable(name string) {
	db, err := sql.Open("sqlite3", "./longevity.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	rows, _ := db.Query("SELECT name, macAddress, twin, version FROM devices")
	var device Device
	for rows.Next() {
		rows.Scan(&device.Name, &device.MacAddress, &device.Twin, &device.Version)
		log.Printf("name: %s, macAddress: %s, twin: %s, version: %s\n",
			device.Name, device.MacAddress, device.Twin, device.Version)
	}

}

func setup() {
	file, err := os.Create("./longevity.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("longevity.db created")
}

func checkIfDeviceExists(address string) bool {
	db, err := sql.Open("sqlite3", "./longevity.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	rows, _ := db.Query("SELECT EXISTS(SELECT 1 FROM devices WHERE macAddress=?);", address)
	var result bool
	for rows.Next() {
		rows.Scan(&result)
	}
	return result
}

func checkMatchingMacAdress(address string, d *Device) error {
	if address != d.MacAddress {
		return errors.New("MacAdresses do not match")
	}
	return nil
}
