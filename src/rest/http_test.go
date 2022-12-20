package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPass(t *testing.T) {
	assert := assert.New(t)
	assert.Truef(true, "this test should pass")
}
