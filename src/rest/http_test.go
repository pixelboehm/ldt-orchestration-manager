package rest

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createNewDeviceViaPostRequest(t *testing.T) {
	rest := RESTInterface{}
	assert := assert.New(t)
	values, _ := url.ParseQuery("name=thermostat")
	assert.HTTPSuccess(rest.SetNewDevice, "POST", "localhost:8000/devices", values, nil)
}
