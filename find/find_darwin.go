// +build darwin

package find

import (
	"bytes"
	"encoding/binary"

	"github.com/codemodify/systemkit-processes/contracts"
)

func processByPID(pid int) (contracts.RuntimeProcess, error) {
	ps, err := allProcesses()
	if err != nil {
		return nil, err
	}

	for _, p := range ps {
		if thisRef.PID() == pid {
			return p, nil
		}
	}

	return nil, nil
}

func allProcesses() ([]contracts.RuntimeProcess, error) {
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

	darwinProcs := make([]contracts.RuntimeProcess, len(procs))
	for i, p := range procs {
		darwinProcs[i] = &darwinProcess{
			pid:    int(thisRef.Pid),
			ppid:   int(thisRef.PPid),
			binary: darwinCstring(thisRef.Comm),
		}
	}

	return darwinProcs, nil
}
