package router

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTarget(t *testing.T) {
	code, body, cleanup := testGET(t, "/api/v1/targets")
	defer cleanup()
	require.Equal(t, 200, code)
	require.JSONEq(t, `{"data":{"activeTargets":[
		{"discoveredLabels":{"__meta_kubernetes_node_name":"node01","job":"node"},"health":"up","lastScrape":"2009-11-10T23:00:00Z"},
		{"discoveredLabels":{"__meta_kubernetes_node_name":"node02","job":"node"},"health":"up","lastScrape":"2009-11-10T21:58:00Z"},
		{"discoveredLabels":{"__meta_kubernetes_namespace":"namespace01","job":"pod"},"health":"up","lastScrape":"2009-11-10T23:00:00Z"},
		{"discoveredLabels":{"__meta_kubernetes_namespace":"namespace02","job":"pod"},"health":"up","lastScrape":"2009-11-10T22:58:00Z"}]},"status":"success"}`, body)
}
