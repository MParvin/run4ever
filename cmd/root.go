/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"
	"strconv"
	"log"

	tools "github.com/mparvin/run4ever/tools"
	"github.com/spf13/cobra"
)

var delay string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "run4ever [flags] [command] [arguments]",
	Short:   "Run a command indefinitely with a specified delay between executions.",
	Example: "run4ever -d 30 echo hello world",
	Long: `run4ever is a command-line tool that allows you to run a specified command repeatedly, with a specified delay between each execution.

The -d flag or --delay flag is used to specify the delay between each execution of the command. The default value is 10 seconds.

Use the -w flag or --watch to show a list of running commands and their PIDs, It will not run the provided command.

You can also enable verbose mode by using the -v flag. This will cause run4ever to print additional output such as errors and confirmation messages.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			tools.Log(args[0], args[1:], os.Getpid())
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		tools.DeleteLog(os.Getpid())
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		tools.DeleteLog(os.Getpid())
		return nil
	},
	DisableFlagParsing: false,
	// BashCompletionFunction:     bashCompletionFunc,
	Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&delay, "delay", "d", "10", "delay between to run of the command")
	rootCmd.Flags().BoolP("verbose", "v", false, "Verbose mode")
	rootCmd.Flags().SetInterspersed(false)
	rootCmd.Flags().BoolP("watch", "w", false, "Watch mode")

	rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if len(args) == 0  && rootCmd.Flags().Lookup("watch").Value.String() == "false" {
			log.Fatal("No command provided")
		}
	}

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		delayInt, err := strconv.Atoi(delay)
		if err != nil {
			log.Fatal("Invalid delay value provided")
		}

		verbose := false
		verbose, err = cmd.Flags().GetBool("verbose")
		if err != nil {
			log.Fatal("Error getting verbose flag")
		}

		if verbose {
			fmt.Println("run4ever called")
			fmt.Println("dely is", delayInt)
		}

		watch, err := cmd.Flags().GetBool("watch")
		if err != nil {
			log.Fatal("Error getting watch flag")
		}
		if watch {
			tools.Watch()
			return
		}
		runInfinitely(delayInt, args, verbose)
	}

}

func runInfinitely(delayInt int, args []string, verbose bool) {
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
