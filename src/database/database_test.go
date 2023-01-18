package database

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/stretchr/testify/assert"
)

var sample = &Device{
	ID:         1,
	Name:       "Foo",
	MacAddress: "00:11:22:33:44",
	Twin:       "general",
	Version:    "0.0.1"}

var sqlite_db string = "./test_db.db"
var db *sql.DB

func TestMain(m *testing.M) {
	// db = SetupSQLiteDB(sqlite_db)
	// defer os.Remove(sqlite_db)

	db = SetupPostgresDB("postgres", "foobar", "postgres")

	defer db.Close()

	code := m.Run()
	clearTable()
	os.Exit(code)
}

func clearTable() {
	db.Exec("ALTER SEQUENCE devices_id_seq RESTART WITH 1")
	db.Exec("DELETE FROM devices")
}

func Test_createDevice(t *testing.T) {
	clearTable()
	assert := assert.New(t)

	err := sample.CreateDevice(db)
	assert.NoError(err)
	assert.True(true)
}

func Test_createAlreadyExistingDevice(t *testing.T) {
	clearTable()
	assert := assert.New(t)

	_ = sample.CreateDevice(db)
	err := sample.CreateDevice(db)
	assert.Error(err)
}

func Test_UpdateExistingDevice(t *testing.T) {
	clearTable()
	assert := assert.New(t)

	_ = sample.CreateDevice(db)

	sample.Name = "Foo Updated"
	err := sample.UpdateDevice(db)
	assert.NoError(err)
}

func Test_DeleteDevice(t *testing.T) {
	clearTable()
	assert := assert.New(t)

	_ = sample.CreateDevice(db)
	err := sample.DeleteDevice(db)
	assert.NoError(err)
}

func Test_GetDevice(t *testing.T) {
	clearTable()
	assert := assert.New(t)

	sample.CreateDevice(db)
	var test_device = &Device{
		ID:         1,
		Name:       "",
		MacAddress: "",
		Twin:       "",
		Version:    "",
	}

	test_device.GetDevice(db)
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
	clearTable()
	assert := assert.New(t)
	_ = sample.CreateDevice(db)
	res := checkIfDeviceExists("00:11:22:33:44", db)
	assert.True(res)
}

func Test_EnsureMacAddressKeyIsUnique(t *testing.T) {
	clearTable()
	assert := assert.New(t)

	var device1 = &Device{
		Name:       "Foo",
		MacAddress: "55:55:55:55:55",
		Twin:       "vs-lite",
		Version:    "0.0.1",
	}
	var device2 = &Device{
		Name:       "Bar",
		MacAddress: "55:55:55:55:55",
		Twin:       "vs-full",
		Version:    "0.0.1",
	}

	_ = device1.CreateDevice(db)
	err := device2.CreateDevice(db)
	assert.Error(err)
}

func Test_GetDevices(t *testing.T) {
	clearTable()
	assert := assert.New(t)

	var device1 = &Device{
		Name:       "Foo",
		MacAddress: "55:55:55:55:55",
		Twin:       "vs-lite",
		Version:    "0.0.1",
	}
	var device2 = &Device{
		Name:       "Bar",
		MacAddress: "44:44:44:44:44",
		Twin:       "vs-full",
		Version:    "0.0.1",
	}

	_ = device1.CreateDevice(db)
	_ = device2.CreateDevice(db)

	devices, err := getDevices(db)
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
