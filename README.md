# process_list
List Process In Table, Search and Filter by Name, PID, PPID, User. Provides Command Line Argumentss and Env Variables.
Idea came from searching for how binaries are executed and loaded on a given OS. 
Can lead to great finds in environment variables and hardcoded passwords on the command line

## To Build

```bash
go build -v -a -ldflags="-w -s" -o proclist main.go

```

### Building on iOS
```
CGO_ENABLED=1 GOOS=ios GOARCH=arm64 go build -v -a -ldflags="-w -s"  -o proclist_ios main.go
codesign -f -o runtime --timestamp -s "Developer ID Application: YOUR NAME (TEAM_ID)" proclist_ios
```

### Building on Android
```bash
GOARCH=arm GOOS=linux go build -v -a -ldflags="-w -s" -o proclist_android main.go
```


## Running

```bash
./proclist -h
Usage of ./proclist:
  -name string
        Process name to search
  -pid int
        PID for process to search (default -1)
  -ppid int
        PPID for parent process to search (default -1)
  -user string
        User to search for processes under
```
