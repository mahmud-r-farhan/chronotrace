//go:build !windows

package main

import (
	"errors"
	"os"
	"syscall"
)

// isProcessAlive returns true if a process with the given PID is running.
// On Unix systems, sending signal 0 checks for process existence without
// actually delivering a signal. EPERM means the process exists but belongs
// to another user, which also counts as alive.
func isProcessAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// Signal 0 checks existence without delivering a signal.
	err = proc.Signal(syscall.Signal(0))
	return err == nil || errors.Is(err, syscall.EPERM)
}
