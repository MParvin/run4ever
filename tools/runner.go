package tools

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func RunInfinitely(delayInt int, args []string, verbose bool, notifyOn string, notifyMethod string, telegramToken string, telegramChatID string, telegramCustomAPI string) {
	for {
		exitStatus := 0

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
		}

		if ShouldNotify(notifyOn, exitStatus) {
			title := "run4ever: Task " + StatusToString(exitStatus)
			maskedArgs := MaskPassword(args)
			message := fmt.Sprintf("Command %s %s exited with status %d", args[0], maskedArgs, exitStatus)

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
		if verbose {
			fmt.Printf("Command `%s` exited with status %d\n", args[0], exitStatus)
			fmt.Printf("Sleeping for %d seconds\n", delayInt)
		}
		time.Sleep(time.Duration(delayInt) * time.Second)
	}
}

func ShouldNotify(notifyOn string, exitStatus int) bool {
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

func StatusToString(exitStatus int) string {
	if exitStatus == 0 {
		return "Success"
	}
	return "Failure"
}
