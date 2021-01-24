// +build !windows

package internal

func processKillHelper(pid int) {}

func sendCtrlC(pid int) error {
	return nil
}
