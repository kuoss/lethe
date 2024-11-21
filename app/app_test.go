package app

import (
	"context"
	"testing"
	"time"

	"github.com/kuoss/common/tester"
	"github.com/stretchr/testify/require"
)

func TestRun_error_invalid(t *testing.T) {
	_, cleanup := tester.SetupDir(t, map[string]string{
		"@/testdata/etc/lethe.error.invalid.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	err := new(App).Run("test")
	require.EqualError(t, err, `failed to load configuration: Config file has error: While parsing config: yaml: did not find expected key`)
}

func TestRun_error_log_data_path(t *testing.T) {
	_, cleanup := tester.SetupDir(t, map[string]string{
		"@/testdata/etc/lethe.error.log_data_path.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	err := new(App).Run("test")
	require.EqualError(t, err, `new fileservice err: os.MkdirAll err: mkdir /dev/null: not a directory`)
}

func TestRun_smokeTest(t *testing.T) {
	_, cleanup := tester.SetupDir(t, map[string]string{})
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()
	panicChan := make(chan interface{}, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicChan <- r
			}
		}()
		err := new(App).Run("test")
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
