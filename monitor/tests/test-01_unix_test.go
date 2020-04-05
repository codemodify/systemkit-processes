// +build !windows

package tests

import (
	"testing"
	"time"

	logging "github.com/codemodify/systemkit-logging"

	"github.com/codemodify/systemkit-processes/contracts"
	procMon "github.com/codemodify/systemkit-processes/monitor"
)

func Test_01(t *testing.T) {
	const logID = "Test_01"

	logging.Instance().Debugf("%s: START", logID)

	processID := "test-id"

	monitor := procMon.New()
	monitor.Spawn(processID, contracts.ProcessTemplate{
		Executable: "sh",
		Args:       []string{"-c", "while :; do echo 'Hit CTRL+C'; echo aaaaaaa 1>&2; sleep 1; done"},
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
	logging.Instance().Debugf("%s: STOP", logID)

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
