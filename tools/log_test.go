package tools

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestLogFile(t *testing.T) string {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "run4ever_test.state")

	// Override HOME env for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})

	return logFile
}

func TestWriteHeader(t *testing.T) {
	logFile := setupTestLogFile(t)

	WriteHeader(logFile)

	content, err := ioutil.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	expected := "Time \t\t\t | Job-ID \t\t | PID \t\t | Command \t | Args \t\t | Status\n"
	if string(content) != expected {
		t.Errorf("WriteHeader failed. Expected: %q, Got: %q", expected, string(content))
	}
}

func TestLog(t *testing.T) {
	logFile := setupTestLogFile(t)
	WriteHeader(logFile)

	command := "test_command"
	args := []string{"arg1", "arg2"}
	pid := 12345
	jobID := "test-job-id-12345"
	name := "test-job"

	LogWithFile(command, args, pid, jobID, name, logFile)

	content, err := ioutil.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) < 2 {
		t.Fatal("Expected at least header and one log entry")
	}

	logLine := lines[1]
	if !strings.Contains(logLine, fmt.Sprintf("%d", pid)) {
		t.Errorf("Log entry should contain PID %d", pid)
	}
	if !strings.Contains(logLine, command) {
		t.Errorf("Log entry should contain command %s", command)
	}
	if !strings.Contains(logLine, jobID) {
		t.Errorf("Log entry should contain job ID %s", jobID)
	}
	for _, arg := range args {
		if !strings.Contains(logLine, arg) {
			t.Errorf("Log entry should contain argument %s", arg)
		}
	}
}

func TestDeleteLog(t *testing.T) {
	logFile := setupTestLogFile(t)
	WriteHeader(logFile)

	command := "test_command"
	args := []string{"arg1", "arg2"}
	pid := 12345
	jobID := "test-job-id-12345"
	name := "test-job"

	LogWithFile(command, args, pid, jobID, name, logFile)

	// Verify the entry exists
	content, err := ioutil.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), jobID) {
		t.Fatal("Test entry should exist before deletion")
	}

	DeleteLogWithFile(jobID, logFile)

	// Verify the entry is deleted
	content, err = ioutil.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if strings.Contains(string(content), jobID) {
		t.Error("Test entry should be deleted")
	}
}

func TestLogWithPasswordMasking(t *testing.T) {
	logFile := setupTestLogFile(t)
	WriteHeader(logFile)

	command := "ssh"
	args := []string{"-password", "secret123", "user@host"}
	pid := 12345
	jobID := "test-job-id-12345"
	name := "test-job"

	LogWithFile(command, args, pid, jobID, name, logFile)

	content, err := ioutil.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if strings.Contains(string(content), "secret123") {
		t.Error("Password should be masked in log")
	}

	if !strings.Contains(string(content), "******") {
		t.Error("Masked password should appear as ******")
	}
}
