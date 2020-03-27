// +build windows

package internal

import (
	"os"
	"syscall"

	"github.com/codemodify/systemkit-processes/contracts"
)

// Windows API functions
var (
	modKernel32                  = syscall.NewLazyDLL("kernel32.dll")
	procCloseHandle              = modKernel32.NewProc("CloseHandle")
	procCreateToolhelp32Snapshot = modKernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32First           = modKernel32.NewProc("Process32FirstW")
	procProcess32Next            = modKernel32.NewProc("Process32NextW")
)

// Some constants from the Windows API
const (
	ERROR_NO_MORE_FILES = 0x12
	MAX_PATH            = 260
)

// PROCESSENTRY32 is the Windows API structure that contains a process's
// information.
type PROCESSENTRY32 struct {
	Size              uint32
	CntUsage          uint32
	ProcessID         uint32
	DefaultHeapID     uintptr
	ModuleID          uint32
	CntThreads        uint32
	ParentProcessID   uint32
	PriorityClassBase int32
	Flags             uint32
	ExeFile           [MAX_PATH]uint16
}

func existingWindowsProcessFromProcEntry(e *PROCESSENTRY32) (contracts.RuningProcess, error) {
	// find where the string ends
	end := 0
	for {
		if e.ExeFile[end] == 0 {
			break
		}
		end++
	}

	// build the object
	osProcess, err := os.FindProcess(int(e.ProcessID))
	if err != nil {
		return NewEmptyRuningProcess(), err
	}

	// FIXME: int(e.ParentProcessID)

	return NewRuningProcessWithOSProc(
		contracts.ProcessTemplate{
			Executable: syscall.UTF16ToString(e.ExeFile[:end]),
		},
		osProcess,
	), nil
}
