package device

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createDevice(t *testing.T) {
	assert := assert.New(t)
	var tests = []struct {
		name  string
		input string
		want  device
	}{
		{"Device should be named Foo", "Foo", device{"Foo"}},
		{"Device should be named 1234", "1234", device{"12344"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := NewDevice(tt.input)
			assert.Equal(ans, tt.want)
		})
	}
}

func Test_getName(t *testing.T) {
	assert := assert.New(t)
	name := "Thermostat"
	device := NewDevice(name)
	want := name
	got := device.getName()
	assert.Equal(want, got)
}
