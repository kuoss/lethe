package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_StringToBytes(t *testing.T) {
	var out int
	out, _ = StringToBytes("100k")
	assert.Equal(t, 102400, out)
	out, _ = StringToBytes("100m")
	assert.Equal(t, 104857600, out)
	out, _ = StringToBytes("100g")
	assert.Equal(t, 107374182400, out)
}
