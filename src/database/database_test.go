package database

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// var sample = &Device{
// 	ID:         1,
// 	Name:       "Foo",
// 	MacAddress: "00:11:22:33:44",
// 	Twin:       "general",
// 	Version:    "0.0.1"}

var sqlite_db string = "./test_db.db"
var db *sql.DB

func TestMain(m *testing.M) {
	// db = SetupSQLiteDB(sqlite_db)
	// defer os.Remove(sqlite_db)

	db = SetupPostgresDB("postgres", "foobar", "postgres")

	defer db.Close()

	code := m.Run()
	// clearTable()
	// CleanUpDatabase(db)
	os.Exit(code)
}

func clearTable() {
	db.Exec("ALTER SEQUENCE devices_id_seq RESTART WITH 1")
	db.Exec("DELETE FROM devices")
}

func Test_CreateDevice(t *testing.T) {
	clearTable()
	assert := assert.New(t)

	var sample *Device = NewDevice("Foo", "00:11:22:33:44", "general", "0.0.1")

	err := sample.CreateDevice(db)
	assert.NoError(err)
}

func Test_CreateAlreadyExistingDevice(t *testing.T) {
	clearTable()
	assert := assert.New(t)

	var sample *Device = NewDevice("Foo", "00:11:22:33:44", "general", "0.0.1")

	_ = sample.CreateDevice(db)
	err := sample.CreateDevice(db)
	assert.Error(err)
}

func Test_UpdateExistingDevice(t *testing.T) {
	clearTable()
	assert := assert.New(t)
	require := require.New(t)

	var sample *Device = NewDevice("Foo", "00:11:22:33:44", "general", "0.0.1")
	err := sample.CreateDevice(db)
	require.NoError(err)

	sample.Name = "Foo Updated"
	err = sample.UpdateDevice(db)
	assert.NoError(err)
}

func Test_DeleteDevice(t *testing.T) {
	clearTable()
	assert := assert.New(t)
	require := require.New(t)

	var sample *Device = NewDevice("Foo", "00:11:22:33:44", "general", "0.0.1")
	err := sample.CreateDevice(db)
	require.NoError(err)

	err = sample.DeleteDevice(db)
	assert.NoError(err)
}

func Test_GetDevice(t *testing.T) {
	clearTable()
	assert := assert.New(t)
	require := require.New(t)

	var sample *Device = NewDevice("Foo", "00:11:22:33:44", "general", "0.0.1")
	err := sample.CreateDevice(db)
	require.NoError(err)

	var result_device = NewDevice("", "", "", "")
	result_device.ID = 1

	_ = result_device.Device(db)
	require.NoError(err)

	assert.Equal("Foo", result_device.Name)
	assert.Equal("00:11:22:33:44", result_device.MacAddress)
	assert.Equal("general", result_device.Twin)
	assert.Equal("0.0.1", result_device.Version)
}

func Test_MatchingMacAddressRaisesError(t *testing.T) {
	assert := assert.New(t)

	var sample *Device = NewDevice("Foo", "00:11:22:33:44", "general", "0.0.1")
	err := checkMatchingMacAddress("11:22:33:44:55", sample)
	assert.Error(err)
}

func Test_MatchingMacAddressSucceeds(t *testing.T) {
	assert := assert.New(t)

	var sample *Device = NewDevice("Foo", "00:11:22:33:44", "general", "0.0.1")
	err := checkMatchingMacAddress("00:11:22:33:44", sample)
	assert.Nil(err)
}

func Test_CheckIfDeviceExists(t *testing.T) {
	clearTable()
	assert := assert.New(t)
	require := require.New(t)

	var sample *Device = NewDevice("Foo", "00:11:22:33:44", "general", "0.0.1")
	err := sample.CreateDevice(db)
	require.NoError(err)

	res := checkIfDeviceExists("00:11:22:33:44", db)
	assert.True(res)
}

func Test_EnsureMacAddressKeyIsUnique(t *testing.T) {
	clearTable()
	assert := assert.New(t)
	require := require.New(t)

	var device1 *Device = NewDevice("Foo", "00:11:22:33:44", "vs-lite", "0.0.1")
	var device2 *Device = NewDevice("Bar", "00:11:22:33:44", "vs-full", "0.0.1")

	err := device1.CreateDevice(db)
	require.NoError(err)
	err = device2.CreateDevice(db)
	assert.Error(err)
}

func Test_GetDevices(t *testing.T) {
	clearTable()
	assert := assert.New(t)
	require := require.New(t)

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

	err := device1.CreateDevice(db)
	require.NoError(err)
	err = device2.CreateDevice(db)
	require.NoError(err)

	devices, err := Devices(db)
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
