package autostart

import (
	"path/filepath"
	"testing"
)

func TestAutostartManager(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	mgr := New()
	if mgr == nil {
		t.Fatal("New() returned nil Manager")
	}

	enabled, err := mgr.IsEnabled()
	if err != nil {
		t.Fatalf("IsEnabled() failed: %v", err)
	}
	if enabled {
		t.Errorf("expected isEnabled to be false initially")
	}

	dummyPath := filepath.Join(tmpDir, "dummy-daemon")
	// Test Enable/Disable on Linux (systemctl commands may fail if systemd is not present in container, but file operations work)
	err = mgr.Enable(dummyPath)
	if err != nil {
		t.Logf("Enable() returned error (expected without systemd user session): %v", err)
	} else {
		enabled, _ := mgr.IsEnabled()
		if !enabled {
			t.Errorf("expected isEnabled to be true after Enable")
		}

		_ = mgr.Disable()
		enabled, _ = mgr.IsEnabled()
		if enabled {
			t.Errorf("expected isEnabled to be false after Disable")
		}
	}
}
