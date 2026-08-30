//go:build windows

package autostart

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows/registry"
)

// Compile-time check that the package is available.
// golang.org/x/sys/windows/registry is part of golang.org/x/sys.

const (
	regKeyPath  = `Software\Microsoft\Windows\CurrentVersion\Run`
	regValueName = "ChronoTraceDaemon"
)

type windowsManager struct{}

func newPlatformManager() Manager {
	return &windowsManager{}
}

// Enable registers the daemon executable in HKCU Run registry key.
func (m *windowsManager) Enable(executablePath string) error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, regKeyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("autostart: open registry key: %w", err)
	}
	defer key.Close()

	// Quote path in case it contains spaces.
	value := fmt.Sprintf(`"%s"`, executablePath)
	if err := key.SetStringValue(regValueName, value); err != nil {
		return fmt.Errorf("autostart: set registry value: %w", err)
	}
	return nil
}

// Disable removes the daemon entry from the Run registry key.
func (m *windowsManager) Disable() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, regKeyPath, registry.SET_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return nil // Already not set.
		}
		return fmt.Errorf("autostart: open registry key: %w", err)
	}
	defer key.Close()

	if err := key.DeleteValue(regValueName); err != nil && err != registry.ErrNotExist {
		return fmt.Errorf("autostart: delete registry value: %w", err)
	}
	return nil
}

// IsEnabled returns true if the registry entry exists.
func (m *windowsManager) IsEnabled() (bool, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, regKeyPath, registry.QUERY_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return false, nil
		}
		return false, fmt.Errorf("autostart: open registry key: %w", err)
	}
	defer key.Close()

	_, _, err = key.GetStringValue(regValueName)
	if err == registry.ErrNotExist {
		return false, nil
	}
	return err == nil, err
}

// SelfExecutablePath returns the absolute path to the running binary.
func SelfExecutablePath() (string, error) {
	return os.Executable()
}
