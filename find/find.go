package find

import (
	"github.com/codemodify/systemkit-processes/contracts"
	"github.com/codemodify/systemkit-processes/internal"
)

// GetRuningProcessByPID - finds process by PID
func GetRuningProcessByPID(pid int) (contracts.RuningProcess, error) {
	return internal.GetRuningProcessByPID(pid)
}

// GetAllRuningProcesses - returns all processes
func GetAllRuningProcesses() ([]contracts.RuningProcess, error) {
	return internal.GetAllRuningProcesses()
}
