package driver

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	testCases := []struct {
		path      string
		err       error
		wantError string
	}{
		{
			"", nil,
			"Path not found: , err: %!s(<nil>)",
		},
		{
			"", errors.New("hello"),
			"Path not found: , err: hello",
		},
		{
			"/tmp", nil,
			"Path not found: /tmp, err: %!s(<nil>)",
		},
		{
			"/tmp", errors.New("hello"),
			"Path not found: /tmp, err: hello",
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			err := PathNotFoundError{Path: tc.path, Err: tc.err}
			assert.EqualError(t, err, tc.wantError)
		})
	}
}
