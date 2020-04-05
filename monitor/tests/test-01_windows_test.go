// +build windows

package tests

import (
	"fmt"
	"testing"
	"time"

	logging "github.com/codemodify/systemkit-logging"

	"github.com/codemodify/systemkit-processes/contracts"
	procMon "github.com/codemodify/systemkit-processes/monitor"
)

func Test_02(t *testing.T) {
	const logID = "Test_01"

	processID := "test-id"

	monitor := procMon.New()
	monitor.Spawn(processID, contracts.ProcessTemplate{
		Executable: "notepad.exe",
		Args:       []string{},
		OnStdOut: func(data []byte) {
			logging.Instance().Debugf("%s: OnStdOut: %v", logID, string(data))
		},
		OnStdErr: func(data []byte) {
			logging.Instance().Debugf("%s: OnStdErr: %v", logID, string(data))
		},
	})

	logging.Instance().Infof(
		"%s: IsRunning: %v, ExitCode: %v, StartedAt: %v, StoppedAt: %v",
		logID,
		monitor.GetRuningProcess(processID).IsRunning(),
		monitor.GetRuningProcess(processID).ExitCode(),
		monitor.GetRuningProcess(processID).StartedAt(),
		monitor.GetRuningProcess(processID).StoppedAt(),
	)

	// WAIT 5 seconds
	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				logging.Instance().Debugf("%s: Tick at, %v", logID, t)
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

	logging.Instance().Infof(
		"%s: IsRunning: %v, ExitCode: %v, StartedAt: %v, StoppedAt: %v",
		logID,
		monitor.GetRuningProcess(processID).IsRunning(),
		monitor.GetRuningProcess(processID).ExitCode(),
		monitor.GetRuningProcess(processID).StartedAt(),
		monitor.GetRuningProcess(processID).StoppedAt(),
	)
}
