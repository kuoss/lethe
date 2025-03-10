package fileservice

import (
	"fmt"
	"io/fs"
	"os"
	"testing"
)

func TestIsEmpty(t *testing.T) {
	// Set up a test directory
	testDir := "testDir"
	emptyDir := testDir + "/empty"
	nonEmptyDir := testDir + "/nonEmpty"

	// Clean up any previous test artifacts
	os.RemoveAll(testDir)

	// Create test directories
	err := os.MkdirAll(emptyDir, fs.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create empty test directory: %v", err)
	}

	err = os.MkdirAll(nonEmptyDir, fs.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create non-empty test directory: %v", err)
	}

	// Add a file to the non-empty directory
	err = os.WriteFile(nonEmptyDir+"/testfile.txt", []byte("content"), fs.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create file in non-empty test directory: %v", err)
	}

	// Test case for empty directory
	t.Run("Empty directory", func(t *testing.T) {
		empty, err := isEmpty(emptyDir)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		if !empty {
			t.Errorf("Expected empty directory to return true, but got false")
		}
	})

	// Test case for non-empty directory
	t.Run("Non-empty directory", func(t *testing.T) {
		empty, err := isEmpty(nonEmptyDir)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		if empty {
			t.Errorf("Expected non-empty directory to return false, but got true")
		}
	})

	// Test case for non-existent directory
	t.Run("Non-existent directory", func(t *testing.T) {
		_, err := isEmpty("nonExistentDir")
		if err == nil {
			t.Errorf("Expected error for non-existent directory, but got none")
		}
	})

	// Clean up test directories
	err = os.RemoveAll(testDir)
	if err != nil {
		fmt.Printf("Warning: failed to clean up test directories: %v", err)
	}
}
