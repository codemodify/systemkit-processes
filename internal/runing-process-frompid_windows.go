// +build windows

package internal

import (
	"fmt"
	"strings"
	"time"

	"github.com/StackExchange/wmi"
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
	var processes []Win32_Process

	wmiQuery := wmi.CreateQuery(&processes, fmt.Sprintf("WHERE ProcessID = %d", pid))
	if err := wmi.Query(wmiQuery, &processes); err != nil {
		return contracts.RuntimeProcess{
			State: contracts.ProcessStateUnknown,
		}, err
	}

	if len(processes) == 0 {
		return contracts.RuntimeProcess{
			State: contracts.ProcessStateNonExistent,
		}, contracts.ErrProcessDoesNotExist
	}

	args := *processes[0].CommandLine
	args = strings.Replace(args, processes[0].Name, "", 1)

	cwd := *processes[0].ExecutablePath
	cwd = strings.Replace(cwd, processes[0].Name, "", 1)

	return contracts.RuntimeProcess{
		Executable:       *processes[0].ExecutablePath,
		ExecutableName:   processes[0].Name,
		Args:             strings.Split(args, " "),
		WorkingDirectory: cwd,
		Environment:      []string{},
		ProcessID:        int(processes[0].ProcessID),
		ParentProcessID:  int(processes[0].ParentProcessID),
		UserID:           0,
		GroupID:          0,
		State:            contracts.ProcessStateRunning,
	}, nil
}

type Win32_Process struct {
	Name            string
	ExecutablePath  *string
	CommandLine     *string
	ProcessID       uint32
	Status          *string
	ParentProcessID uint32

	Priority              uint32
	CreationDate          *time.Time
	ThreadCount           uint32
	ReadOperationCount    uint64
	ReadTransferCount     uint64
	WriteOperationCount   uint64
	WriteTransferCount    uint64
	CSCreationClassName   string
	CSName                string
	Caption               *string
	CreationClassName     string
	Description           *string
	ExecutionState        *uint16
	HandleCount           uint32
	KernelModeTime        uint64
	MaximumWorkingSetSize *uint32
	MinimumWorkingSetSize *uint32
	OSCreationClassName   string
	OSName                string
	OtherOperationCount   uint64
	OtherTransferCount    uint64
	PageFaults            uint32
	PageFileUsage         uint32
	PeakPageFileUsage     uint32
	PeakVirtualSize       uint64
	PeakWorkingSetSize    uint32
	PrivatePageCount      uint64
	TerminationDate       *time.Time
	UserModeTime          uint64
	WorkingSetSize        uint64
}
