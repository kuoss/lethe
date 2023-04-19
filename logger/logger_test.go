package logger

import (
	"bytes"
	"os"
	"runtime"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	assert.NotNil(t, logger)
	assert.NotZero(t, logger)
	assert.Equal(t, logger.Out, os.Stderr)

	assert.NotZero(t, logger.Formatter)
	assert.True(t, logger.Formatter.(*logrus.TextFormatter).FullTimestamp)
	assert.NotZero(t, logger.Formatter.(*logrus.TextFormatter).CallerPrettyfier)
	funcname, filename := logger.Formatter.(*logrus.TextFormatter).CallerPrettyfier(&runtime.Frame{
		PC:       0,
		Func:     &runtime.Func{},
		Function: "github.com/kuoss/lethe/logger.Func1", // not zero
		File:     "/a/b/c/file1.go",                     // not zero
		Line:     0,
		Entry:    0,
	})
	assert.Equal(t, "", funcname)
	assert.Equal(t, "???:1", filename)
}

func TestGetCallerPrettyfier(t *testing.T) {
	callerPrettyfier := getCallerPrettyfier()
	funcname, filename := callerPrettyfier(&runtime.Frame{
		PC:       0,
		Func:     &runtime.Func{},
		Function: "github.com/kuoss/lethe/logger.Func1", // not zero
		File:     "/a/b/c/file1.go",                     // not zero
		Line:     0,
		Entry:    0,
	})
	assert.Equal(t, "", funcname)
	assert.Equal(t, "???:1", filename)
}

func TestSetLevel(t *testing.T) {
	for _, level := range logrus.AllLevels {
		SetLevel(level)
		assert.Equal(t, level, logger.Level)
	}
	SetLevel(logrus.InfoLevel)
}

func captureOutput(f func()) string {
	buf := &bytes.Buffer{}
	logger.SetOutput(buf)
	f()
	logger.SetOutput(os.Stderr)
	return buf.String()
}

func TestDebugf(t *testing.T) {
	output := captureOutput(func() {
		Debugf("hello=%s lorem=%s number=%d", "hello", "ipsum", 42)
	})
	assert.Regexp(t, ``, output)
}

func TestInfof(t *testing.T) {
	output := captureOutput(func() {
		Infof("hello=%s lorem=%s number=%d", "hello", "ipsum", 42)
	})
	assert.Regexp(t, `time="[^"]+" level=info msg="hello=hello lorem=ipsum number=42" file="logger_test.go:[0-9]+"`, output)
}

func TestWarnf(t *testing.T) {
	output := captureOutput(func() {
		Warnf("hello=%s lorem=%s number=%d", "hello", "ipsum", 42)
	})
	assert.Regexp(t, `time="[^"]+" level=warning msg="hello=hello lorem=ipsum number=42" file="logger_test.go:[0-9]+"`, output)
}

func TestErrorf(t *testing.T) {
	output := captureOutput(func() {
		Errorf("hello=%s lorem=%s number=%d", "hello", "ipsum", 42)
	})
	assert.Regexp(t, `time="[^"]+" level=error msg="hello=hello lorem=ipsum number=42" file="logger_test.go:[0-9]+"`, output)
}
