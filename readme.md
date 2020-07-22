# ![](https://fonts.gstatic.com/s/i/materialiconsoutlined/flare/v4/24px.svg) Processes
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

&nbsp;										| &nbsp;
---											| ---
find.ProcessByPID(_pid_)					| Find process by PID
find.AllProcesses()							| Fetches a snapshot of all running processes
&nbsp;										|
procMon := `monitor.New()`					| Create a new process monitor
procMon.`Spawn`(_template_)					| Spawns and monitors a process based on a template, generates a tag
procMon.`SpawnWithTag`(_template_, _tag_)	| Spawns and monitors a process based on a template and custom tag
procMon.`Start`(_tag_)						| Starts the process taged with ID
procMon.`Stop`(_tag_)						| Stop the process taged with ID
procMon.`Restart`(_tag_)					| Restart the process taged with ID
procMon.`StopAl`l()							| Stops all monitored processes
procMon.`GetProcess`(_tag_)					| Gets the running process
procMon.`RemoveFromMonitor`(_tag_)			| Removes a process from being monitred
procMon.`GetAllTags`()						| Returns tags for all monitored processes
&nbsp;										|
proc.`Start`()								| Starts the process
proc.`Stop`()								| Stops the process (kills it if needed)
proc.`IsRunning`()							| `true` if process is running
proc.`Details`()							| Details about the process, like PID, executable name
proc.`ExitCode`()							| Returns the exit code
proc.`StartedAt`()							| Started time
proc.`StoppedAt`()							| Stopped time

proc.`OnStdOut`()							| Set reader for process STDOUT
proc.`OnStdErr`()							| Set reader for process STDERR
proc.`OnStop`()								| Set handler when the process stops
