package rest

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createNewDeviceViaPostRequest(t *testing.T) {
	assert := assert.New(t)
	rest := RESTInterface{}
	values, _ := url.ParseQuery("name=thermostat")
	assert.HTTPSuccess(rest.SetNewDevice, "POST", "localhost:8000/devices", values, nil)
}

func Test_createNewDeviceWithEmptyNameViaPostRequest(t *testing.T) {
	assert := assert.New(t)
	rest := RESTInterface{}
	values, _ := url.ParseQuery("name=")
	assert.HTTPError(rest.SetNewDevice, "POST", "localhost:8000/devices", values, nil)
}

func Test_GetDevicesSucceeds(t *testing.T) {
	assert := assert.New(t)
	rest := RESTInterface{}
	assert.HTTPSuccess(rest.GetDevices, "GET", "localhost:8000/devices", nil, nil)
}
