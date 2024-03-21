package tools

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestRunInfinitely(t *testing.T) {
	f, err := os.Create("/tmp/test.sh")
	if err != nil {
		fmt.Println(err)
	}
	f.WriteString("#!/bin/bash\necho \"Hello World\"")
	f.Close()

	err = os.Chmod("/tmp/test.sh", 0755)
	if err != nil {
		fmt.Println(err)
	}

	go RunInfinitely(1, []string{"/tmp/test.sh"}, false)
	time.Sleep(5 * time.Second)

	cmd := exec.Command("kill", "-9", fmt.Sprintf("%d", os.Getpid()))
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	os.Remove("/tmp/test.sh")
}
