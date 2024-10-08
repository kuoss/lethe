package fileservice

import (
	"testing"

	"github.com/kuoss/common/tester"
	"github.com/kuoss/lethe/config"
	"github.com/stretchr/testify/require"
)

func TestScanner(t *testing.T) {
	_, cleanup := tester.SetupDir(t, map[string]string{
		"@/testdata/log": "data/log",
	})
	defer cleanup()

	cfg, err := config.New("test")
	require.NoError(t, err)

	fileService, err := New(cfg)
	require.NoError(t, err)
	require.NotEmpty(t, fileService)

	scanner, err := fileService.Scanner("node/node01")
	require.NoError(t, err)
	require.NotEmpty(t, scanner)
}
