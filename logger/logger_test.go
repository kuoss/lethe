package logger

import (
	"runtime"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetLogger(t *testing.T) {
	logger := GetLogger()
	assert.NotZero(t, logger)

	assert.NotZero(t, logger.Formatter)
	assert.True(t, logger.Formatter.(*logrus.TextFormatter).FullTimestamp)
	assert.NotZero(t, logger.Formatter.(*logrus.TextFormatter).CallerPrettyfier)

	funcname, filename := logger.Formatter.(*logrus.TextFormatter).CallerPrettyfier(&runtime.Frame{
		File:     "/a/b/c/file1.go",
		Function: "func1",
	})
	assert.Equal(t, "func1", funcname)
	assert.Equal(t, "file1.go", filename)
}

func TestGetLogger_Singleton(t *testing.T) {
	logger1 := GetLogger()
	logger2 := GetLogger()
	assert.Equal(t, logger1, logger2)
}

func TestNewLogger_NotSingleton(t *testing.T) {
	logger1 := newLogger()
	logger2 := newLogger()
	assert.NotEqual(t, logger1, logger2)
}
