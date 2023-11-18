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

You can use the -d flag to specify the delay in seconds between command executions. By default, the delay is 10 seconds.

You can also enable verbose mode by using the -v flag. This will cause run4ever to print additional output such as errors and confirmation messages.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			tools.Log(args[0], os.Getpid())
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
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&delay, "delay", "d", "10", "delay between to run of the command")
	rootCmd.Flags().BoolP("verbose", "v", false, "Verbose mode")
	rootCmd.Flags().SetInterspersed(false)

	rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("No command provided")
			os.Exit(1)
		}
	}

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		delayInt, err := strconv.Atoi(delay)
		if err != nil {
			fmt.Println("Invalid delay value")
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
			fmt.Println("dely is", delayInt)
		}
		runInfinitely(delayInt, args, verbose)
	}

}

func runInfinitely(delayInt int, args []string, verbose bool) {
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
			// Sleep for delayInt seconds
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
