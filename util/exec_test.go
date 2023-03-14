package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Execute(t *testing.T) {
	var str string
	var code int
	var err error

	str, code, err = Execute("asdfajkshfajskf")
	assert.Equal(t, "", str)
	assert.Equal(t, 127, code)
	assert.NotNil(t, err)

	str, code, err = Execute("echo hello")
	assert.Equal(t, "hello\n", str)
	assert.Equal(t, 0, code)
	assert.Nil(t, err)
}
