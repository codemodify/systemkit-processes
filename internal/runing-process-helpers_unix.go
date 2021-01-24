// +build !windows

package internal

import (
	"golang.org/x/sys/unix"
)

var procAttrs = &unix.SysProcAttr{}
