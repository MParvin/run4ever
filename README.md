# run4ever
`run4ever` is a command-line tool that allows you to run a specified command indefinitely, with a specified delay between each execution.

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
    -d or --delay: Specify the delay in seconds between command executions. Default is 10 seconds.
    -v or --verbose: Enable verbose mode. This will cause run4ever to print additional output such as errors and confirmation messages.
```

## Example
```bash
./run4ever ssh -R 80:localhost:80 user@server.example.com
```
The above command will run the command `ssh -R 80:localhost:80 user@server.example.com` for ever, and If  the command fails, it will be restarted after 10 seconds.

```bash
./run4ever -d 5 ssh -D 1080 user@server.example.com
```
The above command will setup a SOCKS proxy on port 1080, and If  the command fails, it will be restarted after 5 seconds.

```bash
./run4ever -d 3600 my-backup-script.sh
```
This will run the script `my-backup-script.sh` every hour.

```bash
run4ever -d 60 check-service-status.sh my-service
```
Monitor the status of a service every minute.

```bash
run4ever -d 300 rsync -avz --delete /home/user/ /mnt/backup
```
This will run the command `rsync -avz --delete /home/user/ /mnt/backup` every 5 minutes.

```bash
./run4ever -d 300 my-long-running-job.sh
```
Run a long-running job that takes a long time to complete, but which needs to be restarted if it fails


## Description
`run4ever` is a command-line tool that allows you to run a specified command indefinitely, with a specified delay between each execution.


#### Todo
- [X] Fix flags conflict between run4ever and command (Fixed)
- [ ] Add bash completion
- [ ] Add github actions
- [ ] Add tests