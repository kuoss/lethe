package config

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultNew(t *testing.T) {
	cfg, err := New("test")
	if err != nil {
		t.Error(err)
	}
	expected := &Config{
		limit:                   1000,
		logDataPath:             "./tmp/log",
		retentionSize:           100 * 1024 * 1024,
		retentionTime:           15 * 24 * time.Hour,
		retentionSizingStrategy: "file",
		timeout:                 20 * time.Second,
		version:                 "test",
		webListenAddress:        ":6060",
	}
	assert.Equal(t, expected, cfg)
}

var customConfigYaml = `
retention:
  time: 1d
  size: 200m
storage:
  driver: filesystem
  log_data_path: /data/log
`

func TestNewFromFile(t *testing.T) {

	appFS := afero.NewOsFs()
	err := afero.WriteFile(appFS, filepath.Join("..", "etc", "lethe.yaml"), []byte(customConfigYaml), 0644)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(filepath.Join("..", "etc", "lethe.yaml"))

	path := filepath.Join("..", "etc", "lethe.yaml")
	_, err = appFS.Stat(path)

	if os.IsNotExist(err) {
		t.Errorf("file %q does not exist.\n", path)
	}

	cfg, err := New("test")
	if err != nil {
		t.Error(err)
	}
	expected := &Config{
		limit:                   1000,
		logDataPath:             "/data/log",
		retentionSize:           200 * 1024 * 1024,
		retentionTime:           24 * time.Hour,
		retentionSizingStrategy: "file",
		timeout:                 20 * time.Second,
		version:                 "test",
		webListenAddress:        ":6060",
	}
	assert.Equal(t, expected, cfg)

}

var customConfigYaml2 = `
retention:
  size: 300m
storage:
  driver: filesystem
  log_data_path: /var/data/log
`

func TestNewFromFilePartial(t *testing.T) {

	appFS := afero.NewOsFs()
	err := afero.WriteFile(appFS, filepath.Join("..", "etc", "lethe.yaml"), []byte(customConfigYaml2), 0644)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(filepath.Join("..", "etc", "lethe.yaml"))

	path := filepath.Join("..", "etc", "lethe.yaml")
	_, err = appFS.Stat(path)

	if os.IsNotExist(err) {
		t.Errorf("file %q does not exist.\n", path)
	}

	cfg, err := New("test")
	if err != nil {
		t.Error(err)
	}
	expected := &Config{
		limit:                   1000,
		logDataPath:             "/var/data/log",
		retentionSize:           300 * 1024 * 1024,
		retentionTime:           15 * 24 * time.Hour, // Only retentionTime feild is initialized with default value
		retentionSizingStrategy: "file",
		timeout:                 20 * time.Second,
		version:                 "test",
		webListenAddress:        ":6060",
	}
	assert.Equal(t, expected, cfg)
}

func TestExportedFunction(t *testing.T) {
	cfg, err := New("DEFAULT")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 1000, cfg.Limit())

	assert.Equal(t, "./tmp/log", cfg.LogDataPath())
	cfg.SetLogDataPath("tmp/hello")
	assert.Equal(t, "tmp/hello", cfg.LogDataPath())
	cfg.SetLogDataPath("/data/log")

	assert.Equal(t, 100*1024*1024, cfg.RetentionSize()) // 100 MiB
	cfg.SetRetentionSize(1 * 1024 * 1024)               // 1 MiB
	assert.Equal(t, 1*1024*1024, cfg.RetentionSize())   // 1 MiB
	cfg.SetRetentionSize(100 * 1024 * 1024)             // 100 MiB

	assert.Equal(t, 15*24*time.Hour, cfg.RetentionTime()) // 15 day
	cfg.SetRetentionTime(1 * 24 * time.Hour)              // 1 days
	assert.Equal(t, 1*24*time.Hour, cfg.RetentionTime())  // 1 days
	cfg.SetRetentionTime(15 * 24 * time.Hour)             // 15 day

	assert.Equal(t, "file", cfg.RetentionSizingStrategy())
	assert.Equal(t, "DEFAULT", cfg.Version())
	assert.Equal(t, ":6060", cfg.WebListenAddress())
}
