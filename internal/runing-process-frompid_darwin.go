// +build darwin

package internal

import (
	"bytes"
	"encoding/binary"
	"os"

	"github.com/codemodify/systemkit-processes/contracts"
)

func processByPID(pid int) (contracts.RuningProcess, error) {
	ps, err := allProcesses()
	if err != nil {
		return NewEmptyRuningProcess(), err
	}

	for _, p := range ps {
		if p.PID() == pid {
			return p, nil
		}
	}

	return NewEmptyRuningProcess(), nil
}

func allProcesses() ([]contracts.RuningProcess, error) {
	buf, err := darwinSyscall()
	if err != nil {
		return nil, err
	}

	procs := make([]*kinfoProc, 0, 50)
	k := 0
	for i := _KINFO_STRUCT_SIZE; i < buf.Len(); i += _KINFO_STRUCT_SIZE {
		proc := &kinfoProc{}
		err = binary.Read(bytes.NewBuffer(buf.Bytes()[k:i]), binary.LittleEndian, proc)
		if err != nil {
			return nil, err
		}

		k = i
		procs = append(procs, proc)
	}

	rps := make([]contracts.RuningProcess, len(procs))
	for i, p := range procs {
		osProcess, err := os.FindProcess(int(p.Pid))
		if err != nil {
			rps[i] = NewEmptyRuningProcess()
		} else {
			rps[i] = NewRuningProcessWithOSProc(
				contracts.ProcessTemplate{
					Executable: darwinCstring(p.Comm),
				},
				osProcess,
			)

			// rps[i].ParentPID() = int(p.PPid)
		}
	}

	return rps, nil
}
