package contracts

import (
	"errors"
	"time"
)

// ErrProcessDoesNotExist -
var ErrProcessDoesNotExist = errors.New("ErrProcessDoesNotExist")

// ProcessState -
type ProcessState int

// ProcessStateWaitingIO -
const (
	// UNIX
	ProcessStateWaitingIO    ProcessState = iota // 0 -> D - Uninterruptible sleep (usually IO)
	ProcessStateRunning                          // 1 -> R - Running or runnable (on run queue)
	ProcessStateWaitingEvent                     // 2 -> S - Interruptible sleep (waiting for an event to complete)
	ProcessStateTraced                           // 3 -> T - Stopped, either by a job control signal or because it is being traced
	ProcessStatePaging                           // 4 -> W - paging (not valid since the 2.6.xx kernel)
	ProcessStateDead                             // 5 -> X - dead (should never be seen)
	ProcessStateObsolete                         // 6 -> Z - Defunct ("zombie" / "obsolete") process, terminated but not reaped by its parent

	// EXTEND
	ProcessStateNonExistent // process does not exist
	ProcessStateUnknown
)

// String - stringer interface
func (thisRef ProcessState) String() string {
	switch thisRef {
	case ProcessStateWaitingIO:
		return "ProcessStateWaitingIO"
	case ProcessStateRunning:
		return "ProcessStateRunning"
	case ProcessStateWaitingEvent:
		return "ProcessStateWaitingEvent"
	case ProcessStateTraced:
		return "ProcessStateTraced"
	case ProcessStatePaging:
		return "ProcessStatePaging"
	case ProcessStateDead:
		return "ProcessStateDead"
	case ProcessStateObsolete:
		return "ProcessStateObsolete"

	default:
		return "ProcessStateNonExistent"
	}
}

// RuntimeProcess -
type RuntimeProcess struct {
	Executable       string       `json:"executable"`
	ExecutableName   string       `json:"executableName"`
	Args             []string     `json:"args"`
	WorkingDirectory string       `json:"workingDirectory"`
	Environment      []string     `json:"environment"`
	ProcessID        int          `json:"processID"`
	ParentProcessID  int          `json:"parentProcessID"`
	UserID           int          `json:"userID"`
	GroupID          int          `json:"groupID"`
	State            ProcessState `json:"state"`

	// FIXME
	sessionID       int `json:"-"`
	effectiveUserID int `json:"-"`
}

// RuningProcess - represents a running process
type RuningProcess interface {
	Start() error
	Stop(tag string, attempts int, waitTimeout time.Duration) error
	IsRunning() bool
	Details() RuntimeProcess

	ExitCode() int
	StartedAt() time.Time
	StoppedAt() time.Time

	OnStdOut(outputReader ProcessOutputReader, params interface{})
	OnStdErr(outputReader ProcessOutputReader, params interface{})
	OnStop(stoppedDelegate ProcessStoppedDelegate, params interface{})
}
