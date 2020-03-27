// +build windows

package tests

import (
	"fmt"
	"testing"
	"time"

	logging "github.com/codemodify/systemkit-logging"
	loggingC "github.com/codemodify/systemkit-logging/contracts"
	loggingP "github.com/codemodify/systemkit-logging/persisters"

	"github.com/codemodify/systemkit-processes/contracts"
	procMon "github.com/codemodify/systemkit-processes/monitor"
)

func Test_01(t *testing.T) {
	const logID = "Test_01"

	logging.Init(logging.NewEasyLoggerForLogger(loggingP.NewFileLogger(loggingC.TypeDebug, "log1.log")))

	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"message": fmt.Sprintf("%s: START", logID),
	})

	processID := "test-id"

	monitor := procMon.New()
	monitor.Spawn(processID, contracts.ProcessTemplate{
		Executable: "notepad.exe",
		Args:       []string{},
		OnStdOut: func(data []byte) {
			logging.Instance().LogDebugWithFields(loggingC.Fields{
				"message": fmt.Sprintf("%s: OnStdOut: %v", logID, string(data)),
			})
		},
		OnStdErr: func(data []byte) {
			logging.Instance().LogDebugWithFields(loggingC.Fields{
				"message": fmt.Sprintf("%s: OnStdErr: %v", logID, string(data)),
			})
		},
	})

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"message": fmt.Sprintf(
			"%s: IsRunning: %v, ExitCode: %v, StartedAt: %v, StoppedAt: %v",
			logID,
			monitor.GetRuningProcess(processID).IsRunning(),
			monitor.GetRuningProcess(processID).ExitCode(),
			monitor.GetRuningProcess(processID).StartedAt(),
			monitor.GetRuningProcess(processID).StoppedAt(),
		),
	})

	// WAIT 5 seconds
	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				logging.Instance().LogDebugWithFields(loggingC.Fields{
					"message": fmt.Sprintf("%s: Tick at, %v", logID, t),
				})
			}
		}
	}()
	time.Sleep(5 * time.Second)
	ticker.Stop()
	done <- true

	// STOP
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"message": fmt.Sprintf("%s: STOP", logID),
	})

	monitor.Stop(processID)

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"message": fmt.Sprintf(
			"%s: IsRunning: %v, ExitCode: %v, StartedAt: %v, StoppedAt: %v",
			logID,
			monitor.GetRuningProcess(processID).IsRunning(),
			monitor.GetRuningProcess(processID).ExitCode(),
			monitor.GetRuningProcess(processID).StartedAt(),
			monitor.GetRuningProcess(processID).StoppedAt(),
		),
	})
}
