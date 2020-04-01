# ![](https://fonts.gstatic.com/s/i/materialicons/bookmarks/v4/24px.svg)Processes
[![GoDoc](https://godoc.org/github.com/codemodify/SystemKit?status.svg)](https://godoc.org/github.com/codemodify/SystemKit)
[![0-License](https://img.shields.io/badge/license-0--license-brightgreen)](https://github.com/codemodify/TheFreeLicense)
[![Go Report Card](https://goreportcard.com/badge/github.com/codemodify/SystemKit)](https://goreportcard.com/report/github.com/codemodify/SystemKit)
[![Test Status](https://github.com/danawoodman/systemservice/workflows/Test/badge.svg)](https://github.com/danawoodman/systemservice/actions)
![code size](https://img.shields.io/github/languages/code-size/codemodify/systemkit-processes?style=flat-square)

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
