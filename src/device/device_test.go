package device

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getName(t *testing.T) {
	assert := assert.New(t)
	name := "Thermostat"
	device := Device{name}
	want := name
	got := device.getName()
	assert.Equal(want, got)
}
