package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadata(t *testing.T) {
	code, body := testGET("/api/v1/metadata")
	assert.Equal(t, 200, code)
	assert.JSONEq(t, `{
		"data": {
			"targets": [
				"node{node=\"node01\"}",
				"node{node=\"node02\"}",
				"pod{namespace=\"namespace01\"}",
				"pod{namespace=\"namespace02\"}"]},
		"status": "success"}`, body)
}
