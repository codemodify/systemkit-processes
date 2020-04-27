// +build windows

package tests

import (
	"sync"
	"testing"

	logging "github.com/codemodify/systemkit-logging"

	"github.com/codemodify/systemkit-processes/contracts"
	procMon "github.com/codemodify/systemkit-processes/monitor"
)

func TestStoppedUnix(t *testing.T) {
	const logID = "TestStoppedUnix"

	logging.Debugf("%s: START", logID)

	monitor := procMon.New()

	wg := sync.WaitGroup{}
	wg.Add(1)

	processTag, _ := monitor.Spawn(contracts.ProcessTemplate{
		Executable: "notepad.exe",
	})
	monitor.GetProcess(processTag).OnStop(func() {
		logging.Debugf("%s: OnStop()", logID)
		wg.Done()
	})

	logging.Infof(
		"%s: pid: %v",
		logID,
		monitor.GetProcess(processTag).Details().ProcessID,
	)

	wg.Wait()
	logging.Debugf("%s: STOP", logID)
}
