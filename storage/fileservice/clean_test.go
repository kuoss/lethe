package fileservice

import (
	"os"
	"testing"

	"github.com/kuoss/common/tester"
	"github.com/kuoss/lethe/config"
	"github.com/stretchr/testify/require"
)

func TestClean(t *testing.T) {
	tmpDir, cleanup := tester.SetupDir(t, map[string]string{
		"@/testdata/log": "data/log",
	})
	defer cleanup()

	cfg, err := config.New("test")
	require.NoError(t, err)
	fileService, err := New(cfg)
	require.NoError(t, err)

	// // setup
	err = os.WriteFile(tmpDir+"/data/log/kube.1", []byte("hello"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(tmpDir+"/data/log/host.1", []byte("hello"), 0644)
	require.NoError(t, err)

	require.FileExists(t, tmpDir+"/data/log/kube.1")
	require.FileExists(t, tmpDir+"/data/log/host.1")
	fileService.Clean()
	require.NoFileExists(t, tmpDir+"/data/log/kube.1")
	require.NoFileExists(t, tmpDir+"/data/log/host.1")
}
