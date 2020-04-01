# ![](https://fonts.gstatic.com/s/i/materialicons/bookmarks/v4/24px.svg) Processes
[![](https://img.shields.io/github/v/release/codemodify/systemkit-processes?style=flat-square)](https://github.com/codemodify/systemkit-processes/releases/latest)
![](https://img.shields.io/github/languages/code-size/codemodify/systemkit-processes?style=flat-square)
![](https://img.shields.io/github/last-commit/codemodify/systemkit-processes?style=flat-square)
[![](https://img.shields.io/badge/license-0--license-brightgreen?style=flat-square)](https://github.com/codemodify/TheFreeLicense)

![](https://img.shields.io/github/workflow/status/codemodify/systemkit-processes/qa?style=flat-square)
![](https://img.shields.io/github/issues/codemodify/systemkit-processes?style=flat-square)
[![](https://goreportcard.com/badge/github.com/codemodify/systemkit-processes?style=flat-square)](https://goreportcard.com/report/github.com/codemodify/systemkit-processes)

[![](https://img.shields.io/badge/godoc-reference-brightgreen?style=flat-square)](https://godoc.org/github.com/codemodify/systemkit-processes)
![](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)
![](https://img.shields.io/gitter/room/codemodify/systemkit-processes?style=flat-square)

![](https://img.shields.io/github/contributors/codemodify/systemkit-processes?style=flat-square)
![](https://img.shields.io/github/stars/codemodify/systemkit-processes?style=flat-square)
![](https://img.shields.io/github/watchers/codemodify/systemkit-processes?style=flat-square)
![](https://img.shields.io/github/forks/codemodify/systemkit-processes?style=flat-square)

#### Robust system process management, spawn and monitor,

#### Supported: Linux, Raspberry Pi, FreeBSD, Mac OS, Windows, Solaris

# ![](https://fonts.gstatic.com/s/i/materialicons/bookmarks/v4/24px.svg) Install
```go
go get github.com/codemodify/systemkit-processes
```


# ![](https://fonts.gstatic.com/s/i/materialicons/bookmarks/v4/24px.svg) API

&nbsp;								| &nbsp;
---									| ---
find.ProcessByPID(_pid_)			| Find process by PID
find.AllProcesses()					| Fetches a snapshot of all running processes
&nbsp;								|
procMon := `monitor.New()`			| Create a new process monitor
procMon.`Spawn`(_id_)				| Spawns and monitors a process, tags it with ID
procMon.`Start`(_id_)				| Starts the process taged with ID
procMon.`Stop`(_id_)				| Stop the process taged with ID
procMon.`Restart`(_id_)				| Restart the process taged with ID
procMon.`StopAl`l()					| Stops all monitored processes
procMon.`GetRuningProcess`(_id_)	| Gets the running process
procMon.`RemoveFromMonitor`(_id_)	| Removes a process from being monitred
procMon.`GetAllIDs`(_id_)			| Returns tags for all monitored processes
&nbsp;								|
proc.`Start`(_id_)					| Starts the process
proc.`IsRunning`(_id_)				| `true` if process is running
proc.`ExitCode`(_id_)				| Returns the exit code
proc.`StartedAt`(_id_)				| Started time
proc.`StoppedAt`(_id_)				| Stopped time
proc.`PID`(_id_)					| Returns PID
proc.`ParentPID`(_id_)				| Returns parent PID
proc.`Stop`(_id_)					| Stops the process (kills it if needed)
proc.`Details`(_id_)				| Details about the process, like executable name
