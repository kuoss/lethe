package app

import (
	"testing"

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
		Retention: config.Retention{SizingStrategy: "file", Size: 214748364800, Time: 86400000000000},
		Rotator:   config.Rotator{Interval: 20000000000},
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

// func TestMustConfig(t *testing.T) {

// }

// func TestMustFileService(t *testing.T) {
// 	defer cleanup()

// 	cfg := mustConfig("test")
// 	fileService := mustFileService(cfg)
// 	require.NotNil(t, fileService, "Expected file service to be non-nil")
// }

// func TestStartRotator(t *testing.T) {
// 	_, cleanup := tester.MustSetupDir(t, map[string]string{
// 		"@/testdata/etc/lethe.main.yaml": "etc/lethe.yaml",
// 	})
// 	defer cleanup()

// 	version := "test-version"
// 	cfg := mustConfig(version)
// 	fileService := mustFileService(cfg)
// 	require.NotPanics(t, func() {
// 		startRotator(cfg, fileService)
// 	}, "Expected startRotator to not panic")
// }
