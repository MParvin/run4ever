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
)

func RunInfinitely(delayInt int, args []string, verbose bool, maxRetries int, notifyOn string, notifyMethod string, token string, chatID string, customAPI string) {
	telegramToken = token
	telegramChatID = chatID
	telegramCustomAPI = customAPI

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
		err := cmd.Run()

		if err != nil {
			if verbose {
				fmt.Println(err)
			}
			exitStatus = cmd.ProcessState.ExitCode()
			retryCount++
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
