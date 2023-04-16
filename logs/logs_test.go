package logs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var hello1 PatternedString = "hello"
var hello2 PatternedString = "hello*"
var hello3 PatternedString = "hello.*"

func Test_Patterned(t *testing.T) {
	assert.Equal(t, hello1.Patterned(), false)
	assert.Equal(t, hello2.Patterned(), true)
	assert.Equal(t, hello3.Patterned(), true)
}

func Test_PatternMatch(t *testing.T) {
	assert.Equal(t, hello1.PatternMatch("hello world"), false)
	assert.Equal(t, hello2.PatternMatch("hello world"), true)
	assert.Equal(t, hello3.PatternMatch("hello world"), true)
}

func Test_withoutPattern(t *testing.T) {
	assert.Equal(t, hello1.withoutPattern(), "hello")
	assert.Equal(t, hello2.withoutPattern(), "hello")
	assert.Equal(t, hello3.withoutPattern(), "hello")
}

func Test_New(t *testing.T) {
	var logStore *LogStore = New()
	assert.Equal(t, logStore, New())
}
