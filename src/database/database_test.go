package database

import (
	. "longevity/src/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

var sample = &Device{
	Name:       "Foo",
	MacAddress: "00:11:22:33:44",
	Twin:       "general",
	Version:    "0.0.1"}

func Test_MatchingMacAddressRaisesError(t *testing.T) {
	assert := assert.New(t)
	err := checkMatchingMacAdress("11:22:33:44:55", sample)
	assert.Error(err)
}

func Test_matchingMacAddressSucceeds(t *testing.T) {
	assert := assert.New(t)
	err := checkMatchingMacAdress("00:11:22:33:44", sample)
	assert.Nil(err)
}

func Test_AddEntryToDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	var sample = NewDevice("Foo", "00:11:22:33:44", "general", "0.0.1")
	var sample2 = NewDevice("Bar", "11:22:33:44:55", "general", "0.0.1")
	Start()
	AddDeviceToDatabase(&sample)
	AddDeviceToDatabase(&sample2)
	PrintTable("devices")
	sample2.Name = "Bar2"
	UpdateDevice("11:22:33:44:55", &sample2)
	PrintTable("devices")
}

func Test_DeleteEntryFromDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	var test_sample = NewDevice("Foo", "00:11:22:33:44", "general", "0.0.1")
	var test_sample2 = NewDevice("Bar", "11:22:33:44:55", "general", "0.0.1")
	Start()
	AddDeviceToDatabase(&test_sample)
	AddDeviceToDatabase(&test_sample2)
	PrintTable("devices")
	RemoveDevice("11:22:33:44:55")
	PrintTable("devices")
}

func Test_CheckIfDeviceExists(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	assert := assert.New(t)
	var sample = NewDevice("Foo", "00:11:22:33:44", "general", "0.0.1")
	Start()
	AddDeviceToDatabase(&sample)
	var tests = []struct {
		name       string
		macAddress string
		want       bool
	}{
		{"Device should exist", "00:11:22:33:44", true},
		{"Device should not exist", "11:22:33:44:55", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := checkIfDeviceExists(tt.macAddress)
			assert.Equal(tt.want, ans)
		})
	}
}
func Test_EnsureMacAddressKeyIsUnique(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	assert := assert.New(t)
	var sample = NewDevice("Foo", "00:11:22:33:44", "general", "0.0.1")
	var sample2 = NewDevice("Bar", "00:11:22:33:44", "general", "0.0.1")
	Start()
	AddDeviceToDatabase(&sample)
	err := AddDeviceToDatabase(&sample2)
	assert.Error(err)
	PrintTable("devices")
}
