package config

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadConfig(t *testing.T) {
	err := LoadConfig()
	assert.NotNil(t, config)
	assert.NotZero(t, config)
	assert.NoError(t, err)
}

func Test_GetConfig(t *testing.T) {
	ret := GetConfig()
	assert.NotNil(t, ret)
}

func Test_SetWriter(t *testing.T) {
	tempWriter := new(bytes.Buffer)
	SetWriter(tempWriter)
	assert.NotNil(t, writer, tempWriter)
}

func Test_GetWriter(t *testing.T) {
	ret := GetWriter()
	assert.NotNil(t, ret)
}

func Test_GetLimit(t *testing.T) {
	limit := GetLimit()
	assert.NotNil(t, limit)
}

func Test_SetLimit(t *testing.T) {
	assert.NotNil(t, limit)
}

func Test_GetLogRoot(t *testing.T) {
	logRoot := GetLogRoot()
	assert.NotNil(t, logRoot)
}

func Test_SetLogRoot(t *testing.T) {
	assert.NotNil(t, logRoot)
}
