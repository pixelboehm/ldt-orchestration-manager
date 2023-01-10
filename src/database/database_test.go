package database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	model "longevity/src/model"

	_ "github.com/mattn/go-sqlite3"

	"github.com/stretchr/testify/assert"
)

var sample = &model.Device{
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
	err = CreateDevice(sql_db, sample)
	assert.NoError(err)
}

func Test_createAlreadyExistingDevice(t *testing.T) {
	assert := assert.New(t)
	sql_db, err := sql.Open("sqlite3", test_db.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer sql_db.Close()
	_ = CreateDevice(sql_db, sample)
	err = CreateDevice(sql_db, sample)
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
	err = updateDevice(sql_db, sample)
	assert.NoError(err)
}

func Test_DeleteDevice(t *testing.T) {
	assert := assert.New(t)
	sql_db, err := sql.Open("sqlite3", test_db.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer sql_db.Close()
	err = deleteDevice(sql_db, sample)
	assert.NoError(err)
}

func Test_GetDevice(t *testing.T) {
	assert := assert.New(t)
	sql_db, err := sql.Open("sqlite3", test_db.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer sql_db.Close()
	CreateDevice(sql_db, sample)
	var test_device = &model.Device{
		ID:         1,
		Name:       "",
		MacAddress: "",
		Twin:       "",
		Version:    "",
	}

	GetDevice(sql_db, test_device)
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
	_ = CreateDevice(sql_db, sample)
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

	var device1 = &model.Device{
		Name:       "Foo",
		MacAddress: "55:55:55:55:55",
		Twin:       "vs-lite",
		Version:    "0.0.1",
	}
	var device2 = &model.Device{
		Name:       "Bar",
		MacAddress: "55:55:55:55:55",
		Twin:       "vs-full",
		Version:    "0.0.1",
	}

	_ = CreateDevice(sql_db, device1)
	err = CreateDevice(sql_db, device2)
	assert.Error(err)
}

func Test_GetDevices(t *testing.T) {
	assert := assert.New(t)
	sql_db, err := sql.Open("sqlite3", test_db.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer sql_db.Close()

	var device1 = &model.Device{
		Name:       "Foo",
		MacAddress: "55:55:55:55:55",
		Twin:       "vs-lite",
		Version:    "0.0.1",
	}
	var device2 = &model.Device{
		Name:       "Bar",
		MacAddress: "44:44:44:44:44",
		Twin:       "vs-full",
		Version:    "0.0.1",
	}

	_ = CreateDevice(sql_db, device1)
	_ = CreateDevice(sql_db, device2)

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
