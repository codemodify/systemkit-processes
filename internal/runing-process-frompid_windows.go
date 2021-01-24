// +build windows

package internal

import (
	"path/filepath"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/codemodify/systemkit-processes/contracts"
)

func getAllRuningProcesses() ([]contracts.RuningProcess, error) {
	// list all processes
	var allProcessesPids []int32
	var psSize uint32 = 1024
	const dwordSize uint32 = 4

	for {
		pids := make([]uint32, psSize)
		var read uint32 = 0
		if err := windows.EnumProcesses(pids, &read); err != nil {
			return nil, err
		}

		if uint32(len(pids)) == read { // ps buffer was too small to host every results, retry with a bigger one
			psSize += 1024
			continue
		}

		for _, pid := range pids[:read/dwordSize] {
			allProcessesPids = append(allProcessesPids, int32(pid))
		}

		break
	}

	// build RuningProcess list
	results := []contracts.RuningProcess{}
	for _, pid := range allProcessesPids {
		rp, err := getRuningProcessByPID(int(pid))
		if err != nil {
			continue
		}

		results = append(results, rp)
	}

	return results, nil
}

func getRuntimeProcessByPID(pid int) (contracts.RuntimeProcess, error) {
	handle, err := windows.CreateToolhelp32Snapshot(0x00000002, 0)
	if handle < 0 || err != nil {
		return contracts.RuntimeProcess{
			State: contracts.ProcessStateUnknown,
		}, err
	}

	var processEntry windows.ProcessEntry32
	processEntry.Size = uint32(unsafe.Sizeof(processEntry))

	err = windows.Process32First(handle, &processEntry)
	if err != nil {
		return contracts.RuntimeProcess{
			State: contracts.ProcessStateUnknown,
		}, err
	}

	for {
		if processEntry.ProcessID == uint32(pid) {
			executable := getExecutabe(&processEntry)

			return contracts.RuntimeProcess{
				Executable:       executable,
				ExecutableName:   filepath.Base(executable),
				Args:             []string{},
				WorkingDirectory: "",
				Environment:      []string{},
				ProcessID:        int(processEntry.ProcessID),
				ParentProcessID:  int(processEntry.ParentProcessID),
				UserID:           0,
				GroupID:          0,
				State:            contracts.ProcessStateRunning,
			}, nil
		}

		err = windows.Process32Next(handle, &processEntry)
		if err != nil {
			break
		}
	}

	return contracts.RuntimeProcess{
		State: contracts.ProcessStateUnknown,
	}, nil
}

func getExecutabe(processEntry *windows.ProcessEntry32) string {
	end := 0
	for {
		if processEntry.ExeFile[end] == 0 {
			break
		}
		end++
	}

	return syscall.UTF16ToString(processEntry.ExeFile[:end])
}
