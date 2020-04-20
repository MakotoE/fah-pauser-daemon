Pauses Folding@home while specified program is running.

```
go get -ldflags -Hwindowsgui github.com/MakotoE/fah-pauser
```

```
usage: fah-pauser <path>
Stops Folding@home when <path> is running
```

After `go get`ing, edit your Shortcuts and .desktop files to prefix `fah-pauser` before the command.

FYI: Windows taskbar shortcuts are located in `%AppData%\Microsoft\Internet Explorer\Quick Launch\User Pinned\TaskBar`