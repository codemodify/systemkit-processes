// +build linux

package internal

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/codemodify/systemkit-processes/contracts"
)

func getAllRuningProcesses() ([]contracts.RuningProcess, error) {
	d, err := os.Open("/proc")
	if err != nil {
		return nil, err
	}
	defer d.Close()

	results := []contracts.RuningProcess{}
	for {
		fis, err := d.Readdir(10)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for _, fi := range fis {
			// We only care about directories, since all pids are dirs
			if !fi.IsDir() {
				continue
			}

			// We only care if the name starts with a numeric
			name := fi.Name()
			if name[0] < '0' || name[0] > '9' {
				continue
			}

			// From this point forward, any errors we just ignore, because
			// it might simply be that the process doesn't exist anymore.
			pid, err := strconv.ParseInt(name, 10, 0)
			if err != nil {
				continue
			}

			p, err := getRuningProcessByPID(int(pid))
			if err != nil {
				continue
			}

			results = append(results, p)
		}
	}

	return results, nil
}

func getRuntimeProcessByPID(pid int) (contracts.RuntimeProcess, error) {
	folder := fmt.Sprintf("/proc/%d", pid)

	// 1 - check folder exists
	_, err := os.Stat(folder)
	if err != nil {
		if os.IsNotExist(err) {
			return contracts.RuntimeProcess{
				State: contracts.ProcessStateNonExistent,
			}, contracts.ErrProcessDoesNotExist
		}

		return contracts.RuntimeProcess{
			State: contracts.ProcessStateUnknown,
		}, err
	}

	//
	// /proc/%d/*
	// 		cwd			-> sym link to the working dir
	//		environ		-> env vars
	//		status 		-> Name, Pid, PPid, Uid, Gid
	//		cmdline		-> full path with args
	//
	// 		comm		-> executable name
	//		loginuid 	-> ID of the running-as user
	//

	procMedata := contracts.RuntimeProcess{
		Args:        []string{},
		Environment: []string{},
		State:       contracts.ProcessStateRunning,
	}

	// 2 - read cwd
	fi, err := os.Lstat(path.Join(folder, "cwd"))
	if err == nil && fi != nil && fi.Mode()&os.ModeSymlink != 0 {
		procMedata.WorkingDirectory, _ = os.Readlink(path.Join(folder, "cwd"))
	}

	// 3 - read environ
	data, _ := ioutil.ReadFile(path.Join(folder, "environ"))
	lines := strings.Split(string(data), "\x00")
	for _, line := range lines {
		procMedata.Environment = append(procMedata.Environment, line)
	}

	// 4 - read status
	data, _ = ioutil.ReadFile(path.Join(folder, "status"))
	lines = strings.Split(string(data), "\n")
	for _, line := range lines {
		props := strings.Split(line, ":")
		if len(props) > 1 {
			key := strings.TrimSpace(strings.ToLower(props[0]))
			val := strings.TrimSpace(strings.ToLower(props[1]))
			switch key {
			case "name":
				procMedata.ExecutableName = val
			case "state":
				states := strings.Split(val, " ")
				if len(states) > 0 {
					switch states[0] {
					case "d":
						procMedata.State = contracts.ProcessStateWaitingIO
					case "r":
						procMedata.State = contracts.ProcessStateRunning
					case "s":
						procMedata.State = contracts.ProcessStateWaitingEvent
					case "t":
						procMedata.State = contracts.ProcessStateTraced
					case "w":
						procMedata.State = contracts.ProcessStatePaging
					case "x":
						procMedata.State = contracts.ProcessStateDead
					case "z":
						procMedata.State = contracts.ProcessStateObsolete
					}
				}
			case "pid":
				fetchedPid, _ := strconv.Atoi(val)
				procMedata.ProcessID = fetchedPid
			case "ppid":
				ppid, _ := strconv.Atoi(val)
				procMedata.ParentProcessID = ppid
			case "uid":
				uids := strings.Split(val, "\t")
				if len(uids) > 0 {
					uid, _ := strconv.Atoi(uids[0])
					procMedata.UserID = uid
				}
			case "gid":
				gids := strings.Split(val, "\t")
				if len(gids) > 0 {
					gid, _ := strconv.Atoi(gids[0])
					procMedata.GroupID = gid
				}
			}
		}
	}

	// 4 - read cmdline
	data, _ = ioutil.ReadFile(path.Join(folder, "cmdline"))
	lines = strings.Split(string(data), "\x00")
	for index, line := range lines {
		if index == 0 {
			procMedata.Executable = line
		} else {
			trimmedLine := strings.TrimSpace(line)
			if len(trimmedLine) > 0 {
				procMedata.Args = append(procMedata.Args, trimmedLine)
			}
		}
	}

	return procMedata, err
}
