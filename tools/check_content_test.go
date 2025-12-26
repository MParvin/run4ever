package tools

import (
	"os"
	"testing"
)

func TestIsEmpty(t *testing.T) {
	tempDir := t.TempDir()
	testFile := tempDir + "/test_empty.txt"

	// Test empty file
	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	f.Close()

	if !IsEmpty(testFile) {
		t.Error("IsEmpty should return true for empty file")
	}

	// Test file with content
	err = os.WriteFile(testFile, []byte("content"), 0644)
	if err != nil {
		t.Fatalf("Failed to write to test file: %v", err)
	}

	if IsEmpty(testFile) {
		t.Error("IsEmpty should return false for non-empty file")
	}

	// Test non-existent file (should return false)
	nonExistentFile := tempDir + "/nonexistent.txt"
	if IsEmpty(nonExistentFile) {
		t.Error("IsEmpty should return false for non-existent file")
	}
}

func TestIsEmptyWithError(t *testing.T) {
	tempDir := t.TempDir()
	testFile := tempDir + "/test_empty.txt"

	// Test empty file
	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	f.Close()

	empty, err := IsEmptyWithError(testFile)
	if err != nil {
		t.Fatalf("IsEmptyWithError should not error for existing file: %v", err)
	}
	if !empty {
		t.Error("IsEmptyWithError should return true for empty file")
	}

	// Test file with content
	err = os.WriteFile(testFile, []byte("content"), 0644)
	if err != nil {
		t.Fatalf("Failed to write to test file: %v", err)
	}

	empty, err = IsEmptyWithError(testFile)
	if err != nil {
		t.Fatalf("IsEmptyWithError should not error for existing file: %v", err)
	}
	if empty {
		t.Error("IsEmptyWithError should return false for non-empty file")
	}

	// Test non-existent file
	nonExistentFile := tempDir + "/nonexistent.txt"
	empty, err = IsEmptyWithError(nonExistentFile)
	if err == nil {
		t.Error("IsEmptyWithError should return error for non-existent file")
	}
	if empty {
		t.Error("IsEmptyWithError should return false for non-existent file")
	}
}
