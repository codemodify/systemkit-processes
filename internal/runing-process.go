package internal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"

	logging "github.com/codemodify/systemkit-logging"
	"github.com/codemodify/systemkit-processes/contracts"
	"github.com/codemodify/systemkit-processes/helpers"
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
	if !helpers.IsNullOrEmpty(thisRef.processTemplate.WorkingDirectory) {
		thisRef.osCmd.Dir = thisRef.processTemplate.WorkingDirectory
	}

	// set env
	if thisRef.processTemplate.Environment != nil {
		thisRef.osCmd.Env = thisRef.processTemplate.Environment
	}

	// capture STDERR
	stdOutPipe, err := thisRef.osCmd.StdoutPipe()
	if err != nil {
		logging.Errorf("%s: get-StdOut-FAIL for [%s], [%s]", logID, thisRef.processTemplate.Executable, err.Error())
		return err
	}
	thisRef.stdOut = stdOutPipe

	// capture STDERR
	stdErrPipe, err := thisRef.osCmd.StderrPipe()
	if err != nil {
		logging.Errorf("%s: get-StdErr-FAIL for [%s], [%s]", logID, thisRef.processTemplate.Executable, err.Error())
		return err
	}
	thisRef.stdErr = stdErrPipe

	thisRef.osCmd.SysProcAttr = procAttrs

	// start
	logging.Debugf("%s: start %s", logID, helpers.AsJSONString(thisRef.processTemplate))

	err = thisRef.osCmd.Start()
	if err != nil {
		thisRef.stoppedAt = time.Now()

		detailedErr := fmt.Errorf("%s: start-FAILED %s, %s", logID, helpers.AsJSONString(thisRef.processTemplate), err.Error())
		logging.Error(detailedErr.Error())

		return detailedErr
	}

	thisRef.startedAt = time.Now()

	return nil
}

// Stop - stops the process
func (thisRef *runingProcess) Stop(tag string, attempts int, waitTimeout time.Duration) error {
	if thisRef.osCmd == nil || thisRef.osCmd.Process == nil {
		return nil
	}

	if !thisRef.IsRunning() {
		return nil
	}

	// go func() {
	// 	if thisRef.stdOut != nil {
	// 		thisRef.stdOut.Close()
	// 	}

	// 	if thisRef.stdErr != nil {
	// 		thisRef.stdErr.Close()
	// 	}
	// }()

	defer func() {
		logging.Debugf("%s: STOP-END %s", logID, tag)
	}()

	logging.Debugf("%s: STOP-START %s", logID, tag)

	var err error
	count := 0
	maxStopAttempts := 20
	for {
		// try #
		count++
		if count > maxStopAttempts {
			logging.Errorf("%s: stop-FAIL [%s] with PID [%d]", logID, thisRef.processTemplate.Executable, thisRef.processID())
			break
		}

		// log the STOP attempt #

		for i := 0; i < attempts; i++ {
			logging.Debugf("%s: stop-ATTEMPT-SIGINT #%d to stop [%s]", logID, i, thisRef.processTemplate.Executable)
			thisRef.osCmd.Process.Signal(syscall.SIGINT) // this works on all except on Windows
			sendCtrlC(thisRef.osCmd.Process.Pid)         // this works on Windows

			time.Sleep(waitTimeout)
			if !thisRef.IsRunning() {
				thisRef.osCmd.Process.Wait()
				thisRef.stoppedAt = time.Now()
				logging.Debugf("%s: stop-SUCCESS [%s]", logID, thisRef.processTemplate.Executable)
				return nil
			}
		}

		for i := 0; i < attempts; i++ {
			logging.Debugf("%s: stop-ATTEMPT-SIGTERM #%d to stop [%s]", logID, i, thisRef.processTemplate.Executable)
			thisRef.osCmd.Process.Signal(syscall.SIGTERM)
			time.Sleep(waitTimeout)
			if !thisRef.IsRunning() {
				thisRef.osCmd.Process.Wait()
				thisRef.stoppedAt = time.Now()
				logging.Debugf("%s: stop-SUCCESS [%s]", logID, thisRef.processTemplate.Executable)
				return nil
			}
		}

		for i := 0; i < attempts; i++ {
			logging.Debugf("%s: stop-ATTEMPT-SIGKILL #%d to stop [%s]", logID, i, thisRef.processTemplate.Executable)
			thisRef.osCmd.Process.Signal(syscall.SIGKILL)
			time.Sleep(waitTimeout)
			if !thisRef.IsRunning() {
				thisRef.osCmd.Process.Wait()
				thisRef.stoppedAt = time.Now()
				logging.Debugf("%s: stop-SUCCESS [%s]", logID, thisRef.processTemplate.Executable)
				return nil
			}
		}

		for i := 0; i < attempts; i++ {
			logging.Debugf("%s: stop-ATTEMPT-aggressive-kill-1 #%d to stop [%s]", logID, i, thisRef.processTemplate.Executable)
			processKillHelper(thisRef.osCmd.Process.Pid)
			time.Sleep(waitTimeout)
			if !thisRef.IsRunning() {
				thisRef.osCmd.Process.Wait()
				thisRef.stoppedAt = time.Now()
				logging.Debugf("%s: stop-SUCCESS [%s]", logID, thisRef.processTemplate.Executable)
				return nil
			}
		}

		for i := 0; i < attempts; i++ {
			logging.Debugf("%s: stop-ATTEMPT-aggressive-kill-2 #%d to stop [%s]", logID, i, thisRef.processTemplate.Executable)
			err = thisRef.osCmd.Process.Kill()
			time.Sleep(waitTimeout)
			if !thisRef.IsRunning() {
				thisRef.osCmd.Process.Wait()
				thisRef.stoppedAt = time.Now()
				logging.Debugf("%s: stop-SUCCESS [%s]", logID, thisRef.processTemplate.Executable)
				return nil
			}
		}
	}

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
		rp.State != contracts.ProcessStateDead &&
		rp.State != contracts.ProcessStateUnknown)
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

func (thisRef runingProcess) OnStdOut(outputReader contracts.ProcessOutputReader, params interface{}) {
	logging.Debugf("%s: read-StdOut for [%s]", logID, thisRef.processTemplate.Executable)

	if outputReader != nil {
		go func() {
			err := readOutput(thisRef.stdOut, outputReader, params)
			if err != nil {
				logging.Warningf("%s: read-StdOut-FAIL for [%s], [%s]", logID, thisRef.processTemplate.Executable, err.Error())
			}

			logging.Debugf("%s: read-StdOut-SUCESS for [%s]", logID, thisRef.processTemplate.Executable)
		}()
	}
}

func (thisRef runingProcess) OnStdErr(outputReader contracts.ProcessOutputReader, params interface{}) {
	logging.Debugf("%s: read-StdErr for [%s]", logID, thisRef.processTemplate.Executable)

	if outputReader != nil {
		go func() {
			err := readOutput(thisRef.stdErr, outputReader, params)
			if err != nil {
				logging.Warningf("%s: read-StdErr-FAIL for [%s], [%s]", logID, thisRef.processTemplate.Executable, err.Error())
			}

			logging.Debugf("%s: read-StdErr-SUCESS for [%s]", logID, thisRef.processTemplate.Executable)
		}()
	}
}

func (thisRef *runingProcess) OnStop(stoppedDelegate contracts.ProcessStoppedDelegate, params interface{}) {
	go func(paramsToPass interface{}) {
		for {
			time.Sleep(1 * time.Second)

			if !thisRef.IsRunning() {
				thisRef.Stop("", 1, 100*time.Millisecond) // call this because .osCmd.Process.Wait() is needed
				if stoppedDelegate != nil {
					stoppedDelegate(paramsToPass)
				}

				return
			}
		}
	}(params)
}

func (thisRef runingProcess) processID() int {
	if thisRef.osCmd == nil || thisRef.osCmd.Process == nil {
		return processDoesNotExist
	}

	return thisRef.osCmd.Process.Pid
}

func readOutput(readerCloser io.ReadCloser, outputReader contracts.ProcessOutputReader, params interface{}) error {
	reader := bufio.NewReader(readerCloser)
	line, _, err := reader.ReadLine()
	for {
		if err != nil {
			break
		}

		outputReader(params, line)
		line, _, err = reader.ReadLine()
	}

	if err == io.EOF {
		return nil
	}

	return err
}
