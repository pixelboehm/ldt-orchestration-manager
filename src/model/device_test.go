package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var sample = &device{"Foo", "00:11:22:33:44", "general", "0.0.1"}

func Test_createDevice(t *testing.T) {
	assert := assert.New(t)
	var tests = []struct {
		test_name  string
		name       string
		macAddress string
		twin       string
		version    string
		want       device
	}{
		{"Device should be named Foo", "Foo", "00:11:22:33:44", "", "", device{"Foo", "00:11:22:33:44", "", ""}},
		{"Device should have macAdress 00:11:22:33:44", "Foo", "00:11:22:33:44", "", "", device{"Foo", "00:11:22:33:44", "", ""}},
		{"Device should have default twin called none", "Foo", "00:11:22:33:44", "", "", device{"Foo", "00:11:22:33:44", "", ""}},
		{"Device should have default version called none", "Foo", "00:11:22:33:44", "", "", device{"Foo", "00:11:22:33:44", "", ""}},
	}
	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			ans := NewDevice(tt.name, tt.macAddress, tt.twin, tt.version)
			assert.Equal(ans, tt.want)
		})
	}
}
func Test_getterFunctions(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		name     string
		function func() string
		want     string
	}{
		{"Device name should be", sample.getName, sample.name},
		{"Device macAddress should be", sample.getMacAddress, sample.macAddress},
		{"Device twin should be", sample.getTwin, sample.twin},
		{"Device version should be", sample.getVersion, sample.version},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := tt.function()
			assert.Equal(ans, tt.want)
		})
	}
}
