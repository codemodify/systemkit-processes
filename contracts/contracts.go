package contracts

import "time"

// ProcessDoesNotExist -
const ProcessDoesNotExist = -1

// ProcessOutputReader -
type ProcessOutputReader func([]byte)

// ProcessTemplate -
type ProcessTemplate struct {
	Executable       string
	Args             []string
	WorkingDirectory string
	Environment      []string
	OnStdOut         ProcessOutputReader `json:"-"`
	OnStdErr         ProcessOutputReader `json:"-"`
}

// RuningProcess - represents a running process
type RuningProcess interface {
	Start() error
	IsRunning() bool
	ExitCode() int
	StartedAt() time.Time
	StoppedAt() time.Time
	PID() int
	ParentPID() int
	Stop() error
	Details() ProcessTemplate
}

// Monitor - process monitor
type Monitor interface {
	Spawn(id string, process ProcessTemplate) error
	Start(id string) error
	Stop(id string) error
	Restart(id string) error
	StopAll() []error
	GetRuningProcess(id string) RuningProcess
	RemoveFromMonitor(id string)
	GetAllIDs() []string
}
