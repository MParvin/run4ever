package tools

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func Log(command string, pid int) {
	HomeDir := os.Getenv("HOME")
	logFile := HomeDir + "/.run4ever/run4ever.state"

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	t := time.Now()
	tf := t.Format("2006-01-02 15:04:05")
	if _, err := f.WriteString(fmt.Sprintf("%s %d %s\n", tf, pid, command)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func DeleteLog(pid int) {
	HomeDir := os.Getenv("HOME")
	logFile := HomeDir + "/.run4ever/run4ever.state"
	tempFileDir := HomeDir + "/.run4ever/run4ever.state.temp"

	f, err := os.Open(logFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	tempFile, err := os.OpenFile(tempFileDir, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer tempFile.Close()

	scanner := bufio.NewScanner(f)
	writer := bufio.NewWriter(tempFile)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, strconv.Itoa(pid)) {
			continue
		}
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	if err := writer.Flush(); err != nil {
		log.Fatal(err)
	}
	if err := os.Rename(tempFile.Name(), logFile); err != nil {
		log.Fatal(err)
	}
}