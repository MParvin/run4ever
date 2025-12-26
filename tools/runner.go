package tools

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

var (
	telegramToken     string
	telegramChatID    string
	telegramCustomAPI string
	slackWebhookURL   string
	emailTo           string
	emailFrom         string
	emailPassword     string
	emailSMTPHost     string
	emailSMTPPort     int
)

func RunInfinitely(delayInt int, timeoutInt int, args []string, verbose bool, maxRetries int, notifyOn string, notifyMethod string, token string, chatID string, customAPI string, exitOnSuccess bool, slackWebhook string, emailToAddr string, emailFromAddr string, emailPass string, emailSMTP string, emailPort int) {
	telegramToken = token
	telegramChatID = chatID
	telegramCustomAPI = customAPI
	slackWebhookURL = slackWebhook
	emailTo = emailToAddr
	emailFrom = emailFromAddr
	emailPassword = emailPass
	emailSMTPHost = emailSMTP
	emailSMTPPort = emailPort

	retryCount := 0
	for {
		exitStatus := 0

		if maxRetries != -1 && retryCount >= maxRetries {
			if verbose {
				fmt.Println("Max retries reached, exiting")
			}
			os.Exit(1)
		}

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		// Set up timeout if specified
		var err error
		if timeoutInt > 0 {
			if verbose {
				fmt.Printf("Running command with timeout: %d seconds\n", timeoutInt)
			}
			err = runWithTimeout(cmd, timeoutInt)
		} else {
			err = cmd.Run()
		}

		if err != nil {
			if verbose {
				fmt.Println(err)
			}
			if cmd.ProcessState != nil {
				exitStatus = cmd.ProcessState.ExitCode()
			} else {
				// Command was killed due to timeout
				exitStatus = 124 // Standard timeout exit code
			}
			retryCount++
		}

		// Handle exit-on-success: if command succeeded, exit
		if exitOnSuccess && exitStatus == 0 {
			if shouldNotify(notifyOn, exitStatus) {
				title := "run4ever: Task " + statusToString(exitStatus)
				maskedArgs := MaskPassword(args)
				message := fmt.Sprintf("Command %s %s exited with status %d", args[0], maskedArgs, exitStatus)
				if verbose {
					fmt.Printf("Sending notification\nTitle: %s\nMessage: %s\n", title, message)
				}
				doNotify(notifyOn, notifyMethod, verbose, title, message)
			}
			if verbose {
				fmt.Printf("Command `%s` succeeded, exiting as requested\n", args[0])
			}
			os.Exit(0)
		}

		if shouldNotify(notifyOn, exitStatus) {
			title := "run4ever: Task " + statusToString(exitStatus)
			maskedArgs := MaskPassword(args)
			message := fmt.Sprintf("Command %s %s exited with status %d", args[0], maskedArgs, exitStatus)
			if verbose {
				fmt.Printf("Sending notification\nTitle: %s\nMessage: %s\n", title, message)
			}
			doNotify(notifyOn, notifyMethod, verbose, title, message)
		}
		if verbose {
			fmt.Printf("Command `%s` exited with status %d\n", args[0], exitStatus)
			fmt.Printf("Sleeping for %d seconds\n", delayInt)
		}
		time.Sleep(time.Duration(delayInt) * time.Second)
	}
}

// runWithTimeout runs a command with a timeout
func runWithTimeout(cmd *exec.Cmd, timeoutSeconds int) error {
	// Start the command
	err := cmd.Start()
	if err != nil {
		return err
	}

	// Create a channel to signal when the command completes
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	// Wait for either completion or timeout
	select {
	case err := <-done:
		return err
	case <-time.After(time.Duration(timeoutSeconds) * time.Second):
		// Timeout occurred, kill the process
		if err := cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill process after timeout: %w", err)
		}
		return fmt.Errorf("command timed out after %d seconds", timeoutSeconds)
	}
}

func doNotify(notifyOn string, notifyMethod string, verbose bool, title string, message string) {
	switch notifyMethod {
	case "desktop":
		if err := SendDesktopNotification(title, message, verbose); err != nil && verbose {
			fmt.Println("Error sending desktop notification: ", err)
		}
	case "telegram":
		if err := SendTelegramNotification(telegramToken, telegramChatID, message, telegramCustomAPI, verbose); err != nil && verbose {
			fmt.Println("Error sending Telegram notification: ", err)
		}
	case "slack":
		if err := SendSlackNotification(slackWebhookURL, message, verbose); err != nil && verbose {
			fmt.Println("Error sending Slack notification: ", err)
		}
	case "email":
		if err := SendEmailNotification(emailTo, emailFrom, emailPassword, emailSMTPHost, emailSMTPPort, title, message, verbose); err != nil && verbose {
			fmt.Println("Error sending email notification: ", err)
		}
	}
}

func shouldNotify(notifyOn string, exitStatus int) bool {
	switch notifyOn {
	case "always":
		return true
	case "success":
		return exitStatus == 0
	case "failure":
		return exitStatus != 0
	default:
		return false
	}
}

func statusToString(exitStatus int) string {
	if exitStatus == 0 {
		return "Success"
	}
	return "Failure"
}
