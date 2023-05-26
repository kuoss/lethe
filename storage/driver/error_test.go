package driver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	err := PathNotFoundError{Path: "/tmp"}
	got := err.Error()
	assert.Equal(t, "Path not found: /tmp", got)
}
