// Test for mkdir.go
package tools

import (
	"os"
	"testing"
)

func TestCreateDir(t *testing.T) {
	// Use a temporary directory for testing
	tempDir := t.TempDir()
	testDir := tempDir + "/.run4ever_test"

	CreateDir(testDir)
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Fatalf("Directory was not created: %v", err)
	}

	// Test creating the same directory again (should not error)
	CreateDir(testDir)
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Fatalf("Directory should still exist: %v", err)
	}
}

func TestCreateDirNested(t *testing.T) {
	tempDir := t.TempDir()

	// First create parent directories
	parentDir := tempDir + "/level1/level2"
	err := os.MkdirAll(parentDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create parent directories: %v", err)
	}

	nestedDir := parentDir + "/.run4ever_test"

	CreateDir(nestedDir)
	if _, err := os.Stat(nestedDir); os.IsNotExist(err) {
		t.Fatalf("Nested directory was not created: %v", err)
	}
}
