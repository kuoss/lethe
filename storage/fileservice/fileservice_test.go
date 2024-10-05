package fileservice

import (
	"testing"

	"github.com/kuoss/lethe/config"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cfg, err := config.New("test")
	require.NoError(t, err)

	fileService, err := New(cfg)
	require.NoError(t, err)
	require.NotEmpty(t, fileService)
	require.Equal(t, cfg, fileService.Config)
}
