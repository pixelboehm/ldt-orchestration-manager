package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Stringer interface {
	String() string
}

func TestFail(t *testing.T) {
	assert := assert.New(t)
	assert.Fail("This test should fail")
}
