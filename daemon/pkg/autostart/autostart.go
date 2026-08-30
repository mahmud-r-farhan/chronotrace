// Package autostart provides native OS autostart registration for the daemon.
// The daemon registers itself to auto-launch on user login.
package autostart

// Manager can enable or disable the daemon's autostart entry.
type Manager interface {
	Enable(executablePath string) error
	Disable() error
	IsEnabled() (bool, error)
}

// New returns the platform-appropriate autostart Manager.
// Implementation is in platform-specific files.
func New() Manager {
	return newPlatformManager()
}
