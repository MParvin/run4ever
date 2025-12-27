# run4ever
`run4ever` is a command-line tool that allows you to run a specified command indefinitely, with a specified delay between each execution.

## Run4ever vs Bash while true

Bash while true vs run4ever — at this maturity

### Bash while true

```bash
while true; do
    echo "Hello, world!"
    sleep 10
done
```

### Run4ever

But your run4ever already provides:

[X] Persistent execution

[X] Delay management

[X] Retry counter

[X] State file with PIDs    

[X] --ps process visibility

[X] Signal handling

[X] Notifications (Telegram/Desktop)

[X] Password masking

[X] Structured CLI

[X] Container support

This is no longer “just a loop”.
It is a lightweight self-healing job runner.


## Build
```bash
git clone github.com/mparvin/run4ever
go build -o run4ever
```

## Usage
```bash
run4ever [flags] [command]
```

## Flags
```bash
    --ps : Show a list of running commands and their PIDs.
    -d or --delay: Specify the delay in seconds between command executions. Default is 10 seconds.
    -t or --timeout: Specify the timeout in seconds for command execution. Default is no timeout.
    -v or --verbose: Enable verbose mode. This will cause run4ever to print additional output such as errors and confirmation messages.
    -m or --max-retries: Maximum number of retries before giving up. -1 for infinite retries (default is -1).
    -g or --background: Run command in background (daemon mode).
    --notify-on: Notify on: failure, success, always.
    --notify-method: Notification method: desktop, telegram, slack, email.
    --telegram-token: Telegram bot token (required for Telegram notifications).
    --telegram-chat-id: Telegram chat ID (required for Telegram notifications).
    --telegram-custom-api: Telegram custom API URL (optional).
    --slack-webhook-url: Slack webhook URL (required for Slack notifications).
    --slack-channel: Slack channel (optional, can be set in webhook URL).
    --email-to: Email recipient address (required for email notifications).
    --email-from: Email sender address (required for email notifications).
    --email-password: Email password (required for email notifications).
    --email-smtp: SMTP server hostname (required for email notifications).
    --email-port: SMTP server port (default is 587).
    --exit-on-success: Exit when command succeeds (exit code 0).
    --persist: Save job definition for restore on restart.
    --restore: Restore and run all saved jobs.
```

## Examples

### Basic usage
```bash
run4ever -d 30 echo hello world
```

### Docker pull with retry until success
```bash
run4ever -d 1 --exit-on-success docker pull nginx
```
This will retry every second until the image is successfully pulled, then exit.

### Persist and restore jobs
```bash
# Save a job for later restoration
run4ever -d 60 --persist my-command

# Restore all saved jobs (useful after container restart)
run4ever --restore
```

All examples are in [examples](examples) directory.

## Description
`run4ever` is a command-line tool that allows you to run a specified command indefinitely, with a specified delay between each execution.

## Limitations

### ⚠️ Not recommended for backup scheduling
run4ever is not a scheduler or cron replacement. For critical backups, prefer cron or systemd timers. run4ever can be used only as a retry/notification wrapper.


#### Todo
- [X] Fix flags conflict between run4ever and command (Fixed)
- [X] Add ps flag
- [X] Move runInfinitely function to tools package
- [X] Add github actions
- [X] Add tests
- [X] Add timeout flag
- [X] Add `-g` flag to run command in background
- [X] Add `-l` flag to list all running jobs
- [X] Add more examples - added in examples directory
- [X] Add bash completion
- [X] Add notification on Success/Failure
- [X] Add notification Desktop method
- [X] Add notification Telegram method
- [X] Add notification Slack method
- [X] Add notification Email method
