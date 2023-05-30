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

func loadDescription(path string) (string, error) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var description string = string(buffer)
	return description, nil
}
