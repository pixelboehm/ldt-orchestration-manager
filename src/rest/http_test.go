package rest

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createNewDeviceViaPostRequest(t *testing.T) {
	assert := assert.New(t)
	rest := NewRestInterface()
	values, _ := url.ParseQuery("name=thermostat&macAddress=00:11:22:33:44")
	assert.HTTPSuccess(rest.SetNewDevice, "POST", "/devices", values, nil)
}

func Test_createNewDeviceWithEmptyNameViaPostRequest(t *testing.T) {
	assert := assert.New(t)
	rest := NewRestInterface()
	values, _ := url.ParseQuery("name=")
	assert.HTTPError(rest.SetNewDevice, "POST", "/devices", values, nil)
}

func Test_GetDevicesSucceeds(t *testing.T) {
	assert := assert.New(t)
	rest := RESTInterface{}
	assert.HTTPBodyContains(rest.GetDevices, "GET", "/devices", nil, "foo")
}
