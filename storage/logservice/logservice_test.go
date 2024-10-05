package logservice

import (
	"testing"

	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/storage/fileservice"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cfg, err := config.New("test")
	require.NoError(t, err)
	fileService, err := fileservice.New(cfg)
	require.NoError(t, err)
	logService := New(cfg, fileService)
	require.NotEmpty(t, logService)
}
