package internal

import (
	"github.com/codemodify/systemkit-processes/contracts"
)

// ProcessByPID - finds process by PID
func ProcessByPID(pid int) (contracts.RuningProcess, error) {
	return processByPID(pid)
}
