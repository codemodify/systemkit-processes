// +build windows

package internal

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/codemodify/systemkit-processes/contracts"
)

func processByPID(pid int) (contracts.RuningProcess, error) {
	ps, err := allProcesses()
	if err != nil {
		return nil, err
	}

	for _, p := range ps {
		if p.PID() == pid {
			return p, nil
		}
	}

	return NewEmptyRuningProcess(), nil
}

func allProcesses() ([]contracts.RuningProcess, error) {
	handle, _, _ := procCreateToolhelp32Snapshot.Call(
		0x00000002,
		0)
	if handle < 0 {
		return nil, syscall.GetLastError()
	}
	defer procCloseHandle.Call(handle)

	var entry PROCESSENTRY32
	entry.Size = uint32(unsafe.Sizeof(entry))
	ret, _, _ := procProcess32First.Call(handle, uintptr(unsafe.Pointer(&entry)))
	if ret == 0 {
		return nil, fmt.Errorf("Error retrieving process info")
	}

	results := make([]contracts.RuningProcess, 0, 50)
	for {
		rp, err := existingWindowsProcessFromProcEntry(&entry)
		if err != nil {
			results = append(results, NewEmptyRuningProcess())
		} else {
			results = append(results, rp)
		}

		ret, _, _ := procProcess32Next.Call(handle, uintptr(unsafe.Pointer(&entry)))
		if ret == 0 {
			break
		}
	}

	return results, nil
}
