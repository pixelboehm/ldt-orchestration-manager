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
	Run()
	CreateTable()
	AddDevice()
	UpdateDevice()
	RemoveDevice()
	PrintTable()
	Initialize()
}

type DB struct {
	Path string
}

type DB_Device struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	MacAddress string `json:"macAddress"`
	Twin       string `json:"twin"`
	Version    string `json:"version"`
}

var app_db = &DB{
	Path: "./longevity.db",
}

func Run() {
	Initialize(app_db.Path)
	sqliteDatabase, err := sql.Open("sqlite3", app_db.Path)
	if err != nil {
		log.Fatal(err.Error())
	}
	createTable(app_db.Path)
	defer sqliteDatabase.Close()
}

func (d *DB_Device) getDevice(db *sql.DB) error {
	return db.QueryRow(getDeviceByIDQuery, d.ID).Scan(&d.Name, d.MacAddress, d.Twin, d.Version)
}

func (d *DB_Device) updateDevice(db *sql.DB) error {
	statement, _ := db.Prepare(updateDeviceQuery)
	_, err := statement.Exec(d.Name, d.MacAddress, d.Twin, d.Version, d.ID)

	return err
}

func (d *DB_Device) deleteDevice(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (d *DB_Device) createDevice(db *sql.DB) error {
	statement, _ := db.Prepare(insertDeviceQuery)
	_, err := statement.Exec(d.Name, d.MacAddress, d.Twin, d.Version)

	if err != nil {
		log.Println("Device with macAddress already exists")
		return err
	}
	log.Printf("New Device %s added!\n", d.Name)
	return nil
}

func getDevices(db *sql.DB, start int, count int) error {
	return errors.New("Not implemented")
}

func createTable(db_name string) {
	db, err := sql.Open("sqlite3", db_name)

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

func Initialize(db_name string) {
	file, err := os.Create(db_name)
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Printf("%s created\n", db_name)
}

func checkIfDeviceExists(address string) bool {
	db, err := sql.Open("sqlite3", app_db.Path)

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
