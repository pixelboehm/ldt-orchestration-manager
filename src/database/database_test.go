package database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/stretchr/testify/assert"
)

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
	test_db.CreateTable()
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
	err = sample.UpdateDevice(sql_db)
	assert.NoError(err)
}

func Test_DeleteDevice(t *testing.T) {
	assert := assert.New(t)
	sql_db, err := sql.Open("sqlite3", test_db.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer sql_db.Close()
	err = sample.DeleteDevice(sql_db)
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

func Test_EnsureMacAddressKeyIsUnique(t *testing.T) {
	assert := assert.New(t)
	sql_db, err := sql.Open("sqlite3", test_db.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer sql_db.Close()

	var device1 = &DB_Device{
		Name:       "Foo",
		MacAddress: "55:55:55:55:55",
		Twin:       "vs-lite",
		Version:    "0.0.1",
	}
	var device2 = &DB_Device{
		Name:       "Bar",
		MacAddress: "55:55:55:55:55",
		Twin:       "vs-full",
		Version:    "0.0.1",
	}

	_ = device1.CreateDevice(sql_db)
	err = device2.CreateDevice(sql_db)
	assert.Error(err)
}

func Test_GetDevices(t *testing.T) {
	assert := assert.New(t)
	sql_db, err := sql.Open("sqlite3", test_db.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer sql_db.Close()

	var device1 = &DB_Device{
		Name:       "Foo",
		MacAddress: "55:55:55:55:55",
		Twin:       "vs-lite",
		Version:    "0.0.1",
	}
	var device2 = &DB_Device{
		Name:       "Bar",
		MacAddress: "44:44:44:44:44",
		Twin:       "vs-full",
		Version:    "0.0.1",
	}

	_ = device1.CreateDevice(sql_db)
	_ = device2.CreateDevice(sql_db)

	devices, err := getDevices(sql_db)
	assert.NoError(err)

	var containsDevice1 bool = false
	var containsDevice2 bool = false
	for _, device := range devices {
		if device.Name == device1.Name {
			containsDevice1 = true
		} else if device.Name == device2.Name {
			containsDevice2 = true
		}
	}
	assert.True(containsDevice1)
	assert.True(containsDevice2)
}
