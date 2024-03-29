package wotmanager

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const path = "../../examples/wotm-description.json"

func Test_fetchWoTDescription(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	description, err := loadDescription(path)
	require.NoError(err)

	wotm := &WoTManager{description_raw: description}
	got, err := wotm.FetchWoTDescription()
	require.NoError(err)
	assert.NotNil(got)
}

func Test_getDeviceAddressFromDescription(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	var want string = "111.111.111.111"

	description_raw, err := loadDescription(path)
	require.NoError(err)

	wotm := &WoTManager{description_raw: description_raw}
	got := wotm.GetDeviceIPv4AddressFromDescription()
	assert.Equal(want, got)
}

func Test_getDeviceMACAddressFromDescription(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	var want string = "AA:BB:CC:DD:EE:FF"

	description_raw, err := loadDescription(path)
	require.NoError(err)

	wotm := &WoTManager{description_raw: description_raw}
	got := wotm.GetDeviceMACAddressFromDescription()
	assert.Equal(want, got)
}

func Test_getLDTPortFromDescription(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	var want int = 98765

	description_raw, err := loadDescription(path)
	require.NoError(err)

	wotm := &WoTManager{description_raw: description_raw}
	got := wotm.GetLdtPortFromDescription()
	assert.Equal(want, got)
}

func loadDescription(path string) (string, error) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var description string = string(buffer)
	return description, nil
}
