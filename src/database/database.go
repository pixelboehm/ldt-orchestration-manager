package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
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
	getDevice()
}

func SetupSQLiteDB(db_name string) *sql.DB {
	file, err := os.Create(db_name)
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	db, err := sql.Open("sqlite3", db_name)
	if err != nil {
		log.Fatal(err)
	}
	CreateTable(db, tableCreationQuerySQLite)
	return db
}

func SetupPostgresDB(user string, password string, db_name string) *sql.DB {
	connection := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, db_name)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	CreateTable(db, tableCreationQuery)
	return db
}

func (d *Device) GetDevice(db *sql.DB) error {
	return db.QueryRow(getDeviceByIDQuery, d.ID).Scan(&d.Name, &d.MacAddress, &d.Twin, &d.Version)
}

func (d *Device) UpdateDevice(db *sql.DB) error {
	statement, _ := db.Prepare(updateDeviceQuery)
	_, err := statement.Exec(d.Name, d.MacAddress, d.Twin, d.Version, d.ID)

	return err
}

func (d *Device) DeleteDevice(db *sql.DB) error {
	statement, _ := db.Prepare(deleteDeviceQuery)
	_, err := statement.Exec(d.ID)

	return err
}

func (d *Device) CreateDevice(db *sql.DB) error {
	statement, _ := db.Prepare(insertDeviceQuery)
	_, err := statement.Exec(d.Name, d.MacAddress, d.Twin, d.Version)

	if err != nil {
		log.Println("Device with macAddress already exists")
		return err
	}
	log.Printf("New Device %s added!\n", d.Name)
	return nil
}

func getDevices(db *sql.DB) ([]Device, error) {
	rows, err := db.Query("SELECT name, macAddress, twin, version FROM devices")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	devices := []Device{}

	for rows.Next() {
		var d Device
		err := rows.Scan(&d.Name, &d.MacAddress, &d.Twin, &d.Version)
		if err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}

	return devices, nil
}

func CreateTable(db *sql.DB, query string) {
	statement, err := db.Prepare(query)

	if err != nil {
		log.Fatal(err)
	} else {
		statement.Exec()
		fmt.Println("Created Table devices")
	}
}

func checkIfDeviceExists(address string, db *sql.DB) bool {
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
