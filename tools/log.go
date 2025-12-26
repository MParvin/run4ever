package tools

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// JobState represents a single job entry in the state file
type JobState struct {
	JobID     string
	Name      string
	PID       int
	Command   string
	Args      string
	StartTime time.Time
	IsStale   bool
}

var (
	stateMutex sync.Mutex
)

// GenerateJobID generates a unique job ID (UUID-like)
func GenerateJobID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// GetStateFile returns the path to the state file
func GetStateFile() string {
	homeDir := os.Getenv("HOME")
	return filepath.Join(homeDir, ".run4ever", "run4ever.state")
}

// WriteHeader writes the header line to the state file
func WriteHeader(LogFile string) {
	stateMutex.Lock()
	defer stateMutex.Unlock()

	// Check if file exists and has content
	if info, err := os.Stat(LogFile); err == nil && info.Size() > 0 {
		return // Header already exists
	}

	header := "Time \t\t\t | Job-ID \t\t | PID \t\t | Command \t | Args \t\t | Status\n"
	if err := atomicWriteFile(LogFile, []byte(header), 0644); err != nil {
		log.Fatal(err)
	}
}

// Log adds a new job entry to the state file
func Log(command string, args []string, pid int) {
	jobID, err := GenerateJobID()
	if err != nil {
		log.Fatalf("Failed to generate job ID: %v", err)
	}
	LogWithJobID(command, args, pid, jobID, "")
}

// LogWithJobID adds a new job entry with a specific job ID
func LogWithJobID(command string, args []string, pid int, jobID string, name string) {
	LogFile := GetStateFile()
	LogWithFile(command, args, pid, jobID, name, LogFile)
}

// LogWithFile adds a new job entry to a specific state file
func LogWithFile(command string, args []string, pid int, jobID string, name string, logFile string) {
	stateMutex.Lock()
	defer stateMutex.Unlock()

	// Read existing state
	jobs, err := readStateFile(logFile)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	// Add new job
	maskedArgs := MaskPassword(args)
	t := time.Now()

	newJob := JobState{
		JobID:     jobID,
		Name:      name,
		PID:       pid,
		Command:   command,
		Args:      strings.Join(maskedArgs, " "),
		StartTime: t,
		IsStale:   false,
	}
	jobs = append(jobs, newJob)

	// Write state atomically
	if err := writeStateFile(logFile, jobs); err != nil {
		log.Fatal(err)
	}
}

// DeleteLog removes a job entry by job ID
func DeleteLog(jobID string) {
	LogFile := GetStateFile()
	DeleteLogWithFile(jobID, LogFile)
}

// DeleteLogWithFile removes a job entry from a specific state file
func DeleteLogWithFile(jobID string, logFile string) {
	stateMutex.Lock()
	defer stateMutex.Unlock()

	jobs, err := readStateFile(logFile)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Fatal(err)
	}

	// Remove job with matching job ID
	var filteredJobs []JobState
	for _, job := range jobs {
		if job.JobID != jobID {
			filteredJobs = append(filteredJobs, job)
		}
	}

	// Write state atomically
	if err := writeStateFile(logFile, filteredJobs); err != nil {
		log.Fatal(err)
	}
}

// DeleteLogByPID removes a job entry by PID (for backward compatibility)
func DeleteLogByPID(pid int) {
	LogFile := GetStateFile()
	stateMutex.Lock()
	defer stateMutex.Unlock()

	jobs, err := readStateFile(LogFile)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Fatal(err)
	}

	// Remove job with matching PID
	var filteredJobs []JobState
	for _, job := range jobs {
		if job.PID != pid {
			filteredJobs = append(filteredJobs, job)
		}
	}

	// Write state atomically
	if err := writeStateFile(LogFile, filteredJobs); err != nil {
		log.Fatal(err)
	}
}

// readStateFile reads the state file and returns all job entries
func readStateFile(logFile string) ([]JobState, error) {
	var jobs []JobState

	f, err := os.Open(logFile)
	if err != nil {
		return jobs, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if lineNum == 1 || line == "" {
			continue // Skip header and empty lines
		}

		// Parse line: Time | Job-ID | PID | Command | Args | Status
		parts := strings.Split(line, "|")
		if len(parts) < 6 {
			continue // Skip malformed lines
		}

		timeStr := strings.TrimSpace(parts[0])
		jobID := strings.TrimSpace(parts[1])
		pidStr := strings.TrimSpace(parts[2])
		command := strings.TrimSpace(parts[3])
		args := strings.TrimSpace(parts[4])
		status := strings.TrimSpace(parts[5])

		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue
		}

		startTime, err := time.Parse("2006-01-02 15:04:05", timeStr)
		if err != nil {
			startTime = time.Now() // Fallback to current time
		}

		isStale := status == "STALE"

		// Check if process is actually running
		if !isStale && !isProcessRunning(pid) {
			isStale = true
		}

		jobs = append(jobs, JobState{
			JobID:     jobID,
			PID:       pid,
			Command:   command,
			Args:      args,
			StartTime: startTime,
			IsStale:   isStale,
		})
	}

	return jobs, scanner.Err()
}

// writeStateFile writes all job entries to the state file atomically
func writeStateFile(logFile string, jobs []JobState) error {
	var lines []string
	lines = append(lines, "Time \t\t\t | Job-ID \t\t | PID \t\t | Command \t | Args \t\t | Status\n")

	for _, job := range jobs {
		tf := job.StartTime.Format("2006-01-02 15:04:05")
		status := "RUNNING"
		if job.IsStale {
			status = "STALE"
		}
		line := fmt.Sprintf("%s \t | %s \t | %d \t | %s \t\t | %s \t\t | %s\n",
			tf, job.JobID, job.PID, job.Command, job.Args, status)
		lines = append(lines, line)
	}

	content := []byte(strings.Join(lines, ""))
	return atomicWriteFile(logFile, content, 0644)
}

// atomicWriteFile writes content to a file atomically using a temp file and rename
func atomicWriteFile(filename string, content []byte, perm os.FileMode) error {
	dir := filepath.Dir(filename)
	tmpFile, err := os.CreateTemp(dir, filepath.Base(filename)+".tmp.*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpName := tmpFile.Name()

	// Write content to temp file
	if _, err := tmpFile.Write(content); err != nil {
		tmpFile.Close()
		os.Remove(tmpName)
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// Set permissions
	if err := tmpFile.Chmod(perm); err != nil {
		tmpFile.Close()
		os.Remove(tmpName)
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	// Close temp file before rename
	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tmpName, filename); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// isProcessRunning checks if a process with the given PID is running
func isProcessRunning(pid int) bool {
	// Use ps command to check if process exists
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid))
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// Ps displays all running jobs in a continuous loop
func Ps() {
	LogFile := GetStateFile()

	for {
		stateMutex.Lock()
		jobs, err := readStateFile(LogFile)
		stateMutex.Unlock()

		if err != nil && !os.IsNotExist(err) {
			log.Fatal(err)
		}

		// Print header
		fmt.Println("Time \t\t\t | Job-ID \t\t | PID \t\t | Command \t | Args \t\t | Status")
		fmt.Println(strings.Repeat("-", 120))

		// Print jobs
		for _, job := range jobs {
			status := "RUNNING"
			if job.IsStale {
				status = "STALE"
			}
			tf := job.StartTime.Format("2006-01-02 15:04:05")
			fmt.Printf("%s \t | %s \t | %d \t | %s \t\t | %s \t\t | %s\n",
				tf, job.JobID, job.PID, job.Command, job.Args, status)
		}

		time.Sleep(3 * time.Second)
		fmt.Print("\033[H\033[2J")
	}
}

// ListJobs displays all running jobs once and exits
func ListJobs() {
	LogFile := GetStateFile()

	stateMutex.Lock()
	jobs, err := readStateFile(LogFile)
	stateMutex.Unlock()

	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	if len(jobs) == 0 {
		fmt.Println("No running jobs found.")
		return
	}

	// Print header
	fmt.Println("Time \t\t\t | Job-ID \t\t | PID \t\t | Command \t | Args \t\t | Status")
	fmt.Println(strings.Repeat("-", 120))

	// Print jobs
	for _, job := range jobs {
		status := "RUNNING"
		if job.IsStale {
			status = "STALE"
		}
		tf := job.StartTime.Format("2006-01-02 15:04:05")
		fmt.Printf("%s \t | %s \t | %d \t | %s \t\t | %s \t\t | %s\n",
			tf, job.JobID, job.PID, job.Command, job.Args, status)
	}
}
