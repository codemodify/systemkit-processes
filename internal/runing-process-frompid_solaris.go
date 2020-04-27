// +build solaris

package internal

import "github.com/codemodify/systemkit-processes/contracts"

func getAllRuningProcesses() ([]contracts.RuningProcess, error) {
	// FIXME:
}

func getRuntimeProcessByPID(pid int) (contracts.RuntimeProcess, error) {
	// FIXME:
}
