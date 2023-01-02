package device

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

func Test_getName(t *testing.T) {
	assert := assert.New(t)
	device := NewDevice(sample.name, "", "", "")
	want := sample.name
	got := device.getName()
	assert.Equal(want, got)
}

func Test_getMacAddress(t *testing.T) {
	assert := assert.New(t)
	device := NewDevice("", sample.macAddress, "", "")
	want := sample.macAddress
	got := device.getMacAddress()
	assert.Equal(want, got)
}

func Test_getTwin(t *testing.T) {
	assert := assert.New(t)
	device := NewDevice("", "", sample.twin, "")
	want := sample.twin
	got := device.getTwin()
	assert.Equal(want, got)
}

func Test_getVersion(t *testing.T) {
	assert := assert.New(t)
	device := NewDevice("", "", "", sample.version)
	want := sample.version
	got := device.getVersion()
	assert.Equal(want, got)
}
