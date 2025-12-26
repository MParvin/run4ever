package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// JobDefinition represents a job that can be persisted and restored
type JobDefinition struct {
	Command          []string `json:"command"`
	Delay            int      `json:"delay"`
	MaxRetries       int      `json:"max_retries"`
	Timeout          int      `json:"timeout"`
	NotifyOn         string   `json:"notify_on"`
	NotifyMethod     string   `json:"notify_method"`
	TelegramToken    string   `json:"telegram_token,omitempty"`
	TelegramChatID   string   `json:"telegram_chat_id,omitempty"`
	TelegramCustomAPI string  `json:"telegram_custom_api,omitempty"`
	ExitOnSuccess    bool     `json:"exit_on_success"`
}

// GetJobsFile returns the path to the jobs persistence file
func GetJobsFile() string {
	homeDir := os.Getenv("HOME")
	return filepath.Join(homeDir, ".run4ever", "jobs.json")
}

// SaveJobDefinition saves a job definition to the jobs file
func SaveJobDefinition(job JobDefinition) error {
	jobsFile := GetJobsFile()
	dir := filepath.Dir(jobsFile)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create jobs directory: %w", err)
	}

	// Read existing jobs
	jobs, err := loadJobDefinitions(jobsFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to load existing jobs: %w", err)
	}

	// Append new job
	jobs = append(jobs, job)

	// Write jobs atomically
	return saveJobDefinitions(jobsFile, jobs)
}

// loadJobDefinitions loads all job definitions from the jobs file
func loadJobDefinitions(jobsFile string) ([]JobDefinition, error) {
	var jobs []JobDefinition

	data, err := os.ReadFile(jobsFile)
	if err != nil {
		return jobs, err
	}

	if len(data) == 0 {
		return jobs, nil
	}

	if err := json.Unmarshal(data, &jobs); err != nil {
		return jobs, fmt.Errorf("failed to parse jobs file: %w", err)
	}

	return jobs, nil
}

// saveJobDefinitions saves job definitions to the jobs file atomically
func saveJobDefinitions(jobsFile string, jobs []JobDefinition) error {
	data, err := json.MarshalIndent(jobs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal jobs: %w", err)
	}

	return atomicWriteFile(jobsFile, data, 0600)
}

// RestoreJobs restores and runs all saved jobs
func RestoreJobs(verbose bool) error {
	jobsFile := GetJobsFile()

	jobs, err := loadJobDefinitions(jobsFile)
	if err != nil {
		if os.IsNotExist(err) {
			if verbose {
				fmt.Println("No saved jobs found")
			}
			return nil
		}
		return fmt.Errorf("failed to load jobs: %w", err)
	}

	if len(jobs) == 0 {
		if verbose {
			fmt.Println("No jobs to restore")
		}
		return nil
	}

	if verbose {
		fmt.Printf("Restoring %d job(s)\n", len(jobs))
	}

	// Start each job in the background
	for i, job := range jobs {
		if verbose {
			fmt.Printf("Restoring job %d: %v\n", i+1, job.Command)
		}

		// Build command arguments
		args := []string{
			"-g", // Run in background
			"-d", fmt.Sprintf("%d", job.Delay),
		}

		if job.MaxRetries != -1 {
			args = append(args, "-m", fmt.Sprintf("%d", job.MaxRetries))
		}

		if job.Timeout > 0 {
			args = append(args, "-t", fmt.Sprintf("%d", job.Timeout))
		}

		if job.NotifyOn != "" {
			args = append(args, "--notify-on", job.NotifyOn)
		}

		if job.NotifyMethod != "" {
			args = append(args, "--notify-method", job.NotifyMethod)
		}

		if job.TelegramToken != "" {
			args = append(args, "--telegram-token", job.TelegramToken)
		}

		if job.TelegramChatID != "" {
			args = append(args, "--telegram-chat-id", job.TelegramChatID)
		}

		if job.TelegramCustomAPI != "" {
			args = append(args, "--telegram-custom-api", job.TelegramCustomAPI)
		}

		if job.ExitOnSuccess {
			args = append(args, "--exit-on-success")
		}

		// Add command
		args = append(args, job.Command...)

		// Start run4ever in background for this job
		cmd := exec.Command(os.Args[0], args...)
		cmd.Stdout = nil
		cmd.Stderr = nil
		cmd.Stdin = nil

		if err := cmd.Start(); err != nil {
			if verbose {
				fmt.Printf("Warning: failed to start job %d: %v\n", i+1, err)
			}
			continue
		}

		if verbose {
			fmt.Printf("Started job %d with PID: %d\n", i+1, cmd.Process.Pid)
		}
	}

	return nil
}

