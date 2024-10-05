package router

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetadata(t *testing.T) {
	code, body, cleanup := testGET(t, "/api/v1/metadata")
	defer cleanup()

	require.Equal(t, 200, code)
	require.JSONEq(t, `{
		"data": {
			"targets": [
				"node{node=\"node01\"}",
				"node{node=\"node02\"}",
				"pod{namespace=\"namespace01\"}",
				"pod{namespace=\"namespace02\"}"]},
		"status": "success"}`, body)
}
