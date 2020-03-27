// +build linux solaris

package internal

import (
	"os"

	"github.com/codemodify/systemkit-processes/contracts"
)

type unixProcMedata struct {
	State      rune
	ParentPID  int
	Pgrp       int
	Sid        int
	Executable string
}

func existingUnixProcessByPID(pid int) (contracts.RuningProcess, error) {
	upm, err := fetchProcMedata(pid)
	if err != nil {
		return NewRuningProcess(contracts.ProcessTemplate{}), err
	}

	osProcess, err := os.FindProcess(pid)
	if err != nil {
		return NewRuningProcess(contracts.ProcessTemplate{}), err
	}

	return NewRuningProcessWithOSProc(
		contracts.ProcessTemplate{
			Executable: upm.Executable,
		},
		osProcess,
	), nil
}
