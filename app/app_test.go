package app

import (
	"testing"
	"time"

	"github.com/kuoss/common/tester"
	"github.com/kuoss/lethe/config"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	_, cleanup := tester.MustSetupDir(t, map[string]string{
		"@/testdata/etc/lethe.main.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	wantConfig := &config.Config{
		Version:   "test",
		Query:     config.Query{Limit: 1000, Timeout: 20000000000},
		Retention: config.Retention{SizingStrategy: "file", Size: 214748364800, Time: 86400000000000, RotationInterval: 20 * time.Second},
		Storage:   config.Storage{LogDataPath: "/tmp"},
		Web:       config.Web{ListenAddress: ":6060"},
	}

	app, err := New("test")
	require.NoError(t, err)
	require.NotEmpty(t, app)
	require.NotEmpty(t, app.rotator)
	require.NotEmpty(t, app.router)
	require.Equal(t, wantConfig, app.Config)
}

func TestNew_error1(t *testing.T) {
	_, cleanup := tester.MustSetupDir(t, map[string]string{
		"@/testdata/etc/lethe.error1.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	app, err := New("test")
	require.Error(t, err)
	require.Nil(t, app)
}

func TestNew_error2(t *testing.T) {
	_, cleanup := tester.MustSetupDir(t, map[string]string{
		"@/testdata/etc/lethe.error1.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	app, err := New("test")
	require.Error(t, err)
	require.Nil(t, app)
}
