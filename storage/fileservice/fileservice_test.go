package fileservice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.NotEmpty(t, fileService)
	assert.Equal(t, "filesystem", fileService.driver.Name())
	assert.Equal(t, "tmp/init", fileService.driver.RootDirectory())
}
