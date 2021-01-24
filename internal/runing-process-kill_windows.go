// +build windows

package internal

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func processKillHelper(pid int) {
	sendCtrlBreak(pid)
	sendWMQuit(pid)

	h, err := openProcessHandle(pid)
	if err != nil {
		return
	}
	defer syscall.CloseHandle(h)

	const exitCode = 1
	syscall.TerminateProcess(h, uint32(exitCode))
}

func openProcessHandle(pid int) (syscall.Handle, error) {
	const da = syscall.STANDARD_RIGHTS_READ |
		syscall.PROCESS_QUERY_INFORMATION |
		syscall.SYNCHRONIZE |
		syscall.PROCESS_TERMINATE
	return syscall.OpenProcess(da, false, uint32(pid))
}

// Used to nicely quit console applications
func sendCtrlBreak(pid int) error {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procGenerateConsoleCtrlEvent := kernel32.NewProc("GenerateConsoleCtrlEvent")
	r, _, _ := procGenerateConsoleCtrlEvent.Call(syscall.CTRL_BREAK_EVENT, uintptr(pid))
	if r == 0 {
		return fmt.Errorf("Error calling GenerateConsoleCtrlEvent")
	}
	return nil
}

// Used to nicely quit gui applications
func sendWMQuit(pid int) error {
	user32 := syscall.NewLazyDLL("user32.dll")
	procEnumWindows := user32.NewProc("EnumWindows")
	procGetWindowThreadProcessId := user32.NewProc("GetWindowThreadProcessId")
	procPostMessage := user32.NewProc("PostMessageW")

	// FIXME: Do we need to unregister the callback?
	quitCallback := syscall.NewCallback(func(hwnd syscall.Handle, lparam uintptr) uintptr {
		pid := int(lparam)
		// Does the window belong to our PID?
		var windowPID int
		procGetWindowThreadProcessId.Call(uintptr(hwnd),
			uintptr(unsafe.Pointer(&windowPID)))
		if windowPID == pid {
			const WM_CLOSE = 16
			procPostMessage.Call(uintptr(hwnd), uintptr(WM_CLOSE), 0, 0)
		}
		return 1 // continue enumeration
	})
	ret, _, _ := procEnumWindows.Call(quitCallback, uintptr(pid))
	if ret == 0 {
		return fmt.Errorf("Error called EnumWindows")
	}
	return nil
}

func sendCtrlC(pid int) error {
	dll, err := windows.LoadDLL("kernel32.dll")
	if err != nil {
		return err
	}

	defer dll.Release()

	f, err := dll.FindProc("AttachConsole")
	if err != nil {
		return err
	}

	r1, _, err := f.Call(uintptr(pid))
	if r1 == 0 && err != syscall.ERROR_ACCESS_DENIED {
		return err
	}

	f, err = dll.FindProc("SetConsoleCtrlHandler")
	if err != nil {
		return err
	}

	r1, _, err = f.Call(0, 1)
	if r1 == 0 {
		return err
	}

	f, err = dll.FindProc("GenerateConsoleCtrlEvent")
	if err != nil {
		return err
	}

	r1, _, err = f.Call(windows.CTRL_BREAK_EVENT, uintptr(pid))
	if r1 == 0 {
		return err
	}

	return nil
}
