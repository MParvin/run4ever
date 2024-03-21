/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mparvin/run4ever/cmd"
	tools "github.com/mparvin/run4ever/tools"
)

func main() {
	HomeDir := os.Getenv("HOME")
	tools.CreateDir(HomeDir + "/.run4ever")
	LogFile := HomeDir + "/.run4ever/run4ever.state"

	if _, err := os.Stat(LogFile); os.IsNotExist(err) || tools.IsEmpty(LogFile) {
		tools.WriteHeader(LogFile)
	}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		tools.DeleteLog(os.Getpid())
		os.Exit(1)
	}()
	cmd.Execute()
	select {}

}
