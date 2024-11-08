package config

import (
	"testing"
	"time"

	"github.com/kuoss/common/tester"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	want := &Config{
		Version: "test",
		Query: Query{
			Limit:   1000,
			Timeout: 20 * time.Second,
		},
		Retention: Retention{
			Size:             0,
			Time:             15 * 24 * time.Hour,
			SizeStrategy:     "file",
			RotationInterval: 20 * time.Second,
		},
		Storage: Storage{LogDataPath: "data/log"},
		Web: Web{
			ListenAddress: ":6060",
			GinMode:       "release",
		},
	}

	got, err := New("test")
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestNew_example(t *testing.T) {
	_, cleanup := tester.SetupDir(t, map[string]string{
		"@/etc/lethe.example.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	want := &Config{
		Version: "test",
		Query:   Query{Limit: 1000, Timeout: 20 * time.Second},
		Retention: Retention{
			Size:             1000 * 1024 * 1024 * 1024, // 1000GB
			Time:             15 * 24 * time.Hour,       // 15d
			SizeStrategy:     "file",
			RotationInterval: 20 * time.Second,
		},
		Storage: Storage{LogDataPath: "/data/log"},
		Web: Web{
			ListenAddress: ":6060",
			GinMode:       "release",
		},
	}

	got, err := New("test")
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestNew_ok1(t *testing.T) {
	_, cleanup := tester.SetupDir(t, map[string]string{
		"@/testdata/etc/lethe.ok1.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	want := &Config{
		Version: "ok1",
		Query:   Query{Limit: 1000, Timeout: 20 * time.Second},
		Retention: Retention{
			Size:             200 * 1024 * 1024 * 1024, // 200GB
			Time:             24 * time.Hour,           // 1d
			SizeStrategy:     "file",
			RotationInterval: 20 * time.Second,
		},
		Storage: Storage{LogDataPath: "/tmp/log"},
		Web: Web{
			ListenAddress: ":6060",
			GinMode:       "release",
		},
	}

	got, err := New("ok1")
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestNew_ok2(t *testing.T) {
	_, cleanup := tester.SetupDir(t, map[string]string{
		"@/testdata/etc/lethe.ok2.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	want := &Config{
		Version: "test2",
		Query:   Query{Limit: 1000, Timeout: 20 * time.Second},
		Retention: Retention{
			Size:             300 * 1024 * 1024 * 1024, // 300GB
			Time:             15 * 24 * time.Hour,      // 15d
			SizeStrategy:     "file",
			RotationInterval: 20 * time.Second,
		},
		Storage: Storage{LogDataPath: "/tmp/log"},
		Web: Web{
			ListenAddress: ":6060",
			GinMode:       "release",
		},
	}

	got, err := New("test2")
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestNew_legacy(t *testing.T) {
	_, cleanup := tester.SetupDir(t, map[string]string{
		"@/testdata/etc/lethe.legacy.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	want := &Config{
		Version: "test2",
		Query:   Query{Limit: 1000, Timeout: 20 * time.Second},
		Retention: Retention{
			Size:             100 * 1024 * 1024 * 1024, // 100GB
			Time:             15 * 24 * time.Hour,      // 15d
			SizeStrategy:     "file",
			RotationInterval: 20 * time.Second,
		},
		Storage: Storage{LogDataPath: "/tmp/log"},
		Web: Web{
			ListenAddress: ":6060",
			GinMode:       "release",
		},
	}

	got, err := New("test2")
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestNew_error_invalid(t *testing.T) {
	_, cleanup := tester.SetupDir(t, map[string]string{
		"@/testdata/etc/lethe.error.invalid.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	got, err := New("test")
	require.EqualError(t, err, `Config file has error: While parsing config: yaml: did not find expected key`)
	require.Nil(t, got)
}

func TestNew_error_size(t *testing.T) {
	_, cleanup := tester.SetupDir(t, map[string]string{
		"@/testdata/etc/lethe.error.size.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	got, err := New("test")
	require.EqualError(t, err, `parse retention size err: failed to parse size '200x': units: unknown unit x in 200x, size: 200x`)
	require.Nil(t, got)
}

func TestNew_error_time(t *testing.T) {
	_, cleanup := tester.SetupDir(t, map[string]string{
		"@/testdata/etc/lethe.error.time.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	got, err := New("test")
	require.EqualError(t, err, `parse retention time err: not a valid duration string: "1", time: 1`)
	require.Nil(t, got)
}
