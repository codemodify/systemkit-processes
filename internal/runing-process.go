package internal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"

	helpersStrings "github.com/codemodify/systemkit-helpers-conv"
	helpersReflect "github.com/codemodify/systemkit-helpers-reflection"
	logging "github.com/codemodify/systemkit-logging"
	"github.com/codemodify/systemkit-processes/contracts"
)

const logID = "PROCESS"

// processDoesNotExist -
const processDoesNotExist = -1

type runingProcess struct {
	processTemplate contracts.ProcessTemplate
	osCmd           *exec.Cmd
	startedAt       time.Time
	stoppedAt       time.Time
	stdOut          io.ReadCloser
	stdErr          io.ReadCloser
}

// NewEmptyRuningProcess -
func NewEmptyRuningProcess() contracts.RuningProcess {
	return NewRuningProcess(contracts.ProcessTemplate{})
}

// NewRuningProcess -
func NewRuningProcess(processTemplate contracts.ProcessTemplate) contracts.RuningProcess {
	return &runingProcess{
		processTemplate: processTemplate,
		osCmd:           nil,
		startedAt:       time.Unix(0, 0),
		stoppedAt:       time.Unix(0, 0),
	}
}

// NewRuningProcessWithOSProc -
func NewRuningProcessWithOSProc(processTemplate contracts.ProcessTemplate, osProc *os.Process) contracts.RuningProcess {
	r := &runingProcess{
		processTemplate: processTemplate,
		osCmd:           exec.Command(processTemplate.Executable, processTemplate.Args...),
		startedAt:       time.Unix(0, 0),
		stoppedAt:       time.Unix(0, 0),
	}

	r.osCmd.Process = osProc

	return r
}

// Start -
func (thisRef *runingProcess) Start() error {
	thisRef.osCmd = exec.Command(thisRef.processTemplate.Executable, thisRef.processTemplate.Args...)

	// set working folder
	if !helpersStrings.IsNullOrEmpty(thisRef.processTemplate.WorkingDirectory) {
		thisRef.osCmd.Dir = thisRef.processTemplate.WorkingDirectory
	}

	// set env
	if thisRef.processTemplate.Environment != nil {
		thisRef.osCmd.Env = thisRef.processTemplate.Environment
	}

	// capture STDOUT, STDERR
	stdOutPipe, err := thisRef.osCmd.StdoutPipe()
	if err != nil {
		logging.Instance().Errorf("%s: get-StdOut-FAIL for [%s], [%s] @ %s", logID, thisRef.processTemplate.Executable, err.Error(), helpersReflect.GetThisFuncName())
		return err
	}
	thisRef.stdOut = stdOutPipe

	stdErrPipe, err := thisRef.osCmd.StderrPipe()
	if err != nil {
		logging.Instance().Errorf("%s: get-StdErr-FAIL for [%s], [%s] @ %s", logID, thisRef.processTemplate.Executable, err.Error(), helpersReflect.GetThisFuncName())
		return err
	}
	thisRef.stdErr = stdErrPipe

	// start
	logging.Instance().Debugf("%s: start %s @ %s", logID, helpersStrings.AsJSONString(thisRef.processTemplate), helpersReflect.GetThisFuncName())

	err = thisRef.osCmd.Start()
	if err != nil {
		thisRef.stoppedAt = time.Now()

		detailedErr := fmt.Errorf("%s: start-FAILED %s, %s @ %s", logID, helpersStrings.AsJSONString(thisRef.processTemplate), err.Error(), helpersReflect.GetThisFuncName())
		logging.Instance().Error(detailedErr.Error())

		return detailedErr
	}

	return nil
}

// Stop - stops the process
func (thisRef *runingProcess) Stop() error {
	if thisRef.osCmd == nil || thisRef.osCmd.Process == nil {
		return nil
	}

	var err error

	count := 0
	maxStopAttempts := 20
	for {
		// try #
		count++
		if count > maxStopAttempts {
			logging.Instance().Errorf("%s: stop-FAIL [%s] with PID [%d] @ %s", logID, thisRef.processTemplate.Executable, thisRef.processID(), helpersReflect.GetThisFuncName())
			break
		}

		// break if DONE
		if !thisRef.IsRunning() {
			logging.Instance().Debugf("%s: stop-SUCCESS [%s] @ %s", logID, thisRef.processTemplate.Executable, helpersReflect.GetThisFuncName())
			break
		}

		// log the attempt #
		logging.Instance().Debugf("%s: stop-ATTEMPT #%d to stop [%s] @ %s", logID, count, thisRef.processTemplate.Executable, helpersReflect.GetThisFuncName())

		thisRef.osCmd.Process.Signal(syscall.SIGINT)
		thisRef.osCmd.Process.Signal(syscall.SIGTERM)
		thisRef.osCmd.Process.Signal(syscall.SIGKILL)
		processKillHelper(thisRef.osCmd.Process.Pid)

		err = thisRef.osCmd.Process.Kill()

		time.Sleep(500 * time.Millisecond)
		thisRef.osCmd.Process.Wait()
	}

	thisRef.stoppedAt = time.Now()

	return err
}

// IsRunning - tells if the process is running
func (thisRef runingProcess) IsRunning() bool {
	pid := thisRef.processID()
	if pid == processDoesNotExist {
		return false
	}

	rp := thisRef.Details()

	return (rp.State != contracts.ProcessStateNonExistent &&
		rp.State != contracts.ProcessStateObsolete &&
		rp.State != contracts.ProcessStateDead)
}

// Details - return processTemplate about the process
func (thisRef runingProcess) Details() contracts.RuntimeProcess {
	rpByPID, err := getRuntimeProcessByPID(thisRef.processID())
	if err != nil {
		return contracts.RuntimeProcess{
			State: contracts.ProcessStateNonExistent,
		}
	}

	return rpByPID
}

// ExitCode -
func (thisRef runingProcess) ExitCode() int {
	if thisRef.osCmd == nil || thisRef.osCmd.Process == nil || thisRef.osCmd.ProcessState == nil {
		return 0
	}

	return thisRef.osCmd.ProcessState.ExitCode()
}

// StartedAt - returns the time when the process was started
func (thisRef runingProcess) StartedAt() time.Time {
	if thisRef.osCmd == nil || thisRef.osCmd.Process == nil {
		return time.Unix(0, 0)
	}

	return thisRef.startedAt
}

// StoppedAt - returns the time when the process was stopped
func (thisRef runingProcess) StoppedAt() time.Time {
	if thisRef.osCmd == nil || thisRef.osCmd.Process == nil {
		return time.Unix(0, 0)
	}

	return thisRef.stoppedAt
}

func (thisRef runingProcess) OnStdOut(outputReader contracts.ProcessOutputReader) {
	logging.Instance().Debugf("%s: read-StdOut for [%s] @ %s", logID, thisRef.processTemplate.Executable, helpersReflect.GetThisFuncName())

	if outputReader != nil {
		go func() {
			err := readOutput(thisRef.stdOut, outputReader)
			if err != nil {
				logging.Instance().Warningf("%s: read-StdOut-FAIL for [%s], [%s] @ %s", logID, thisRef.processTemplate.Executable, err.Error(), helpersReflect.GetThisFuncName())
			}

			logging.Instance().Debugf("%s: read-StdOut-SUCESS for [%s]  @ %s", logID, thisRef.processTemplate.Executable, helpersReflect.GetThisFuncName())
		}()
	}
}

func (thisRef runingProcess) OnStdErr(outputReader contracts.ProcessOutputReader) {
	logging.Instance().Debugf("%s: read-StdErr for [%s] @ %s", logID, thisRef.processTemplate.Executable, helpersReflect.GetThisFuncName())

	if outputReader != nil {
		go func() {
			err := readOutput(thisRef.stdErr, outputReader)
			if err != nil {
				logging.Instance().Warningf("%s: read-StdErr-FAIL for [%s], [%s] @ %s", logID, thisRef.processTemplate.Executable, err.Error(), helpersReflect.GetThisFuncName())
			}

			logging.Instance().Debugf("%s: read-StdErr-SUCESS for [%s]  @ %s", logID, thisRef.processTemplate.Executable, helpersReflect.GetThisFuncName())
		}()
	}
}

func (thisRef *runingProcess) OnStop(stoppedDelegate contracts.ProcessStoppedDelegate) {
	go func() {
		for {
			time.Sleep(1 * time.Second)

			if !thisRef.IsRunning() {
				thisRef.Stop() // call this because .osCmd.Process.Wait() is needed
				if stoppedDelegate != nil {
					stoppedDelegate()
				}

				return
			}
		}
	}()
}

func (thisRef runingProcess) processID() int {
	if thisRef.osCmd == nil || thisRef.osCmd.Process == nil {
		return processDoesNotExist
	}

	return thisRef.osCmd.Process.Pid
}

func readOutput(readerCloser io.ReadCloser, outputReader contracts.ProcessOutputReader) error {
	reader := bufio.NewReader(readerCloser)
	line, _, err := reader.ReadLine()
	for err != io.EOF {
		outputReader(line)
		line, _, err = reader.ReadLine()
	}

	if err == io.EOF {
		return nil
	}

	return err
}
