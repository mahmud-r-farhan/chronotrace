// Package tracker defines the interface for active window tracking
// and provides platform-specific implementations.
package tracker

// ActiveWindowInfo holds information about the currently focused window.
type ActiveWindowInfo struct {
	AppName     string // Human-readable application name (e.g. "Google Chrome")
	WindowTitle string // Full window title text
	ExePath     string // Full path to the process executable
	PID         uint32 // Process ID
}

// Tracker is implemented per platform to retrieve the current foreground window.
type Tracker interface {
	GetActiveWindow() (*ActiveWindowInfo, error)
}

// New returns the platform-appropriate Tracker implementation.
// This function is defined in the platform-specific files (tracker_windows.go, etc.)
func New() Tracker {
	return newPlatformTracker()
}
