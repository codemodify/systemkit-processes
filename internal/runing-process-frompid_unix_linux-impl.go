// +build linux

package internal

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func fetchProcMedata(pid int) (unixProcMedata, error) {
	statAsBytes, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return unixProcMedata{}, err
	}
	statAsString := string(statAsBytes)

	binStart := strings.IndexRune(statAsString, '(') + 1
	binEnd := strings.IndexRune(statAsString[binStart:], ')')

	result := unixProcMedata{
		Executable: statAsString[binStart : binStart+binEnd],
	}

	// Move past the image name and start parsing the rest
	statAsString = statAsString[binStart+binEnd+2:]
	_, err = fmt.Sscanf(statAsString,
		"%c %d %d %d",
		&result.State,
		&result.ParentPID,
		&result.Pgrp,
		&result.Sid,
	)

	return result, err
}
