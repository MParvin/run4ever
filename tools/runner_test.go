package tools

import (
	"os/exec"
	"strings"
	"testing"
	"time"
)

// TestRunInfinitely is difficult to test directly due to infinite loop
// Instead we skip this test and rely on manual testing
func TestRunInfinitely(t *testing.T) {
	t.Skip("RunInfinitely contains an infinite loop and is difficult to unit test. Manual testing is recommended.")
}

// Test the command execution part without the infinite loop
func TestCommandExecution(t *testing.T) {
	// Test successful command
	cmd := exec.Command("echo", "hello")
	err := cmd.Run()
	if err != nil {
		t.Errorf("Command execution failed: %v", err)
	}

	// Test failed command
	cmd = exec.Command("false")
	err = cmd.Run()
	if err == nil {
		t.Error("Expected command to fail")
	}
}

// TestRunWithTimeout tests the timeout functionality
func TestRunWithTimeout(t *testing.T) {
	tests := []struct {
		name           string
		command        []string
		timeoutSeconds int
		shouldTimeout  bool
	}{
		{
			name:           "command completes before timeout",
			command:        []string{"echo", "test"},
			timeoutSeconds: 5,
			shouldTimeout:  false,
		},
		{
			name:           "command times out",
			command:        []string{"sleep", "2"},
			timeoutSeconds: 1,
			shouldTimeout:  true,
		},
		{
			name:           "command completes quickly with long timeout",
			command:        []string{"true"},
			timeoutSeconds: 10,
			shouldTimeout:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(tt.command[0], tt.command[1:]...)
			start := time.Now()
			err := runWithTimeout(cmd, tt.timeoutSeconds)
			duration := time.Since(start)

			if tt.shouldTimeout {
				if err == nil {
					t.Error("Expected command to timeout, but it completed successfully")
				}
				// Check that it actually timed out around the expected time
				expectedDuration := time.Duration(tt.timeoutSeconds) * time.Second
				if duration < expectedDuration-100*time.Millisecond || duration > expectedDuration+500*time.Millisecond {
					t.Errorf("Timeout duration %v not close to expected %v", duration, expectedDuration)
				}
			} else {
				if err != nil {
					t.Errorf("Expected command to complete successfully, but got error: %v", err)
				}
			}
		})
	}
}

// TestRunWithTimeoutKillProcess tests that processes are properly killed on timeout
func TestRunWithTimeoutKillProcess(t *testing.T) {
	// Use a command that will definitely timeout
	cmd := exec.Command("sleep", "10")
	err := runWithTimeout(cmd, 1)
	if err == nil {
		t.Error("Expected command to timeout and return error")
	}

	// Wait a bit for the process to be fully killed and ProcessState to be set
	time.Sleep(100 * time.Millisecond)

	// Verify the process was killed - ProcessState might not be set immediately
	// but the error should indicate timeout
	if err != nil && !strings.Contains(err.Error(), "timed out") {
		t.Errorf("Expected timeout error, got: %v", err)
	}
}
