package internal

import (
	"os"

	"github.com/codemodify/systemkit-processes/contracts"
)

// GetRuningProcessByPID - finds process by PID
func GetRuningProcessByPID(pid int) (contracts.RuningProcess, error) {
	return getRuningProcessByPID(pid)
}

// GetAllRuningProcesses - returns all processes
func GetAllRuningProcesses() ([]contracts.RuningProcess, error) {
	return getAllRuningProcesses()
}

func getRuningProcessByPID(pid int) (contracts.RuningProcess, error) {
	rp, err := getRuntimeProcessByPID(pid)
	if err != nil {
		return NewEmptyRuningProcess(), contracts.ErrProcessDoesNotExist
	}

	osProcess, err := os.FindProcess(pid)
	if err != nil {
		return NewEmptyRuningProcess(), contracts.ErrProcessDoesNotExist
	}

	return NewRuningProcessWithOSProc(
		contracts.ProcessTemplate{
			Executable:       rp.Executable,
			Args:             rp.Args,
			WorkingDirectory: rp.WorkingDirectory,
			Environment:      rp.Environment,
		},
		osProcess,
	), nil
}
