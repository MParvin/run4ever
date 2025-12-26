package tools

import (
	"strings"
)

var passwordFlags = []string{
	"-password", "--password", "-p", "--pass",
	"--secret", "--key", "--token",
	"--db-pass", "--api-key", "--api-token",
	"--auth", "--auth-token", "--auth-key",
	"--telegram-token", "--slack-token",
	"--email-pass", "--email-token",
}

func MaskPassword(args []string) []string {
	var maskedArgs []string
	skipNext := false

	for i, arg := range args {
		if skipNext {
			skipNext = false
			continue
		}
	
		isPasswordFlag := false
		for _, flag := range passwordFlags {
			if strings.HasPrefix(arg, flag+"=") || arg == flag {
				isPasswordFlag = true
				break
			}
		}

		if isPasswordFlag {
			if strings.Contains(arg, "=") {
				parts := strings.SplitN(arg, "=", 2)
				maskedArgs = append(maskedArgs, parts[0]+"=******")
			} else {
				maskedArgs = append(maskedArgs, arg)
				if i+1 < len(args) {
					maskedArgs = append(maskedArgs, "******")
					skipNext = true
				}
			}
		} else {
			maskedArgs = append(maskedArgs, arg)
		}
	}

	return maskedArgs
}