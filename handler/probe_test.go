package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthy(t *testing.T) {
	code, body := testGET("/-/healthy")
	assert.Equal(t, 200, code)
	assert.Equal(t, "Venti is Healthy.\n", body)
}

func TestReady(t *testing.T) {
	code, body := testGET("/-/ready")
	assert.Equal(t, 200, code)
	assert.Equal(t, "Venti is Ready.\n", body)
}
