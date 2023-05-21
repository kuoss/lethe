package handler

import (
	"testing"

	"github.com/kuoss/lethe/storage/fileservice"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	handler1 := New(&fileservice.FileService{})
	assert.NotEmpty(t, handler1)
}
