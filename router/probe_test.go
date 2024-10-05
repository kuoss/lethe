package router

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHealthy(t *testing.T) {
	code, body, cleanup := testGET(t, "/-/healthy")
	defer cleanup()
	require.Equal(t, 200, code)
	require.Equal(t, "Venti is Healthy.\n", body)
}

func TestReady(t *testing.T) {
	code, body, cleanup := testGET(t, "/-/ready")
	defer cleanup()
	require.Equal(t, 200, code)
	require.Equal(t, "Venti is Ready.\n", body)
}
