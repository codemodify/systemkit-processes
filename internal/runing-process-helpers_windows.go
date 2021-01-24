// +build windows

package internal

import (
	"golang.org/x/sys/windows"
)

var procAttrs = &windows.SysProcAttr{
	CreationFlags: windows.CREATE_UNICODE_ENVIRONMENT |
		windows.CREATE_NEW_PROCESS_GROUP |
		windows.CREATE_NEW_CONSOLE |
		windows.CREATE_NO_WINDOW,
}
