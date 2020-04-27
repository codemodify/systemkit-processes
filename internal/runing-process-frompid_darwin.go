// +build darwin cgo

package internal

// #include <libproc.h>
// #include <sys/sysctl.h>
import "C"

import (
	"github.com/codemodify/systemkit-processes/contracts"
)

func getAllRuningProcesses() ([]contracts.RuningProcess, error) {
	pids, err := listAllPids()
	if err != nil {
		return []contracts.RuningProcess{}, err
	}

	results := []contracts.RuningProcess{}

	for _, pid := range pids {
		p, err := getRuningProcessByPID(int(pid))
		if err != nil {
			continue
		}

		results = append(results, p)
	}

	return results, nil
}

func getRuntimeProcessByPID(pid int) (contracts.RuntimeProcess, error) {
	// https://fergofrog.com/code/cbowser/xnu/bsd/sys/proc_info.h.html#proc_bsdinfo
	info := C.struct_proc_taskallinfo{}
	if err := fromPidGetProcInfo(pid, &info); err != nil {
		return contracts.RuntimeProcess{
			State: contracts.ProcessStateUnknown,
		}, err
	}

	result := contracts.RuntimeProcess{}

	path, _ := fromPidGetProcPath(pid)
	result.Executable = path

	name, _ := fromPidGetProcName(pid)
	result.ExecutableName = name

	// result.Args
	// result.WorkingDirectory
	// result.Environment

	result.ProcessID = pid
	result.ParentProcessID = int(info.pbsd.pbi_ppid)
	result.UserID = int(info.pbsd.pbi_uid)
	result.GroupID = int(info.pbsd.pbi_gid)

	switch info.pbsd.pbi_status {
	case C.SIDL:
		result.State = contracts.ProcessStateWaitingEvent
	case C.SRUN:
		result.State = contracts.ProcessStateRunning
	case C.SSLEEP:
		result.State = contracts.ProcessStateWaitingEvent
	case C.SSTOP:
		result.State = contracts.ProcessStateDead
	case C.SZOMB:
		result.State = contracts.ProcessStateObsolete
	default:
		result.State = contracts.ProcessStateUnknown
	}

	return result, nil
}
