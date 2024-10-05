package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type MockApp struct {
	newErr error
	runErr error
}

func (m *MockApp) New(version string) error {
	return m.newErr
}

func (m *MockApp) Run() error {
	return m.runErr
}

func TestMainFunc(t *testing.T) {
	originalApp := myApp
	originalExit := exit
	defer func() {
		myApp = originalApp
		exit = originalExit
	}()

	var gotExitCode int
	exit = func(code int) {
		gotExitCode = code
	}

	t.Run("ok", func(t *testing.T) {
		myApp = &MockApp{}
		main()
		require.Equal(t, 0, gotExitCode)
	})
	t.Run("new error", func(t *testing.T) {
		myApp = &MockApp{newErr: errors.New("fake new error")}
		main()
		require.Equal(t, 1, gotExitCode)
	})

	t.Run("run error", func(t *testing.T) {
		myApp = &MockApp{runErr: errors.New("fake run error")}
		main()
		require.Equal(t, 1, gotExitCode)
	})
}
