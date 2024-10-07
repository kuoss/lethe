package app

import (
	"context"
	"testing"
	"time"

	"github.com/kuoss/common/tester"
	"github.com/stretchr/testify/require"
)

func TestRun_error1(t *testing.T) {
	_, cleanup := tester.SetupDir(t, map[string]string{
		"@/testdata/etc/lethe.error1.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	err := Run("test")
	require.Error(t, err)
}

func TestRun_error2(t *testing.T) {
	_, cleanup := tester.SetupDir(t, map[string]string{
		"@/testdata/etc/lethe.error1.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	err := Run("test")
	require.Error(t, err)
}

func TestRun_error4(t *testing.T) {
	_, cleanup := tester.SetupDir(t, map[string]string{
		"@/testdata/etc/lethe.error4.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	err := Run("test")
	require.Error(t, err)
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
		err := Run("test")
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
