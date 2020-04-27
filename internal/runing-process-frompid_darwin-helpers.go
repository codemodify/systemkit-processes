// +build darwin cgo

package internal

// #include <libproc.h>
// #include <stdlib.h>
import "C"

import (
	"syscall"
	"unsafe"
)

// See https://opensource.apple.com/source/xnu/xnu-2782.40.9/libsyscall/wrappers/libproc/libproc.c
const pathBufferSize = C.PROC_PIDPATHINFO_MAXSIZE

// See https://opensource.apple.com/source/xnu/xnu-1699.24.23/bsd/sys/proc_internal.h
const pidsBufferSize = 99999 // #define Pid_MAX 99999

func listAllPids() ([]int, error) {
	buffer := make([]uint32, pidsBufferSize)

	_, err := C.proc_listallpids(unsafe.Pointer(&buffer[0]), C.int(pidsBufferSize))
	if err != nil {
		return nil, err
	}

	pids := []int{}
	for _, pid := range buffer {
		if pid != 0 {
			pids = append(pids, int(pid))
		}
	}

	return pids, nil
}

func fromPidGetProcPath(pid int) (string, error) {
	// Allocate in the C heap a string (char* terminated with `/0`) of size `pathBufferSize`
	// Make sure that we free that memory that gets allocated in C (see the `defer` below)
	buffer := C.CString(string(make([]byte, pathBufferSize)))
	defer C.free(unsafe.Pointer(buffer))

	// Call libproc -> proc_pidpath
	ret, err := C.proc_pidpath(C.int(pid), unsafe.Pointer(buffer), pathBufferSize)
	if ret <= 0 {
		return "", err
	}

	// Convert the C string back to a Go string.
	path := C.GoString(buffer)

	return path, nil
}

func fromPidGetProcName(pid int) (string, error) {
	// Allocate in the C heap a string (char* terminated with `/0`) of size `pathBufferSize`
	// Make sure that we free that memory that gets allocated in C (see the `defer` below)
	buffer := C.CString(string(make([]byte, pathBufferSize)))
	defer C.free(unsafe.Pointer(buffer))

	// Call libproc -> proc_name
	ret, err := C.proc_name(C.int(pid), unsafe.Pointer(buffer), pathBufferSize)
	if ret <= 0 {
		return "", err
	}

	// Convert the C string back to a Go string.
	path := C.GoString(buffer)

	return path, nil
}

func fromPidGetProcInfo(pid int, info *C.struct_proc_taskallinfo) error {
	size := C.int(unsafe.Sizeof(*info))
	ptr := unsafe.Pointer(info)

	n := C.proc_pidinfo(C.int(pid), C.PROC_PIDTASKALLINFO, 0, ptr, size)
	if n != size {
		return syscall.ENOMEM
	}

	return nil
}
