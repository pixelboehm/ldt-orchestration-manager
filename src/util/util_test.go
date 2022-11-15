package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	assert := assert.New(t)
	first := 1
	second := 2
	expected := 3
	assert.Equal(expected, Add(first, second))
}
