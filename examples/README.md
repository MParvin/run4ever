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

TODO:
### Slack Notification
Send a slack message on command success.
```bash
run4ever -d 60 --notify-on success --notify-method slack --slack-token <token> --slack-channel <channel> ./my-backup-script.sh
```

TODO:
### Email Notification
Send an email on command failure.
```bash
run4ever -d 60 --notify-on failure --notify-method email --email-to <to> --email-from <from> --email-password <password> --email-smtp <smtp> --email-port <port> ./my-backup-script.sh
```
