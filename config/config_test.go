package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	config1 *Config
)

func init() {
	var err error
	config1, err = New("test")
	if err != nil {
		panic(err)
	}
}

func TestNew(t *testing.T) {
	assert.NotEmpty(t, config1)
}

func TestLimit(t *testing.T) {
	assert.Equal(t, 1000, config1.Limit())
}

func TestLogDataPath(t *testing.T) {
	assert.Equal(t, "/data/log", config1.LogDataPath())
	config1.SetLogDataPath("tmp/hello")
	assert.Equal(t, "tmp/hello", config1.LogDataPath())
	config1.SetLogDataPath("/data/log")
}

func TestRetentionSize(t *testing.T) {
	assert.Equal(t, 100*1024*1024, config1.RetentionSize()) // 100 MiB
	config1.SetRetentionSize(1 * 1024 * 1024)               // 1 MiB
	assert.Equal(t, 1*1024*1024, config1.RetentionSize())   // 1 MiB
	config1.SetRetentionSize(100 * 1024 * 1024)             // 100 MiB
}
func TestRetentionTime(t *testing.T) {
	assert.Equal(t, 24*time.Hour, config1.RetentionTime())    // 1 day
	config1.SetRetentionTime(15 * 24 * time.Hour)             // 15 days
	assert.Equal(t, 15*24*time.Hour, config1.RetentionTime()) // 15 days
	config1.SetRetentionTime(24 * time.Hour)                  // 1 day
}

func TestRetentionSizingStrategy(t *testing.T) {
	assert.Equal(t, "file", config1.RetentionSizingStrategy())
}

func TestVersion(t *testing.T) {
	assert.Equal(t, "test", config1.Version())
}

func TestWebListenAddress(t *testing.T) {
	assert.Equal(t, ":6060", config1.WebListenAddress())
}
