package fileservice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	assert.NotEmpty(t, fileService)
	assert.Equal(t, "filesystem", fileService.driver.Name())
	assert.Equal(t, "tmp/init", fileService.driver.RootDirectory())
}
