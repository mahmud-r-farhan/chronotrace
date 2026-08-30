//go:build !windows

package main

import (
	"os"
	"syscall"
)

// isProcessAlive returns true if a process with the given PID is running.
// On Unix systems, sending signal 0 checks for process existence without
// actually delivering a signal.
func isProcessAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// Signal 0 checks existence without delivering a signal.
	err = proc.Signal(syscall.Signal(0))
	return err == nil
}
