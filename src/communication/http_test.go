package communication

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	. "longevity/src/database"

	"github.com/stretchr/testify/assert"
)

var rest API
var sqlite_db string = "./test_db.db"
var db *sql.DB

func TestMain(m *testing.M) {
	// db = SetupSQLiteDB(sqlite_db)
	// defer os.Remove(sqlite_db)

	db = SetupPostgresDB("postgres", "foobar", "postgres")

	defer db.Close()

	rest = NewRestInterface(db)
	rest.initialize()

	code := m.Run()
	clearTable()
	os.Exit(code)
}

func clearTable() {
	rest.Database().Exec("ALTER SEQUENCE devices_id_seq RESTART WITH 1")
	rest.Database().Exec("DELETE FROM devices")
}

func Test_EmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/devices", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func Test_GetNonExistentDevice(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/device/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "device not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'device not found'. Got '%s'", m["error"])
	}
}

func Test_GetDevice(t *testing.T) {
	clearTable()
	err := addTestDevices(1)

	if err != nil {
		log.Fatal(err)
	}

	req, _ := http.NewRequest("GET", "/device/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func Test_CreateDevice(t *testing.T) {
	clearTable()
	assert := assert.New(t)

	var jsonStr = []byte(`{
		"name":"Device101",
		"macAddress": "11:22:33:44:55",
		"twin":"vs-full",
		"version":"0.0.1"
	}`)
	req, _ := http.NewRequest("POST", "/device", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var res map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &res)

	assert.Equal("Device101", res["name"])
	assert.Equal("11:22:33:44:55", res["macAddress"])
	assert.Equal("vs-full", res["twin"])
	assert.Equal("0.0.1", res["version"])
}

func Test_UpdateDevice(t *testing.T) {
	clearTable()
	assert := assert.New(t)

	err := addTestDevices(1)
	if err != nil {
		log.Fatal(err)
	}

	req, _ := http.NewRequest("GET", "/device/1", nil)
	response := executeRequest(req)
	var originalDevice map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalDevice)

	var jsonStr = []byte(`{
		"name":"new name for device", 
		"macAddress": "11:22:33:44:55",
		"twin": "vs-lite",
		"version": "0.0.1"
		}`)
	req, _ = http.NewRequest("PUT", "/device/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var res map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &res)

	assert.Equal("new name for device", res["name"])
	assert.Equal("11:22:33:44:55", res["macAddress"])
	assert.Equal("vs-lite", res["twin"])
	assert.Equal("0.0.1", res["version"])
}

func Test_DeleteDevice(t *testing.T) {
	clearTable()

	err := addTestDevices(1)
	if err != nil {
		log.Fatal(err)
	}

	req, _ := http.NewRequest("GET", "/device/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/device/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/device/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	rest.Router().ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, want int, got int) {
	assert := assert.New(t)
	assert.Equalf(want, got, "Expected Response Code %d, Got %d\n", want, got)
}

func addTestDevices(amount int) error {
	for i := 0; i < amount; i++ {
		var device1 = &Device{
			Name:       "Device" + strconv.Itoa(i),
			MacAddress: createDummyMacAddress(i),
			Twin:       "vs-lite",
			Version:    "0.0.1",
		}
		err := device1.CreateDevice(rest.Database())
		if err != nil {
			return err
		}
	}
	return nil
}

func createDummyMacAddress(i int) string {
	result := fmt.Sprintf("%d%d:%d%d:%d%d:%d%d", i, i, i, i, i, i, i, i)
	return result
}
