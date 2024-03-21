package tools

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func WriteHeader(LogFile string) {
	f, err := os.OpenFile(LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if _, err := f.WriteString("Time \t\t\t | PID \t\t | Command \t | Args \t\t\n"); err != nil {
		log.Fatal(err)
	}
}

func Log(command string, args []string, pid int) {
	HomeDir := os.Getenv("HOME")
	LogFile := HomeDir + "/.run4ever/run4ever.state"

	f, err := os.OpenFile(LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	t := time.Now()
	tf := t.Format("2006-01-02 15:04:05")
	if _, err := f.WriteString(fmt.Sprintf("%s \t | %d \t | %s \t\t | %s\n", tf, pid, command, strings.Join(args, " "))); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func DeleteLog(pid int) {
	HomeDir := os.Getenv("HOME")
	// LogFile := HomeDir + "/.run4ever/run4ever.state"
	tempFileDir := HomeDir + "/.run4ever/run4ever.state.temp"

	f, err := os.Open(LogFile)
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
	if err := os.Rename(tempFile.Name(), LogFile); err != nil {
		log.Fatal(err)
	}
}

func Watch() {
	HomeDir := os.Getenv("HOME")
	LogFile := HomeDir + "/.run4ever/run4ever.state"

	for {
		f, err := os.Open(LogFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		_, err = f.Seek(0, io.SeekStart)
		if err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		time.Sleep(3 * time.Second)
		fmt.Print("\033[H\033[2J")
	}

}
