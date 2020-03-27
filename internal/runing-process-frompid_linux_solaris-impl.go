// +build linux solaris

package internal

import (
	"os"

	"github.com/codemodify/systemkit-processes/contracts"
)

func existingUnixProcessByPID(pid int) (contracts.RuningProcess, error) {
	upm, err := fetchProcMedata(pid)
	if err != nil {
		return NewEmptyRuningProcess(), err
	}

	osProcess, err := os.FindProcess(pid)
	if err != nil {
		return NewEmptyRuningProcess(), err
	}

	return NewRuningProcessWithOSProc(
		contracts.ProcessTemplate{
			Executable: upm.Executable,
		},
		osProcess,
	), nil
}
