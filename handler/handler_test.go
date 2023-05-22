package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.NotEmpty(t, handler1)
}

func TestPingRoute(t *testing.T) {
	code, body := testGET("/ping")
	assert.Equal(t, 200, code)
	assert.Equal(t, `{"message":"pong"}`, body)
}
