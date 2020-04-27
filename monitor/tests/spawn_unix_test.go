// +build !windows

package tests

import (
	"testing"
	"time"

	logging "github.com/codemodify/systemkit-logging"

	"github.com/codemodify/systemkit-processes/contracts"
	procMon "github.com/codemodify/systemkit-processes/monitor"
)

func TestSpawnUnix(t *testing.T) {
	const logID = "TestSpawnUnix"

	logging.Debugf("%s: START", logID)

	monitor := procMon.New()

	processTag, _ := monitor.Spawn(contracts.ProcessTemplate{
		Executable: "sh",
		Args:       []string{"-c", "while :; do echo 'Hit CTRL+C'; echo aaaaaaa 1>&2; sleep 1; done"},
	})
	monitor.GetProcess(processTag).OnStdOut(func(data []byte) {
		logging.Debugf("%s: OnStdOut: %v", logID, string(data))
	})
	monitor.GetProcess(processTag).OnStdErr(func(data []byte) {
		logging.Debugf("%s: OnStdErr: %v", logID, string(data))
	})

	logging.Infof(
		"%s: IsRunning: %v, ExitCode: %v, StartedAt: %v, StoppedAt: %v",
		logID,
		monitor.GetProcess(processTag).IsRunning(),
		monitor.GetProcess(processTag).ExitCode(),
		monitor.GetProcess(processTag).StartedAt(),
		monitor.GetProcess(processTag).StoppedAt(),
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
				logging.Debugf("%s: Tick at, %v", logID, t)
			}
		}
	}()
	time.Sleep(5 * time.Second)
	ticker.Stop()
	done <- true

	// STOP
	logging.Debugf("%s: STOP", logID)

	monitor.Stop(processTag)

	logging.Infof(
		"%s: IsRunning: %v, ExitCode: %v, StartedAt: %v, StoppedAt: %v",
		logID,
		monitor.GetProcess(processTag).IsRunning(),
		monitor.GetProcess(processTag).ExitCode(),
		monitor.GetProcess(processTag).StartedAt(),
		monitor.GetProcess(processTag).StoppedAt(),
	)
}
