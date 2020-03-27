// +build linux solaris

package internal

import (
	"fmt"
	"os"

	"github.com/codemodify/systemkit-processes/contracts"
)

func processByPID(pid int) (contracts.RuningProcess, error) {
	dir := fmt.Sprintf("/proc/%d", pid)
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return NewRuningProcess(contracts.ProcessTemplate{}), nil
		}

		return nil, err
	}

	return existingUnixProcessByPID(pid)
}
