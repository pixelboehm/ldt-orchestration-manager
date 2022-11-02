package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MyObject struct {
	Value string
}

func TestSomething(t *testing.T) {
	var object = MyObject{"Something"}

	require := require.New(t)
	assert := assert.New(t)

	// assert equality
	require.Equal(123, 123, "they should be equal")

	// assert inequality
	require.NotEqual(123, 456, "they should not be equal")

	// assert for nil (good for errors)
	// require.Nil(object)

	// assert for not nil (good when you expect something)
	if assert.NotNil(object) {
		// now we know that object isn't nil, we are safe to make
		// further assertions without causing any errors
		require.Equal("Something", object.Value)
	}
}
