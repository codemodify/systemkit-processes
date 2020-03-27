// +build freebsd

package internal

import (
	"bytes"
	"encoding/binary"
	"os"
	"syscall"
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

// copied from sys/user.h
type Kinfo_proc struct {
	Ki_structsize   int32
	Ki_layout       int32
	Ki_args         int64
	Ki_paddr        int64
	Ki_addr         int64
	Ki_tracep       int64
	Ki_textvp       int64
	Ki_fd           int64
	Ki_vmspace      int64
	Ki_wchan        int64
	Ki_pid          int32
	Ki_ppid         int32
	Ki_pgid         int32
	Ki_tpgid        int32
	Ki_sid          int32
	Ki_tsid         int32
	Ki_jobc         [2]byte
	Ki_spare_short1 [2]byte
	Ki_tdev         int32
	Ki_siglist      [16]byte
	Ki_sigmask      [16]byte
	Ki_sigignore    [16]byte
	Ki_sigcatch     [16]byte
	Ki_uid          int32
	Ki_ruid         int32
	Ki_svuid        int32
	Ki_rgid         int32
	Ki_svgid        int32
	Ki_ngroups      [2]byte
	Ki_spare_short2 [2]byte
	Ki_groups       [64]byte
	Ki_size         int64
	Ki_rssize       int64
	Ki_swrss        int64
	Ki_tsize        int64
	Ki_dsize        int64
	Ki_ssize        int64
	Ki_xstat        [2]byte
	Ki_acflag       [2]byte
	Ki_pctcpu       int32
	Ki_estcpu       int32
	Ki_slptime      int32
	Ki_swtime       int32
	Ki_cow          int32
	Ki_runtime      int64
	Ki_start        [16]byte
	Ki_childtime    [16]byte
	Ki_flag         int64
	Ki_kiflag       int64
	Ki_traceflag    int32
	Ki_stat         [1]byte
	Ki_nice         [1]byte
	Ki_lock         [1]byte
	Ki_rqindex      [1]byte
	Ki_oncpu        [1]byte
	Ki_lastcpu      [1]byte
	Ki_ocomm        [17]byte
	Ki_wmesg        [9]byte
	Ki_login        [18]byte
	Ki_lockname     [9]byte
	Ki_comm         [20]byte
	Ki_emul         [17]byte
	Ki_sparestrings [68]byte
	Ki_spareints    [36]byte
	Ki_cr_flags     int32
	Ki_jid          int32
	Ki_numthreads   int32
	Ki_tid          int32
	Ki_pri          int32
	Ki_rusage       [144]byte
	Ki_rusage_ch    [144]byte
	Ki_pcb          int64
	Ki_kstack       int64
	Ki_udata        int64
	Ki_tdaddr       int64
	Ki_spareptrs    [48]byte
	Ki_spareint64s  [96]byte
	Ki_sflag        int64
	Ki_tdflags      int64
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

func fetchProcMedata(pid int) (unixProcMedata, error) {
	mib := []int32{CTL_KERN, KERN_PROC, KERN_PROC_PID, int32(pid)}

	buf, length, err := call_syscall(mib)
	if err != nil {
		return unixProcMedata{}, err
	}

	proc_k := Kinfo_proc{}
	if length != uint64(unsafe.Sizeof(proc_k)) {
		return unixProcMedata{}, err
	}

	k, err := parse_kinfo_proc(buf)
	if err != nil {
		return unixProcMedata{}, err
	}

	ppid, pgrp, sid, binary := copy_params(&k)

	return unixProcMedata{
		ParentPID:  ppid,
		Pgrp:       pgrp,
		Sid:        sid,
		Executable: binary,
	}, nil
}

func copy_params(k *Kinfo_proc) (int, int, int, string) {
	n := -1
	for i, b := range k.Ki_comm {
		if b == 0 {
			break
		}
		n = i + 1
	}
	comm := string(k.Ki_comm[:n])

	return int(k.Ki_ppid), int(k.Ki_pgid), int(k.Ki_sid), comm
}

func parse_kinfo_proc(buf []byte) (Kinfo_proc, error) {
	var k Kinfo_proc
	br := bytes.NewReader(buf)
	err := binary.Read(br, binary.LittleEndian, &k)
	if err != nil {
		return k, err
	}

	return k, nil
}

func call_syscall(mib []int32) ([]byte, uint64, error) {
	miblen := uint64(len(mib))

	// get required buffer size
	length := uint64(0)
	_, _, err := syscall.RawSyscall6(
		syscall.SYS___SYSCTL,
		uintptr(unsafe.Pointer(&mib[0])),
		uintptr(miblen),
		0,
		uintptr(unsafe.Pointer(&length)),
		0,
		0)
	if err != 0 {
		b := make([]byte, 0)
		return b, length, err
	}
	if length == 0 {
		b := make([]byte, 0)
		return b, length, err
	}
	// get proc info itself
	buf := make([]byte, length)
	_, _, err = syscall.RawSyscall6(
		syscall.SYS___SYSCTL,
		uintptr(unsafe.Pointer(&mib[0])),
		uintptr(miblen),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&length)),
		0,
		0)
	if err != 0 {
		return buf, length, err
	}

	return buf, length, nil
}
