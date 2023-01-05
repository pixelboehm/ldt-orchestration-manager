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

type database interface {
	Start()
	CreateTable()
	AddDevice()
	UpdateDevice()
	RemoveDevice()
	PrintTable()
	setup()
}

type db struct{}

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

	statement, err := db.Prepare(tableCreationQuery)

	if err != nil {
		log.Fatal(err)
	} else {
		statement.Exec()
		fmt.Println("Created Table devices")
	}
}

func AddDevice(d *Device) error {
	db, err := sql.Open("sqlite3", "./longevity.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	statement, _ := db.Prepare(insertDeviceQuery)
	_, err = statement.Exec(d.Name, d.MacAddress, d.Twin, d.Version)
	if err != nil {
		log.Println("Device with macAddress already exists")
		return err
	}
	log.Printf("Inserted device %s!\n", d.Name)
	return nil
}

func UpdateDevice(macAddress string, d *Device) {
	if !checkIfDeviceExists(macAddress) {
		log.Printf("Device with macAddress %s does not exist", macAddress)
		return
	}

	err := checkMatchingMacAddress(macAddress, d)
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("sqlite3", "./longevity.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	statement, _ := db.Prepare(updateDeviceQuery)
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

	statement, _ := db.Prepare(deleteDeviceQuery)
	statement.Exec(macAddress)
	log.Println("Successfully deleted the device in database!")
}

func PrintTable(name string) {
	db, err := sql.Open("sqlite3", "./longevity.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	rows, _ := db.Query(getDeviceTableQuery)
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

	rows, _ := db.Query(checkIfDeviceExistsQuery, address)
	var result bool
	for rows.Next() {
		rows.Scan(&result)
	}
	return result
}

func checkMatchingMacAddress(address string, d *Device) error {
	if address != d.MacAddress {
		return errors.New("MacAdresses do not match")
	}
	return nil
}
