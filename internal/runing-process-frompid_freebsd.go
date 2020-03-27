// +build freebsd

package internal

import (
	"unsafe"

	"github.com/codemodify/systemkit-processes/contracts"
)

func processByPID(pid int) (contracts.RuningProcess, error) {
	mib := []int32{CTL_KERN, KERN_PROC, KERN_PROC_PATHNAME, int32(pid)}

	_, _, err := call_syscall(mib)
	if err != nil {
		return nil, err
	}

	return existingUnixProcessByPID(pid)
}

func allProcesses() ([]contracts.RuningProcess, error) {
	results := make([]contracts.RuningProcess, 0, 50)

	mib := []int32{CTL_KERN, KERN_PROC, KERN_PROC_PROC, 0}
	buf, length, err := call_syscall(mib)
	if err != nil {
		return results, err
	}

	// get kinfo_proc size
	k := Kinfo_proc{}
	procinfo_len := int(unsafe.Sizeof(k))
	count := int(length / uint64(procinfo_len))

	// parse buf to procs
	for i := 0; i < count; i++ {
		b := buf[i*procinfo_len : i*procinfo_len+procinfo_len]
		k, err := parse_kinfo_proc(b)
		if err != nil {
			continue
		}
		p, err := existingUnixProcessByPID(int(k.Ki_pid))
		if err != nil {
			continue
		}
		// p.ppid, p.pgrp, p.sid, p.binary = copy_params(&k)

		results = append(results, p)
	}

	return results, nil
}
