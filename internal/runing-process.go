package internal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"

	helpersReflect "github.com/codemodify/systemkit-helpers"
	helpersStrings "github.com/codemodify/systemkit-helpers"
	logging "github.com/codemodify/systemkit-logging"
	loggingC "github.com/codemodify/systemkit-logging/contracts"
	"github.com/codemodify/systemkit-processes/contracts"
)

const logID = "SKIT-PROCESS"

type runingProcessImpl struct {
	processTemplate contracts.ProcessTemplate
	osCmd           *exec.Cmd
	startedAt       time.Time
	stoppedAt       time.Time
	parentPID       int
	lastError       error
}

// NewRuningProcess -
func NewRuningProcess(processTemplate contracts.ProcessTemplate) contracts.RuningProcess {
	return &runingProcessImpl{
		processTemplate: processTemplate,
		osCmd:           nil,
		startedAt:       time.Unix(0, 0),
		stoppedAt:       time.Unix(0, 0),
		parentPID:       -1,
		lastError:       nil,
	}
}

// NewRuningProcessWithOSProc -
func NewRuningProcessWithOSProc(processTemplate contracts.ProcessTemplate, osProc *os.Process) contracts.RuningProcess {
	r := &runingProcessImpl{
		processTemplate: processTemplate,
		osCmd:           exec.Command(processTemplate.Executable, processTemplate.Args...),
		startedAt:       time.Unix(0, 0),
		stoppedAt:       time.Unix(0, 0),
		parentPID:       -1,
		lastError:       nil,
	}

	r.osCmd.Process = osProc

	return r
}

// Start -
func (thisRef *runingProcessImpl) Start() error {
	thisRef.osCmd = exec.Command(thisRef.processTemplate.Executable, thisRef.processTemplate.Args...)

	// set working folder
	if !helpersStrings.IsNullOrEmpty(thisRef.processTemplate.WorkingDirectory) {
		thisRef.osCmd.Dir = thisRef.processTemplate.WorkingDirectory
	}

	// set env
	if thisRef.processTemplate.Environment != nil {
		thisRef.osCmd.Env = thisRef.processTemplate.Environment
	}

	// set stderr and stdout
	stdOutPipe, err := thisRef.osCmd.StdoutPipe()
	if err != nil {
		detailedErr := fmt.Errorf("%s: failed to get StdoutPipe for [%s], details [%s]", logID, thisRef.processTemplate.Executable, err.Error())

		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": detailedErr.Error(),
		})

		return detailedErr
	}

	stdErrPipe, err := thisRef.osCmd.StderrPipe()
	if err != nil {
		detailedErr := fmt.Errorf("%s, failed to get StderrPipe for [%s], details [%s]", logID, thisRef.processTemplate.Executable, err.Error())

		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": detailedErr.Error(),
		})

		return detailedErr
	}

	if thisRef.processTemplate.OnStdOut != nil {
		go readStdOutFromProc(stdOutPipe, thisRef.processTemplate)
	}
	if thisRef.processTemplate.OnStdErr != nil {
		go readStdErrFromProc(stdErrPipe, thisRef.processTemplate)
	}

	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: starting [%s]", logID, thisRef.processTemplate.Executable),
	})

	err = thisRef.osCmd.Start()
	if err != nil {
		thisRef.lastError = err
		thisRef.stoppedAt = time.Now()

		detailedErr := fmt.Errorf("%s, failed to start [%s], details [%s]", logID, thisRef.processTemplate.Executable, err.Error())

		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": detailedErr.Error(),
		})

		return detailedErr
	}

	// FIXME: fetchDetail like parentPID

	return nil
}

// IsRunning - tells if the process is running
func (thisRef runingProcessImpl) IsRunning() bool {
	if thisRef.osCmd == nil || thisRef.osCmd.Process == nil {
		return false
	}

	runningProcess, err := ProcessByPID(thisRef.osCmd.Process.Pid)
	if err != nil {
		return false
	}

	if runningProcess.PID() == thisRef.PID() {
		return true
	}

	return runningProcess.IsRunning()
}

// ExitCode -
func (thisRef runingProcessImpl) ExitCode() int {
	if thisRef.osCmd == nil || thisRef.osCmd.Process == nil || thisRef.osCmd.ProcessState == nil {
		return 0
	}

	return thisRef.osCmd.ProcessState.ExitCode()
}

// StartedAt - returns the time when the process was started
func (thisRef runingProcessImpl) StartedAt() time.Time {
	if thisRef.osCmd == nil || thisRef.osCmd.Process == nil {
		return time.Unix(0, 0)
	}

	return thisRef.startedAt
}

// StoppedAt - returns the time when the process was stopped
func (thisRef runingProcessImpl) StoppedAt() time.Time {
	if thisRef.osCmd == nil || thisRef.osCmd.Process == nil {
		return time.Unix(0, 0)
	}

	return thisRef.stoppedAt
}

// PID - returns process ID
func (thisRef *runingProcessImpl) PID() int {
	if thisRef.osCmd == nil || thisRef.osCmd.Process == nil {
		return contracts.ProcessDoesNotExist
	}

	return thisRef.osCmd.Process.Pid
}

// ParentPID - returns parent process ID
func (thisRef *runingProcessImpl) ParentPID() int {
	return thisRef.parentPID
}

// Stop - stops the process
func (thisRef *runingProcessImpl) Stop() error {
	if thisRef.osCmd == nil || thisRef.osCmd.Process == nil {
		return nil
	}

	count := 0
	maxStopAttempts := 20
	for {
		count++
		if count > maxStopAttempts {
			thisRef.lastError = fmt.Errorf("%s: can't stop %s with PID %d", logID, thisRef.processTemplate.Executable, thisRef.PID())

			logging.Instance().LogErrorWithFields(loggingC.Fields{
				"method":  helpersReflect.GetThisFuncName(),
				"message": thisRef.lastError.Error(),
			})

			break
		}

		if !thisRef.IsRunning() {
			break
		}

		logging.Instance().LogDebugWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s: attempt #%d to stop [%s]", logID, count, thisRef.processTemplate.Executable),
		})

		thisRef.osCmd.Process.Signal(syscall.SIGINT)
		thisRef.osCmd.Process.Signal(syscall.SIGTERM)
		thisRef.osCmd.Process.Signal(syscall.SIGKILL)
		processKillHelper(thisRef.osCmd.Process.Pid)

		err := thisRef.osCmd.Process.Kill()
		if err != nil {
			thisRef.lastError = err
		}

		time.Sleep(500 * time.Millisecond)
		thisRef.osCmd.Process.Wait()
	}

	thisRef.stoppedAt = time.Now()

	return thisRef.lastError
}

// Details - return processTemplate about the process
func (thisRef runingProcessImpl) Details() contracts.ProcessTemplate {
	return thisRef.processTemplate
}

func readStdOutFromProc(readerCloser io.ReadCloser, processTemplate contracts.ProcessTemplate) {
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: starting to read StdOut for [%s]", logID, processTemplate.Executable),
	})

	reader := bufio.NewReader(readerCloser)
	line, _, err := reader.ReadLine()
	for err != io.EOF {
		processTemplate.OnStdOut(line)
		line, _, err = reader.ReadLine()
	}

	if err != nil {
		logging.Instance().LogWarningWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s: error reading StdOut for [%s], details [%s]", logID, processTemplate.Executable, err.Error()),
		})
	}

	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: finished to read StdOut for [%s]", logID, processTemplate.Executable),
	})
}

func readStdErrFromProc(readerCloser io.ReadCloser, processTemplate contracts.ProcessTemplate) {
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: starting to read StdErr for [%s]", logID, processTemplate.Executable),
	})

	reader := bufio.NewReader(readerCloser)
	line, _, err := reader.ReadLine()
	for err != io.EOF {
		processTemplate.OnStdOut(line)
		line, _, err = reader.ReadLine()
	}

	if err != nil {
		logging.Instance().LogWarningWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s: error reading StdErr for [%s], details [%s]", logID, processTemplate.Executable, err.Error()),
		})
	}

	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: finished to read StdErr for [%s]", logID, processTemplate.Executable),
	})
}
