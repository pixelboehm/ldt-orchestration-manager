package database

import (
	. "longevity/src/model"
	"testing"
)

func Test_AddEntryToDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	var sample = NewDevice("Foo", "00:11:22:33:44", "general", "0.0.1")
	var sample2 = NewDevice("Bar", "11:22:33:44:55", "general", "0.0.1")
	Start()
	AddDeviceToDatabase(&sample)
	AddDeviceToDatabase(&sample2)
	ReadTable("devices")
}
