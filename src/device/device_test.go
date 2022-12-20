package device

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createDevice(t *testing.T) {
	assert := assert.New(t)
	device := NewDevice("Some Name")
	assert.NotNilf(device, "Unable to create Device")
}

func Test_getName(t *testing.T) {
	assert := assert.New(t)
	name := "Thermostat"
	device := NewDevice(name)
	want := name
	got := device.getName()
	assert.Equal(want, got)
}
