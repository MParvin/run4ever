# run4ever
This program will get a command and run this, if an error occures or program exit, it will run it again

##### Build
```bash
go build -o run4ever
```

##### Usage
```bash
./run4ever [command]
```

##### Example
```bash
./run4ever ssh -R 80:localhost:80 user@server.example.com
```

```bash
./run4ever ssh -D 1080 user@server.example.com
```

#### Known issues
- Unfortunatly, has conflict between this program flags and command flags, so you can't use flags for running command **at this time ** , just use command without flags or use shell script as command
I'm working on it ;)

#### Todo
- [ ] Fix flags conflict
- [ ] Add bash completion
