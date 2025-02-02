/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"

	tools "github.com/mparvin/run4ever/tools"
	"github.com/spf13/cobra"
)

var (
	notifyOn  string
	notifyMethod string
	telegramToken string
	telegramChatID string
	telegramCustomAPI string
)

var delay string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "run4ever [flags] [command] [arguments]",
	Short:   "Run a command indefinitely with a specified delay between executions.",
	Example: "run4ever -d 30 echo hello world",
	Long: `run4ever is a command-line tool that allows you to run a specified command repeatedly, with a specified delay between each execution.

The -d flag or --delay flag is used to specify the delay between each execution of the command. The default value is 10 seconds.

Use the --ps to show a list of running commands and their PIDs, It will not run the provided command.

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
	Run: func(cmd *cobra.Command, args []string) {},
}

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
	rootCmd.Flags().BoolP("ps", "", false, "Running PIDs")
	rootCmd.Flags().StringVar(&notifyOn, "notify-on", "", "Notify on: failure, success, always")
	rootCmd.Flags().StringVar(&notifyMethod, "notify-method", "desktop", "Notification method: desktop, telegram")
	rootCmd.Flags().StringVar(&telegramToken, "telegram-token", "", "Telegram bot token (required for Telegram notifications)")
	rootCmd.Flags().StringVar(&telegramChatID, "telegram-chat-id", "", "Telegram chat ID (required for Telegram notifications)")
	rootCmd.Flags().StringVar(&telegramCustomAPI, "telegram-custom-api", "", "Telegram custom API URL (optional)")


	rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && rootCmd.Flags().Lookup("ps").Value.String() == "false" {
			log.Fatal("No command provided")
		}
	}

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		delayInt, err := strconv.Atoi(delay)
		if err != nil {
			log.Fatal("Invalid delay value provided")
		}

		verbose := false
		verbose, _ = cmd.Flags().GetBool("verbose")

		if verbose {
			fmt.Println("run4ever called")
			fmt.Println("dely is", delayInt)
		}

		psProvided, _ := cmd.Flags().GetBool("ps")
		if psProvided {
			tools.Ps()
			return
		}
		tools.RunInfinitely(
			delayInt, 
			args, 
			verbose,
			notifyOn,
			notifyMethod,
			telegramToken,
			telegramChatID,
			telegramCustomAPI,
		)
	}
}
