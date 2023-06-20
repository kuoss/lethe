package fileservice

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kuoss/lethe/config"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.NotEmpty(t, fileService)
	assert.Equal(t, "filesystem", fileService.driver.Name())
	assert.Equal(t, "tmp/init", fileService.driver.RootDirectory())
}

func TestNewTemp(t *testing.T) {
	cfg, err := config.New("test")
	assert.NoError(t, err)
	assert.NotEmpty(t, cfg)

	dataPath := "tmp/temp" // caution: RemoveAll()

	cfg.SetLogDataPath(dataPath)
	tempFileService, err := New(cfg)
	assert.NoError(t, err)
	assert.NotEmpty(t, tempFileService)
	assert.DirExists(t, dataPath) // exists

	// duplicate ok
	cfg.SetLogDataPath(dataPath)
	tempFileService2, err := New(cfg)
	assert.NoError(t, err)
	assert.NotEmpty(t, tempFileService2)
	assert.DirExists(t, dataPath) // exists

	// clean up
	err = os.RemoveAll(dataPath)
	assert.NoError(t, err)
}

func TestNewFile(t *testing.T) {
	cfg, err := config.New("test")
	assert.NoError(t, err)

	cfg.SetLogDataPath(filepath.Join(".", "tmp", "writer"))

	_, err = New(cfg)

	assert.NoError(t, err)
	assert.DirExists(t, filepath.Join(".", "tmp", "writer"))
}
