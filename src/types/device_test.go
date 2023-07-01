package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var device_sample = &Device{
	Name:       "Foo",
	MacAddress: "00:11:22:33:44",
	Twin:       "general",
	Version:    "0.0.1"}

func Test_createDeviceConstructor(t *testing.T) {
	assert := assert.New(t)
	var tests = []struct {
		test_name  string
		name       string
		macAddress string
		twin       string
		version    string
		want       Device
	}{
		{"Device should be named Foo", "Foo", "00:11:22:33:44", "", "", Device{0, "Foo", "00:11:22:33:44", "", ""}},
		{"Device should have macAdress 00:11:22:33:44", "Foo", "00:11:22:33:44", "", "", Device{0, "Foo", "00:11:22:33:44", "", ""}},
		{"Device should have default twin called none", "Foo", "00:11:22:33:44", "", "", Device{0, "Foo", "00:11:22:33:44", "", ""}},
		{"Device should have default version called none", "Foo", "00:11:22:33:44", "", "", Device{0, "Foo", "00:11:22:33:44", "", ""}},
	}
	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			ans := NewDevice(tt.name, tt.macAddress, tt.twin, tt.version)
			assert.Equal(tt.want, *ans)
		})
	}
}
