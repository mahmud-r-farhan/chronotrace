//go:build !windows

package main

import (
	"os/exec"
	"syscall"
)

// detachCmd configures cmd so the spawned daemon runs in its own session,
// detached from the GUI's process group and controlling terminal. Closing the
// GUI (or Ctrl+C in a terminal) therefore does not kill the daemon.
func detachCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}
