package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockApp struct {
	runFunc func(version string) error
}

func (m *MockApp) Run(version string) error {
	return m.runFunc(version)
}

func TestMainFunctionExitCode(t *testing.T) {
	testCases := []struct {
		name         string
		mockRunFunc  func(version string) error
		wantExitCode int
	}{
		{
			name: "Successful exit",
			mockRunFunc: func(version string) error {
				return nil
			},
			wantExitCode: 0,
		},
		{
			name: "Error exit",
			mockRunFunc: func(version string) error {
				return errors.New("run error")
			},
			wantExitCode: 1,
		},
	}

	originalApp := ap
	originalExit := exit
	defer func() {
		ap = originalApp
		exit = originalExit
	}()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ap = &MockApp{runFunc: tc.mockRunFunc}
			var gotExitCode int
			exit = func(code int) {
				gotExitCode = code
			}

			main()

			assert.Equal(t, tc.wantExitCode, gotExitCode)
		})
	}
}
