package find

import (
	"github.com/codemodify/systemkit-processes/contracts"
	"github.com/codemodify/systemkit-processes/internal"
)

// ProcessByPID - finds process by PID
func ProcessByPID(pid int) (contracts.RuningProcess, error) {
	return internal.ProcessByPID(pid)
}

// AllProcesses - returns all processes
func AllProcesses() ([]contracts.RuningProcess, error) {
	return internal.AllProcesses()
}
