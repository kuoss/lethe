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
			SizingStrategy:   "file",
			RotationInterval: 20 * time.Second,
		},
		Storage: Storage{LogDataPath: "data/log"},
		Web:     Web{ListenAddress: ":6060"},
	}

	got, err := New("test")
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestNew_example(t *testing.T) {
	_, cleanup := tester.MustSetupDir(t, map[string]string{
		"@/etc/lethe.example.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	want := &Config{
		Version: "example",
		Query:   Query{Limit: 1000, Timeout: 20 * time.Second},
		Retention: Retention{
			Size:             1000 * 1024 * 1024 * 1024, // 1000GB
			Time:             15 * 24 * time.Hour,       // 15d
			SizingStrategy:   "file",
			RotationInterval: 20 * time.Second,
		},
		Storage: Storage{LogDataPath: "/data/log"},
		Web:     Web{ListenAddress: ":6060"},
	}

	got, err := New("example")
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestNew_ok1(t *testing.T) {
	_, cleanup := tester.MustSetupDir(t, map[string]string{
		"@/testdata/etc/lethe.ok1.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	want := &Config{
		Version: "ok1",
		Query:   Query{Limit: 1000, Timeout: 20 * time.Second},
		Retention: Retention{
			Size:             200 * 1024 * 1024 * 1024, // 200GB
			Time:             24 * time.Hour,           // 1d
			SizingStrategy:   "file",
			RotationInterval: 20 * time.Second,
		},
		Storage: Storage{LogDataPath: "/data/log"},
		Web:     Web{ListenAddress: ":6060"},
	}

	got, err := New("ok1")
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestNew_ok2(t *testing.T) {
	_, cleanup := tester.MustSetupDir(t, map[string]string{
		"@/testdata/etc/lethe.ok2.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	want := &Config{
		Version: "test2",
		Query:   Query{Limit: 1000, Timeout: 20 * time.Second},
		Retention: Retention{
			Size:             300 * 1024 * 1024 * 1024, // 300GB
			Time:             15 * 24 * time.Hour,      // 15d
			SizingStrategy:   "file",
			RotationInterval: 20 * time.Second,
		},
		Storage: Storage{LogDataPath: "/var/data/log"},
		Web:     Web{ListenAddress: ":6060"},
	}

	got, err := New("test2")
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestNew_legacy(t *testing.T) {
	_, cleanup := tester.MustSetupDir(t, map[string]string{
		"@/testdata/etc/lethe.legacy.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	want := &Config{
		Version: "test2",
		Query:   Query{Limit: 1000, Timeout: 20 * time.Second},
		Retention: Retention{
			Size:             100 * 1024 * 1024 * 1024, // 100GB
			Time:             15 * 24 * time.Hour,      // 15d
			SizingStrategy:   "file",
			RotationInterval: 20 * time.Second,
		},
		Storage: Storage{LogDataPath: "/data/log"},
		Web:     Web{ListenAddress: ":6060"},
	}

	got, err := New("test2")
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestNew_error1(t *testing.T) {
	_, cleanup := tester.MustSetupDir(t, map[string]string{
		"@/testdata/etc/lethe.error1.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	got, err := New("error")
	require.Error(t, err)
	require.Nil(t, got)
}

func TestNew_error2(t *testing.T) {
	_, cleanup := tester.MustSetupDir(t, map[string]string{
		"@/testdata/etc/lethe.error2.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	got, err := New("error")
	require.Error(t, err)
	require.Nil(t, got)
}

func TestNew_error3(t *testing.T) {
	_, cleanup := tester.MustSetupDir(t, map[string]string{
		"@/testdata/etc/lethe.error3.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	got, err := New("error")
	require.Error(t, err)
	require.Nil(t, got)
}
