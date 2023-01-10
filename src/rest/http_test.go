package rest

import (
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
	rest = NewRestInterface()
	db = &DB{Path: "./test_db.db"}
	db.CreateTable()
	os.Exit(m.Run())
}

func clearTable() {
	os.Remove("./test_db.db")
}

func Test_EmptyTable(t *testing.T) {
	assert := assert.New(t)
	clearTable()

	req, _ := http.NewRequest("GET", "/devices", nil)
	response := executeRequest(req)

	assert.Equalf(http.StatusOK, response.Code, "Expected Response Code %d, Got %d\n", http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	rest.Router.ServeHTTP(rr, req)
	return rr
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
