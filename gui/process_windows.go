//go:build windows

package main

import (
	"os/exec"
	"syscall"
)

// detachCmd configures cmd so the spawned daemon runs in a new process group,
// independent of the GUI's console lifecycle. The daemon is built as a
// windowed binary, so it keeps running after the GUI exits.
func detachCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}
