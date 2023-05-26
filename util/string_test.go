package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubstrAfter(t *testing.T) {
	testCases := []struct {
		haystack string
		needle   string
		want     string
	}{
		{"hello world", "w", "orld"},
		{"hello world", " ", "world"},
		{"hello world", "l", "lo world"},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			got := SubstrAfter(tc.haystack, tc.needle)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestSubstrBefore(t *testing.T) {
	assert.Equal(t, "hello ", SubstrBefore("hello world", "w"))
	assert.Equal(t, "hello", SubstrBefore("hello world", " "))
	assert.Equal(t, "he", SubstrBefore("hello world", "l"))
}

func TestCountNewlines(t *testing.T) {
	assert.Equal(t, "1", CountNewlines("hello world"))
	assert.Equal(t, "2", CountNewlines("hello\nworld"))
	assert.Equal(t, "2", CountNewlines("hello\nworld\n"))

}

func TestStringToBytes(t *testing.T) {
	var out int
	out, _ = StringToBytes("100k")
	assert.Equal(t, 102400, out)
	out, _ = StringToBytes("100m")
	assert.Equal(t, 104857600, out)
	out, _ = StringToBytes("100g")
	assert.Equal(t, 107374182400, out)
}
