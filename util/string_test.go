package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SubstrAfter(t *testing.T) {
	assert.Equal(t, "orld", SubstrAfter("hello world", "w"))
	assert.Equal(t, "world", SubstrAfter("hello world", " "))
	assert.Equal(t, "lo world", SubstrAfter("hello world", "l"))
}

func Test_SubstrBefore(t *testing.T) {
	assert.Equal(t, "hello ", SubstrBefore("hello world", "w"))
	assert.Equal(t, "hello", SubstrBefore("hello world", " "))
	assert.Equal(t, "he", SubstrBefore("hello world", "l"))
}

func Test_CountNewlines(t *testing.T) {
	assert.Equal(t, "1", CountNewlines("hello world"))
	assert.Equal(t, "2", CountNewlines("hello\nworld"))
	assert.Equal(t, "2", CountNewlines("hello\nworld\n"))

}
