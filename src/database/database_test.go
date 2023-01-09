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
	err = sample.createDevice(sql_db)
	assert.NoError(err)
}

func Test_createAlreadyExistingDevice(t *testing.T) {
	assert := assert.New(t)
	sql_db, err := sql.Open("sqlite3", test_db.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer sql_db.Close()
	_ = sample.createDevice(sql_db)
	err = sample.createDevice(sql_db)
	assert.Error(err)
}

// func Test_MatchingMacAddressRaisesError(t *testing.T) {
// 	assert := assert.New(t)
// 	err := checkMatchingMacAddress("11:22:33:44:55", sample)
// 	assert.Error(err)
// }

// func Test_matchingMacAddressSucceeds(t *testing.T) {
// 	assert := assert.New(t)
// 	err := checkMatchingMacAddress("00:11:22:33:44", sample)
// 	assert.Nil(err)
// }

// func Test_AddEntryToDatabase(t *testing.T) {
// 	t.Skip()
// }

// func Test_DeleteEntryFromDatabase(t *testing.T) {
// 	t.Skip()
// }

// func Test_CheckIfDeviceExists(t *testing.T) {
// 	t.Skip()
// }
// func Test_EnsureMacAddressKeyIsUnique(t *testing.T) {
// 	t.Skip()
// }
