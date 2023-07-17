package fileservice

import (
	"os"
	"testing"

	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/util/testutil"
	"github.com/stretchr/testify/require"
	// "github.com/stretchr/testify/require"
)

var (
	fileService_clean *FileService
	logDataPath_clean string = "tmp/storage_fileservice_clean_test"
)

func init() {
	testutil.ChdirRoot()
	testutil.ResetLogData()

	cfg, err := config.New("test")
	if err != nil {
		panic(err)
	}
	cfg.SetLogDataPath(logDataPath_clean)
	fileService_clean, err = New(cfg)
	if err != nil {
		panic(err)
	}
}

func TestClean(t *testing.T) {
	// setup
	var err error
	err = os.WriteFile("tmp/storage_fileservice_clean_test/kube.1", []byte("hello"), 0644)
	require.NoError(t, err)
	err = os.WriteFile("tmp/storage_fileservice_clean_test/host.1", []byte("hello"), 0644)
	require.NoError(t, err)

	require.FileExists(t, "tmp/storage_fileservice_clean_test/kube.1")
	require.FileExists(t, "tmp/storage_fileservice_clean_test/host.1")
	fileService_clean.Clean()
	require.NoFileExists(t, "tmp/storage_fileservice_clean_test/kube.1")
	require.NoFileExists(t, "tmp/storage_fileservice_clean_test/host.1")
}
