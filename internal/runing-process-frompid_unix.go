// +build !windows

package internal

type unixProcMedata struct {
	State      rune
	ParentPID  int
	Pgrp       int
	Sid        int
	Executable string
}
