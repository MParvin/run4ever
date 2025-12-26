/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	tools "github.com/mparvin/run4ever/tools"
	"github.com/spf13/cobra"
)

var (
	notifyOn          string
	notifyMethod      string
	telegramToken     string
	telegramChatID    string
	telegramCustomAPI string
	slackWebhookURL   string
	slackChannel      string
	emailTo           string
	emailFrom         string
	emailPassword     string
	emailSMTPHost     string
	emailSMTPPort     int
	delay             string
	maxRetries        int
	timeout           string
	background        bool
	exitOnSuccess     bool
	persist           bool
	restore           bool
	currentJobID      string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "run4ever [flags] [command] [arguments]",
	Short:   "Run a command indefinitely with a specified delay between executions.",
	Example: "run4ever -d 30 echo hello world",
	Long: `run4ever is a command-line tool that allows you to run a specified command repeatedly, with a specified delay between each execution.

The -d flag or --delay flag is used to specify the delay between each execution of the command. The default value is 10 seconds.

Use the --ps to show a list of running commands and their PIDs continuously, or -l to list once and exit.

You can also enable verbose mode by using the -v flag. This will cause run4ever to print additional output such as errors and confirmation messages.

Notification methods supported: desktop, telegram, slack, email.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			jobID, err := tools.GenerateJobID()
			if err != nil {
				log.Fatalf("Failed to generate job ID: %v", err)
			}
			currentJobID = jobID
			tools.LogWithJobID(args[0], args[1:], os.Getpid(), jobID, "")
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if currentJobID != "" {
			tools.DeleteLog(currentJobID)
		} else {
			// Fallback to PID for backward compatibility
			tools.DeleteLogByPID(os.Getpid())
		}
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if currentJobID != "" {
			tools.DeleteLog(currentJobID)
		} else {
			// Fallback to PID for backward compatibility
			tools.DeleteLogByPID(os.Getpid())
		}
		return nil
	},
	DisableFlagParsing: false,
	Run:                func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

// runInBackground starts the process in the background using nohup
func runInBackground() error {
	// Remove the -g flag from arguments for the background process
	args := make([]string, 0, len(os.Args)-1)
	for _, arg := range os.Args[1:] {
		if arg != "-g" && arg != "--background" {
			args = append(args, arg)
		}
	}

	cmd := exec.Command("nohup", append([]string{os.Args[0]}, args...)...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start background process: %w", err)
	}

	fmt.Printf("Started background process with PID: %d\n", cmd.Process.Pid)
	return nil
}

// runAsDaemon starts the process as a proper daemon
func runAsDaemon() error {
	// Remove the -D flag from arguments for the daemon process
	args := make([]string, 0, len(os.Args)-1)
	for _, arg := range os.Args[1:] {
		if arg != "-D" && arg != "--daemon" {
			args = append(args, arg)
		}
	}

	// Start the process with nohup and redirect output to /dev/null
	cmd := exec.Command("nohup", append([]string{os.Args[0]}, args...)...)
	
	// Redirect stdin, stdout, stderr to /dev/null for proper daemon behavior
	devNull, err := os.OpenFile("/dev/null", os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("failed to open /dev/null: %w", err)
	}
	// Don't close devNull - the child process needs these file descriptors
	
	cmd.Stdin = devNull
	cmd.Stdout = devNull
	cmd.Stderr = devNull
	
	// Start the process and detach
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start daemon process: %w", err)
	}

	// Don't wait for the child process
	cmd.Process.Release()
	
	fmt.Printf("Started daemon process with PID: %d\n", cmd.Process.Pid)
	return nil
}

func init() {
	// Enable bash completion
	rootCmd.CompletionOptions.DisableDefaultCmd = false
	rootCmd.AddCommand(&cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:
  $ source <(run4ever completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ run4ever completion bash > /etc/bash_completion.d/run4ever
  # macOS:
  $ run4ever completion bash > /usr/local/etc/bash_completion.d/run4ever

Zsh:
  $ source <(run4ever completion zsh)
  $ run4ever completion zsh > "${fpath[1]}/_run4ever"

Fish:
  $ run4ever completion fish | source
  $ run4ever completion fish > ~/.config/fish/completions/run4ever.fish

PowerShell:
  PS> run4ever completion powershell | Out-String | Invoke-Expression
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				rootCmd.GenBashCompletion(os.Stdout)
			case "zsh":
				rootCmd.GenZshCompletion(os.Stdout)
			case "fish":
				rootCmd.GenFishCompletion(os.Stdout, true)
			case "powershell":
				rootCmd.GenPowerShellCompletion(os.Stdout)
			}
		},
	})

	rootCmd.Flags().StringVarP(&delay, "delay", "d", "10", "delay between to run of the command")
	rootCmd.Flags().BoolP("verbose", "v", false, "Verbose mode")
	rootCmd.Flags().SetInterspersed(false)
	rootCmd.Flags().BoolP("ps", "", false, "Show running jobs continuously (like top)")
	rootCmd.Flags().BoolP("list", "l", false, "List all running jobs once and exit")
	rootCmd.Flags().StringVar(&notifyOn, "notify-on", "", "Notify on: failure, success, always")
	rootCmd.Flags().StringVar(&notifyMethod, "notify-method", "desktop", "Notification method: desktop, telegram, slack, email")
	rootCmd.Flags().StringVar(&telegramToken, "telegram-token", "", "Telegram bot token (required for Telegram notifications)")
	rootCmd.Flags().StringVar(&telegramChatID, "telegram-chat-id", "", "Telegram chat ID (required for Telegram notifications)")
	rootCmd.Flags().StringVar(&telegramCustomAPI, "telegram-custom-api", "", "Telegram custom API URL (optional)")
	rootCmd.Flags().StringVar(&slackWebhookURL, "slack-webhook-url", "", "Slack webhook URL (required for Slack notifications)")
	rootCmd.Flags().StringVar(&slackChannel, "slack-channel", "", "Slack channel (optional, can be set in webhook URL)")
	rootCmd.Flags().StringVar(&emailTo, "email-to", "", "Email recipient address (required for email notifications)")
	rootCmd.Flags().StringVar(&emailFrom, "email-from", "", "Email sender address (required for email notifications)")
	rootCmd.Flags().StringVar(&emailPassword, "email-password", "", "Email password (required for email notifications)")
	rootCmd.Flags().StringVar(&emailSMTPHost, "email-smtp", "", "SMTP server hostname (required for email notifications)")
	rootCmd.Flags().IntVar(&emailSMTPPort, "email-port", 587, "SMTP server port (default is 587)")
	rootCmd.Flags().IntVarP(&maxRetries, "max-retries", "m", -1, "Maximum number of retries before giving up, -1 for infinite retries (default is -1)")
	rootCmd.Flags().StringVarP(&timeout, "timeout", "t", "", "Timeout for command execution in seconds (default is no timeout)")
	rootCmd.Flags().BoolP("background", "g", false, "Run command in background (daemon mode)")
	rootCmd.Flags().BoolP("daemon", "D", false, "Run command as a daemon (detached from terminal)")
	rootCmd.Flags().BoolVar(&exitOnSuccess, "exit-on-success", false, "Exit when command succeeds (exit code 0)")
	rootCmd.Flags().BoolVar(&persist, "persist", false, "Save job definition for restore on restart")
	rootCmd.Flags().BoolVar(&restore, "restore", false, "Restore and run all saved jobs")

	rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		// Handle restore flag
		restoreFlag, _ := cmd.Flags().GetBool("restore")
		if restoreFlag {
			verbose, _ := cmd.Flags().GetBool("verbose")
			if err := tools.RestoreJobs(verbose); err != nil {
				log.Fatalf("Failed to restore jobs: %v", err)
			}
			return
		}

		psProvided, _ := cmd.Flags().GetBool("ps")
		listProvided, _ := cmd.Flags().GetBool("list")
		if len(args) == 0 && !psProvided && !listProvided {
			log.Fatal("No command provided")
		}
	}

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		background, _ := cmd.Flags().GetBool("background")
		daemon, _ := cmd.Flags().GetBool("daemon")

		// Handle daemon execution (takes precedence over background)
		if daemon {
			if err := runAsDaemon(); err != nil {
				log.Fatalf("Failed to run as daemon: %v", err)
			}
			return
		}

		// Handle background execution
		if background {
			if err := runInBackground(); err != nil {
				log.Fatalf("Failed to run in background: %v", err)
			}
			return
		}

		verbose, _ := cmd.Flags().GetBool("verbose")

		// Load config with priority: CLI flags > env vars > config files
		config, err := tools.LoadConfig(verbose)
		if err != nil && verbose {
			fmt.Printf("Warning: failed to load config: %v\n", err)
		}

		// Apply config values, but CLI flags take precedence
		if notifyMethod == "" && config.NotifyMethod != "" {
			notifyMethod = config.NotifyMethod
		}
		if notifyOn == "" && config.NotifyOn != "" {
			notifyOn = config.NotifyOn
		}
		if telegramToken == "" && config.TelegramToken != "" {
			telegramToken = config.TelegramToken
		}
		if telegramChatID == "" && config.TelegramChatID != "" {
			telegramChatID = config.TelegramChatID
		}
		if telegramCustomAPI == "" && config.TelegramCustomAPI != "" {
			telegramCustomAPI = config.TelegramCustomAPI
		}

		delayInt, err := strconv.Atoi(delay)
		if err != nil {
			log.Fatal("Invalid delay value provided")
		}

		timeoutInt := 0
		if timeout != "" {
			timeoutInt, err = strconv.Atoi(timeout)
			if err != nil {
				log.Fatal("Invalid timeout value provided")
			}
		}

		if verbose {
			fmt.Println("run4ever called")
			fmt.Println("delay is", delayInt)
			if timeoutInt > 0 {
				fmt.Println("timeout is", timeoutInt, "seconds")
			}
		}

		psProvided, _ := cmd.Flags().GetBool("ps")
		listProvided, _ := cmd.Flags().GetBool("list")
		if psProvided {
			tools.Ps()
			return
		}
		if listProvided {
			tools.ListJobs()
			return
		}

		// Handle persist flag
		persistFlag, _ := cmd.Flags().GetBool("persist")
		if persistFlag {
			jobDef := tools.JobDefinition{
				Command:           args,
				Delay:             delayInt,
				MaxRetries:        maxRetries,
				Timeout:           timeoutInt,
				NotifyOn:          notifyOn,
				NotifyMethod:      notifyMethod,
				TelegramToken:     telegramToken,
				TelegramChatID:    telegramChatID,
				TelegramCustomAPI: telegramCustomAPI,
				ExitOnSuccess:     exitOnSuccess,
			}
			if err := tools.SaveJobDefinition(jobDef); err != nil {
				log.Fatalf("Failed to save job definition: %v", err)
			}
			if verbose {
				fmt.Println("Job definition saved")
			}
		}

		tools.RunInfinitely(
			delayInt,
			timeoutInt,
			args,
			verbose,
			maxRetries,
			notifyOn,
			notifyMethod,
			telegramToken,
			telegramChatID,
			telegramCustomAPI,
			exitOnSuccess,
			slackWebhookURL,
			emailTo,
			emailFrom,
			emailPassword,
			emailSMTPHost,
			emailSMTPPort,
		)
	}
}
