// +build freebsd

package find

import (
	"unsafe"

	"github.com/codemodify/systemkit-processes/contracts"
)

// copied from sys/sysctl.h
const (
	CTL_KERN           = 1  // "high kernel": proc, limits
	KERN_PROC          = 14 // struct: process entries
	KERN_PROC_PID      = 1  // by process id
	KERN_PROC_PROC     = 8  // only return procs
	KERN_PROC_PATHNAME = 12 // path to executable
)

func processByPID(pid int) (contracts.RuntimeProcess, error) {
	mib := []int32{CTL_KERN, KERN_PROC, KERN_PROC_PATHNAME, int32(pid)}

	_, _, err := call_syscall(mib)
	if err != nil {
		return nil, err
	}

	return newUnixProcess(pid)
}

func allProcesses() ([]contracts.RuntimeProcess, error) {
	results := make([]contracts.RuntimeProcess, 0, 50)

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
		p, err := newUnixProcess(int(k.Ki_pid))
		if err != nil {
			continue
		}
		thisRef.ppid, thisRef.pgrp, thisRef.sid, thisRef.binary = copy_params(&k)

		results = append(results, p)
	}

	return results, nil
}
