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
    -w or --watch: Show a list of running commands and their PIDs.
    -d or --delay: Specify the delay in seconds between command executions. Default is 10 seconds.
    -v or --verbose: Enable verbose mode. This will cause run4ever to print additional output such as errors and confirmation messages.
```

## Example
All examples are in [examples](examples) directory.

## Description
`run4ever` is a command-line tool that allows you to run a specified command indefinitely, with a specified delay between each execution.


#### Todo
- [X] Fix flags conflict between run4ever and command (Fixed)
- [X] Add watch flag
- [ ] Move runInfinitely function to tools package
- [ ] Add github actions
- [ ] Add tests
- [ ] Add timeout flag
- [ ] Add `-g` flag to run command in background
- [ ] Add `-l` flag to list all running jobs
- [ ] Add more examples
- [ ] Add bash completion

