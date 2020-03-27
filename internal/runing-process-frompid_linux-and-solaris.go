// +build linux solaris

package internal

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/codemodify/systemkit-processes/contracts"
)

func processByPID(pid int) (contracts.RuningProcess, error) {
	dir := fmt.Sprintf("/proc/%d", pid)
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return NewEmptyRuningProcess(), nil
		}

		return nil, err
	}

	return existingUnixProcessByPID(pid)
}

func allProcesses() ([]contracts.RuningProcess, error) {
	d, err := os.Open("/proc")
	if err != nil {
		return nil, err
	}
	defer d.Close()

	results := []contracts.RuningProcess{}
	for {
		fis, err := d.Readdir(10)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for _, fi := range fis {
			// We only care about directories, since all pids are dirs
			if !fi.IsDir() {
				continue
			}

			// We only care if the name starts with a numeric
			name := fi.Name()
			if name[0] < '0' || name[0] > '9' {
				continue
			}

			// From this point forward, any errors we just ignore, because
			// it might simply be that the process doesn't exist anymore.
			pid, err := strconv.ParseInt(name, 10, 0)
			if err != nil {
				continue
			}

			p, err := processByPID(int(pid))
			if err != nil {
				continue
			}

			results = append(results, p)
		}
	}

	return results, nil
}

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
