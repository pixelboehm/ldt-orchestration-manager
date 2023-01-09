package database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/stretchr/testify/assert"
)

/*
* test: ensure macAddress Key is unique
 */

var sample = &DB_Device{
	ID:         1,
	Name:       "Foo",
	MacAddress: "00:11:22:33:44",
	Twin:       "general",
	Version:    "0.0.1"}

var test_db = &DB{
	Path: "./test_db.db",
}

func TestMain(m *testing.M) {
	Initialize(test_db.Path)
	createTable(test_db.Path)
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func clearTable() {
	os.Remove(test_db.Path)
}

func Test_createDevice(t *testing.T) {
	assert := assert.New(t)
	sql_db, err := sql.Open("sqlite3", test_db.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer sql_db.Close()
	err = sample.CreateDevice(sql_db)
	assert.NoError(err)
}

func Test_createAlreadyExistingDevice(t *testing.T) {
	assert := assert.New(t)
	sql_db, err := sql.Open("sqlite3", test_db.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer sql_db.Close()
	_ = sample.CreateDevice(sql_db)
	err = sample.CreateDevice(sql_db)
	assert.Error(err)
}

func Test_UpdateExistingDevice(t *testing.T) {
	assert := assert.New(t)
	sql_db, err := sql.Open("sqlite3", test_db.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer sql_db.Close()
	sample.Name = "Foo Updated"
	err = sample.updateDevice(sql_db)
	assert.NoError(err)
}

func Test_DeleteDevice(t *testing.T) {
	assert := assert.New(t)
	sql_db, err := sql.Open("sqlite3", test_db.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer sql_db.Close()
	err = sample.deleteDevice(sql_db)
	assert.NoError(err)
}

func Test_GetDevice(t *testing.T) {
	assert := assert.New(t)
	sql_db, err := sql.Open("sqlite3", test_db.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer sql_db.Close()
	sample.CreateDevice(sql_db)
	var test_device = &DB_Device{
		ID:         1,
		Name:       "",
		MacAddress: "",
		Twin:       "",
		Version:    "",
	}

	test_device.GetDevice(sql_db)
	assert.Equal("Foo Updated", test_device.Name)
	assert.Equal("00:11:22:33:44", test_device.MacAddress)
	assert.Equal("general", test_device.Twin)
	assert.Equal("0.0.1", test_device.Version)
}

func Test_MatchingMacAddressRaisesError(t *testing.T) {
	assert := assert.New(t)
	err := checkMatchingMacAddress("11:22:33:44:55", sample)
	assert.Error(err)
}

func Test_matchingMacAddressSucceeds(t *testing.T) {
	assert := assert.New(t)
	err := checkMatchingMacAddress("00:11:22:33:44", sample)
	assert.Nil(err)
}

func Test_CheckIfDeviceExists(t *testing.T) {
	assert := assert.New(t)
	sql_db, err := sql.Open("sqlite3", test_db.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer sql_db.Close()
	_ = sample.CreateDevice(sql_db)
	res := checkIfDeviceExists("00:11:22:33:44", test_db.Path)
	assert.True(res)
}
