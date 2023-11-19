package tools

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func RunInfinitely(delayInt int, args []string, verbose bool) {
	for {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err := cmd.Run()
		if err != nil {
			if verbose {
				fmt.Println(err)
			}
			time.Sleep(time.Duration(delayInt) * time.Second)
			continue
		}
		if verbose {
			fmt.Printf("Command %s exited", args[0])
			fmt.Print("Sleeping for ", delayInt, " seconds")
		}
		time.Sleep(time.Duration(delayInt) * time.Second)
	}

}
