package config

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	err := LoadConfig()
	if err != nil {
		panic(err)
	}
}

func TestLoadConfig(t *testing.T) {
	assert.NotNil(t, vip)
	assert.NotZero(t, vip)
}

func TestViper(t *testing.T) {
	vip := Viper()
	assert.NotNil(t, vip)
	assert.NotZero(t, vip)
}

func TestGetWriter(t *testing.T) {
	assert.Equal(t, os.Stdout, GetWriter())
}

func TestSetWriter(t *testing.T) {
	tempWriter := new(bytes.Buffer)
	SetWriter(tempWriter)
	assert.Equal(t, tempWriter, GetWriter())
	SetWriter(os.Stdout)
}

func TestGetLimit(t *testing.T) {
	assert.Equal(t, 1000, GetLimit())
}

func TestSetLimit(t *testing.T) {
	SetLimit(2000)
	assert.Equal(t, 2000, GetLimit())
	SetLimit(1000)
}

func TestGetLogDataPath(t *testing.T) {
	assert.Equal(t, "/data/log", GetLogDataPath())
}

func TestSetLogDataPath(t *testing.T) {
	SetLogDataPath("hello")
	assert.Equal(t, "hello", GetLogDataPath())
	SetLogDataPath("/data/log")
}
