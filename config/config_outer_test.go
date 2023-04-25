package config_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/kuoss/lethe/config"
	"github.com/stretchr/testify/assert"
)

func init() {
	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
}

func TestViper(t *testing.T) {
	vip := config.Viper()
	assert.NotNil(t, vip)
	assert.NotZero(t, vip)
}

func TestGetWriter(t *testing.T) {
	assert.Equal(t, os.Stdout, config.GetWriter())
}

func TestSetWriter(t *testing.T) {
	tempWriter := new(bytes.Buffer)
	config.SetWriter(tempWriter)
	assert.Equal(t, tempWriter, config.GetWriter())
	config.SetWriter(os.Stdout)
}

func TestGetLimit(t *testing.T) {
	assert.Equal(t, 1000, config.GetLimit())
}

func TestSetLimit(t *testing.T) {
	config.SetLimit(2000)
	assert.Equal(t, 2000, config.GetLimit())
	config.SetLimit(1000)
}

func TestGetLogDataPath(t *testing.T) {
	assert.Equal(t, "/data/log", config.GetLogDataPath())
}

func TestSetLogDataPath(t *testing.T) {
	config.SetLogDataPath("hello")
	assert.Equal(t, "hello", config.GetLogDataPath())
	config.SetLogDataPath("/data/log")
}
