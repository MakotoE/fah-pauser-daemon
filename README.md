Pauses Folding@home when a specified program is running. This allows FAH to run "While I'm working" but also makes FAH stop for resource-intensive apps.

```
go get -u -ldflags -H=windowsgui github.com/MakotoE/fah-pauser
```

Create `~/.config/fah-pauser-daemon.yml` (or `%userprofile%\.config\fah-pauser-daemon.yml`) and list programs that should pause FAH when any of them are running.

```
PauseOn:
- devenv.exe
- rFactor2.exe
```

The system processes list is polled once every five minutes.