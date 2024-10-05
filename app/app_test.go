package app

import (
	"context"
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

func TestRun_smokeTest(t *testing.T) {
	_, cleanup := tester.MustSetupDir(t, map[string]string{})
	defer cleanup()

	app, err := New("test")
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()
	panicChan := make(chan interface{}, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicChan <- r
			}
		}()
		err := app.Run()
		require.NoError(t, err)
		close(panicChan)
	}()

	var done bool
	select {
	case <-ctx.Done():
		done = true
	case p := <-panicChan:
		t.Fatalf("panic occurred: %v", p)
	}
	require.True(t, done)
}

func TestRun_error4(t *testing.T) {
	_, cleanup := tester.MustSetupDir(t, map[string]string{
		"@/testdata/etc/lethe.error4.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	app, err := New("test")
	require.NoError(t, err)

	err = app.Run()
	require.EqualError(t, err, "listen tcp: address foo: missing port in address")
}
