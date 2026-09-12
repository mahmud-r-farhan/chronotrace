//go:build windows

package main

import "golang.org/x/sys/windows"

// stillActive is the exit code reported by GetExitCodeProcess for a process
// that has not exited yet (Win32 STILL_ACTIVE). x/sys/windows does not
// export this constant, so define it locally.
const stillActive = 259

// isProcessAlive returns true if a process with the given PID is running.
func isProcessAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	if err != nil {
		return false
	}
	defer windows.CloseHandle(h)

	var code uint32
	if err := windows.GetExitCodeProcess(h, &code); err != nil {
		return false
	}
	return code == stillActive
}
