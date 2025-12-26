# Example of run4ever usage

Run the command `ssh -R 80:localhost:80 user@server.example.com` for ever, and If  the command fails, it will be restarted after 10 seconds.

```bash
./run4ever ssh -R 80:localhost:80 user@server.example.com
```

Setup a SOCKS proxy on port 1080, and If  the command fails, it will be restarted after 5 seconds.
```bash
./run4ever -d 5 ssh -D 1080 user@server.example.com
```

Run the script `my-backup-script.sh` every hour.
```bash
./run4ever -d 3600 ./my-backup-script.sh
```

Monitor the status of a service every minute.
```bash
run4ever -d 60 check-service-status.sh my-service
```

Run `rsync` command every 5 minutes.
```bash
run4ever -d 300 rsync -avz --delete /home/user/ /mnt/backup
```

## Notifications

### Desktop Notification
Show desktop notification on command failure.
```bash
run4ever -d 30 --notify-on failure --notify-method desktop ./my-backup-script.sh
```

### Telegram Notification
Send a telegram message on command failure.
```bash
run4ever -d 60 --notify-on failure --notify-method telegram --telegram-token <token> --telegram-chat-id <chat-id> ./my-backup-script.sh
```

### Slack Notification
Send a slack message on command success.
```bash
run4ever -d 60 --notify-on success --notify-method slack --slack-token <token> --slack-channel <channel> ./my-backup-script.sh
```

### Email Notification
Send an email on command failure.
```bash
run4ever -d 60 --notify-on failure --notify-method email --email-to <to> --email-from <from> --email-password <password> --email-smtp <smtp> --email-port <port> ./my-backup-script.sh
```

## Advanced Examples

### List all running jobs
```bash
run4ever -l
```

### Monitor jobs continuously
```bash
run4ever --ps
```

### Run command with timeout
```bash
run4ever -d 30 -t 60 my-long-running-command.sh
```
This will run the command every 30 seconds, but kill it if it takes longer than 60 seconds.

### Run in background
```bash
run4ever -g -d 60 my-backup-script.sh
```

### Exit on success
```bash
run4ever -d 5 --exit-on-success docker pull nginx
```
This will retry every 5 seconds until the docker pull succeeds, then exit.

### Persist and restore jobs
```bash
# Save a job for later restoration
run4ever -d 60 --persist my-backup-script.sh

# Restore all saved jobs (useful after container restart)
run4ever --restore
```

### Combine multiple features
```bash
run4ever -d 30 -t 120 -m 10 --notify-on failure --notify-method telegram \
  --telegram-token <token> --telegram-chat-id <chat-id> \
  --exit-on-success my-important-command.sh
```
This will:
- Run the command every 30 seconds
- Kill it if it takes longer than 120 seconds
- Retry up to 10 times
- Send Telegram notification on failure
- Exit when the command succeeds

## More Examples

### Using max-retries
Limit the number of retry attempts before giving up.
```bash
run4ever -d 5 -m 3 my-command.sh
```
This will retry the command every 5 seconds, but give up after 3 failed attempts.

### Retry with infinite retries (default)
```bash
run4ever -d 10 -m -1 my-command.sh
```
This will retry indefinitely (default behavior). The `-m -1` flag is optional as it's the default.

### Run with timeout and max retries
```bash
run4ever -d 30 -t 60 -m 5 my-long-running-command.sh
```
This will:
- Run the command every 30 seconds
- Kill it if it takes longer than 60 seconds
- Give up after 5 failed attempts

### Notification on success
Send notification when command succeeds.
```bash
run4ever -d 60 --notify-on success --notify-method desktop ./backup-script.sh
```

### Notification always
Send notification on both success and failure.
```bash
run4ever -d 60 --notify-on always --notify-method telegram \
  --telegram-token <token> --telegram-chat-id <chat-id> ./monitor-script.sh
```

### Slack notification with channel
```bash
run4ever -d 300 --notify-on failure --notify-method slack \
  --slack-webhook-url https://hooks.slack.com/services/YOUR/WEBHOOK/URL \
  --slack-channel "#alerts" ./deployment-script.sh
```

### Email notification for critical tasks
```bash
run4ever -d 3600 --notify-on failure --notify-method email \
  --email-to admin@example.com \
  --email-from run4ever@example.com \
  --email-password your-password \
  --email-smtp smtp.gmail.com \
  --email-port 587 \
  ./critical-backup.sh
```

### Background job with notifications
Run a job in the background and get notified on failures.
```bash
run4ever -g -d 300 --notify-on failure --notify-method telegram \
  --telegram-token <token> --telegram-chat-id <chat-id> \
  ./long-running-task.sh
```

### Persist job with all options
Save a job definition with all configuration options.
```bash
run4ever -d 60 -t 120 -m 10 --notify-on failure --notify-method slack \
  --slack-webhook-url <webhook-url> \
  --exit-on-success --persist ./important-script.sh
```

### Restore all persisted jobs
After a system restart, restore all previously saved jobs.
```bash
run4ever --restore
```

### Verbose mode for debugging
Enable verbose output to see what's happening.
```bash
run4ever -v -d 30 --notify-on always --notify-method desktop ./test-script.sh
```

### Real-world use cases

#### Keep SSH tunnel alive
```bash
run4ever -d 10 ssh -N -R 8080:localhost:80 user@remote-server.com
```

#### Monitor disk space
```bash
run4ever -d 300 --notify-on failure --notify-method telegram \
  --telegram-token <token> --telegram-chat-id <chat-id> \
  df -h | grep -E '^/dev' | awk '{if ($5+0 > 90) exit 1}'
```

#### Retry database connection
```bash
run4ever -d 5 --exit-on-success --notify-on success --notify-method desktop \
  mysql -h db.example.com -u user -p'password' -e "SELECT 1"
```

#### Keep Docker container running
```bash
run4ever -d 30 docker start my-container
```

#### Health check with timeout
```bash
run4ever -d 60 -t 10 --notify-on failure --notify-method slack \
  --slack-webhook-url <webhook-url> \
  curl -f http://localhost:8080/health
```

#### Scheduled backup with email notification
```bash
run4ever -d 86400 --notify-on failure --notify-method email \
  --email-to admin@example.com \
  --email-from backup@example.com \
  --email-password <password> \
  --email-smtp smtp.example.com \
  --email-port 587 \
  ./daily-backup.sh
```
