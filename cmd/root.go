/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var timeout string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:                "	run4ever [flags] [your command]",
	Short:              "run4ever is a CLI tool to run a command forever",
	Example:            "	run4ever -t 10 echo hello world",
	Long:               `run4ever is a CLI tool to run a command forever.`,
	DisableFlagParsing: false,
	// BashCompletionFunction:     bashCompletionFunc,
	Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&timeout, "timeout", "t", "10", "timeout in seconds")
	rootCmd.Flags().BoolP("verbose", "v", false, "Verbose mode")
	rootCmd.Flags().SetInterspersed(false)

	rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("No command provided")
			os.Exit(1)
		}
	}

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		timeoutInt, err := strconv.Atoi(timeout)
		if err != nil {
			fmt.Println("Invalid timeout value")
			os.Exit(1)
		}

		verbose := false
		//  convert cmd.Flags().Lookup("verbose").Value  to bool and store in verbose
		verbose, err = cmd.Flags().GetBool("verbose")
		if err != nil {
			fmt.Println("Error getting verbose flag")
			os.Exit(1)
		}

		if verbose {
			fmt.Println("run4ever called")
			fmt.Println("timeout is", timeoutInt)
		}
		runInfinitely(timeoutInt, args, verbose)
	}

}

func runInfinitely(timeoutInt int, args []string, verbose bool) {
	for {
		// cmd is all the arguments and flags (Instread of this cmd flags) passed to run4ever
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err := cmd.Run()
		if err != nil {
			// if -v flag is set, print error
			if verbose {
				fmt.Println(err)
			}
			// Sleep for timeout seconds before trying again
			time.Sleep(time.Duration(timeoutInt) * time.Second)
			continue
		}
		if verbose {
			fmt.Printf("Command %s exited", args[0])
			fmt.Print("Sleeping for ", timeoutInt, " seconds")
		}
		time.Sleep(time.Duration(timeoutInt) * time.Second)
	}

}
