package rest

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	. "longevity/src/database"

	"github.com/stretchr/testify/assert"
)

var rest *RESTInterface
var db *DB

func TestMain(m *testing.M) {
	db = &DB{Path: "./test_db.db"}
	db.CreateTable()
	sql_db, err := sql.Open("sqlite3", db.Path)
	if err != nil {
		log.Fatal(err)
	}

	rest = NewRestInterface(sql_db)
	rest.setup()
	db.CreateTable()
	os.Exit(m.Run())
}

func clearTable() {
	rest.DB.Exec("DELETE FROM devices")
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

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	rest.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, want int, got int) {
	assert := assert.New(t)
	assert.Equalf(want, got, "Expected Response Code %d, Got %d\n", want, got)
}

// func Test_createNewDeviceViaPostRequest(t *testing.T) {
// 	assert := assert.New(t)
// 	rest := NewRestInterface()
// 	values, _ := url.ParseQuery("name=thermostat&macAddress=00:11:22:33:44")
// 	assert.HTTPSuccess(rest.SetNewDevice, "POST", "/devices", values, nil)
// }

// func Test_createNewDeviceWithEmptyNameViaPostRequest(t *testing.T) {
// 	assert := assert.New(t)
// 	rest := NewRestInterface()
// 	values, _ := url.ParseQuery("name=")
// 	assert.HTTPError(rest.SetNewDevice, "POST", "/devices", values, nil)
// }

// func Test_GetDevicesSucceeds(t *testing.T) {
// 	assert := assert.New(t)
// 	rest := RESTInterface{}
// 	assert.HTTPBodyContains(rest.GetDevices, "GET", "/devices", nil, "foo")
// }
