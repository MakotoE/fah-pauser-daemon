Pauses Folding@home when a specified program is running. This allows FAH to run "While I'm working" but also makes FAH stop for resource-intensive apps.

Create `~/.config/fah-pauser.yml` (or `%userprofile%\.config\fah-pauser.yml`) and list program names that should pause FAH when running.

```
PauseOn:
- devenv.exe
- rFactor2.exe
```