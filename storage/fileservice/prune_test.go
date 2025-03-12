package fileservice

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/kuoss/common/tester"
	"github.com/kuoss/lethe/config"
	"github.com/stretchr/testify/require"
)

func TestPrune(t *testing.T) {
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
	fileService.Prune()
	require.NoFileExists(t, tmpDir+"/data/log/kube.1")
	require.NoFileExists(t, tmpDir+"/data/log/host.1")
}

func TestIsOldEmptyDir(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after test

	t.Run("Test empty directory older than 24 hours", func(t *testing.T) {
		// Change modification time to older than 24 hours
		oldTime := time.Now().Add(-25 * time.Hour)
		if err := os.Chtimes(tempDir, oldTime, oldTime); err != nil {
			t.Fatalf("Failed to change directory time: %v", err)
		}

		// Call the function
		isOldEmpty, err := isOldEmptyDir(tempDir)
		if err != nil {
			t.Fatalf("Error checking directory: %v", err)
		}
		if !isOldEmpty {
			t.Errorf("Expected directory to be old and empty, got: %v", isOldEmpty)
		}
	})

	t.Run("Test non-empty directory", func(t *testing.T) {
		// Create a new file in the directory
		file, err := os.MkdirTemp(tempDir, "testfile")
		if err != nil {
			t.Fatalf("Failed to create file in temp directory: %v", err)
		}
		defer os.RemoveAll(file) // Clean up after test

		// Call the function
		isOldEmpty, err := isOldEmptyDir(tempDir)
		if err != nil {
			t.Fatalf("Error checking directory: %v", err)
		}
		if isOldEmpty {
			t.Errorf("Expected directory to be non-empty, got: %v", isOldEmpty)
		}
	})

	t.Run("Test empty directory newer than 24 hours", func(t *testing.T) {
		// Remove test file if created previously
		files, _ := os.ReadDir(tempDir)
		for _, f := range files {
			os.Remove(fmt.Sprintf("%s/%s", tempDir, f.Name()))
		}

		// Set the directory modification time to now
		now := time.Now()
		if err := os.Chtimes(tempDir, now, now); err != nil {
			t.Fatalf("Failed to change directory time: %v", err)
		}

		// Call the function
		isOldEmpty, err := isOldEmptyDir(tempDir)
		if err != nil {
			t.Fatalf("Error checking directory: %v", err)
		}
		if isOldEmpty {
			t.Errorf("Expected directory to be not old, got: %v", isOldEmpty)
		}
	})
}
