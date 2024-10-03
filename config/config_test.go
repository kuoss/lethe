package config

import (
	"testing"
	"time"

	"github.com/kuoss/common/tester"
	"github.com/stretchr/testify/assert"
)

func TestNew_default(t *testing.T) {
	want := &Config{
		Version:                 "test",
		limit:                   1000,
		logDataPath:             "/var/data/log",
		retentionSize:           1000 * 1024 * 1024 * 1024,
		retentionTime:           15 * 24 * time.Hour,
		retentionSizingStrategy: "file",
		timeout:                 20 * time.Second,
		webListenAddress:        ":6060",
	}

	got, err := New("test")
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestNew_example(t *testing.T) {
	_, cleanup := tester.MustSetupDir(t, map[string]string{
		"@/etc/lethe.example.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	want := &Config{
		Version:                 "example",
		limit:                   1000,
		logDataPath:             "/data/log",
		retentionSize:           100 * 1024 * 1024,
		retentionTime:           24 * time.Hour,
		retentionSizingStrategy: "file",
		timeout:                 20 * time.Second,
		webListenAddress:        ":6060",
	}

	got, err := New("example")
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestNew_ok1(t *testing.T) {
	_, cleanup := tester.MustSetupDir(t, map[string]string{
		"@/testdata/etc/lethe.ok1.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	want := &Config{
		Version:                 "ok1",
		limit:                   1000,
		logDataPath:             "/data/log",
		retentionSize:           200 * 1024 * 1024,
		retentionTime:           24 * time.Hour,
		retentionSizingStrategy: "file",
		timeout:                 20 * time.Second,
		webListenAddress:        ":6060",
	}

	got, err := New("ok1")
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestNew_ok2(t *testing.T) {
	_, cleanup := tester.MustSetupDir(t, map[string]string{
		"@/testdata/etc/lethe.ok2.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	want := &Config{
		Version:                 "test2",
		limit:                   1000,
		logDataPath:             "/var/data/log",
		retentionSize:           300 * 1024 * 1024,
		retentionTime:           15 * 24 * time.Hour, // Only retentionTime feild is initialized with default value
		retentionSizingStrategy: "file",
		timeout:                 20 * time.Second,
		webListenAddress:        ":6060",
	}

	got, err := New("test2")
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestNew_error1(t *testing.T) {
	_, cleanup := tester.MustSetupDir(t, map[string]string{
		"@/testdata/etc/lethe.error1.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	got, err := New("error")
	assert.EqualError(t, err, `readInConfig err: While parsing config: yaml: did not find expected key`)
	assert.Nil(t, got)
}

func TestNew_error2(t *testing.T) {
	_, cleanup := tester.MustSetupDir(t, map[string]string{
		"@/testdata/etc/lethe.error2.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	got, err := New("error")
	assert.EqualError(t, err, "stringToBytes err: cannot accept unit '0' in '200''. allowed units: [k, m, g]")
	assert.Nil(t, got)
}

func TestNew_error3(t *testing.T) {
	_, cleanup := tester.MustSetupDir(t, map[string]string{
		"@/testdata/etc/lethe.error3.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	got, err := New("error")
	assert.EqualError(t, err, "getDurationFromAge err: not a valid duration string: \"1\"")
	assert.Nil(t, got)
}

func TestFunctions(t *testing.T) {
	cfg, err := New("test")
	assert.NoError(t, err)

	assert.Equal(t, "test", cfg.Version)

	assert.Equal(t, 1000, cfg.Limit())

	assert.Equal(t, "/var/data/log", cfg.LogDataPath())
	cfg.SetLogDataPath("tmp/hello")
	assert.Equal(t, "tmp/hello", cfg.LogDataPath())
	cfg.SetLogDataPath("/data/log")

	assert.Equal(t, 1000*1024*1024*1024, cfg.RetentionSize()) // 1000 GiB
	cfg.SetRetentionSize(1 * 1024 * 1024)                     // 1 MiB
	assert.Equal(t, 1*1024*1024, cfg.RetentionSize())         // 1 MiB
	cfg.SetRetentionSize(100 * 1024 * 1024)                   // 100 MiB

	assert.Equal(t, 15*24*time.Hour, cfg.RetentionTime()) // 15 day
	cfg.SetRetentionTime(1 * 24 * time.Hour)              // 1 days
	assert.Equal(t, 1*24*time.Hour, cfg.RetentionTime())  // 1 days
	cfg.SetRetentionTime(15 * 24 * time.Hour)             // 15 day

	assert.Equal(t, 20*time.Second, cfg.Timeout())

	assert.Equal(t, "file", cfg.RetentionSizingStrategy())
	assert.Equal(t, ":6060", cfg.WebListenAddress())
}
