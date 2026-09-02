package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestAppInstantiation(t *testing.T) {
	app := NewApp()
	if app == nil {
		t.Fatal("NewApp() returned nil")
	}
}

func TestFindDaemonExecutable(t *testing.T) {
	// findDaemonExecutable checks various locations for daemon executable.
	// In test environment without build binaries, it returns an error.
	_, err := findDaemonExecutable()
	if err != nil {
		t.Logf("findDaemonExecutable error (expected when not built): %v", err)
	}
}

func TestAppGetStatusOffline(t *testing.T) {
	app := NewApp()
	status := app.GetStatus()
	if connected, ok := status["connected"].(bool); !ok || connected {
		t.Errorf("expected connected=false when daemon not running, got %+v", status)
	}
}

func TestFindDaemonExecutablePath(t *testing.T) {
	tmpDir := t.TempDir()
	daemonName := "chronotrace-daemon"
	if runtime.GOOS == "windows" {
		daemonName = "chronotrace-daemon.exe"
	}
	exePath := filepath.Join(tmpDir, daemonName)
	if err := os.WriteFile(exePath, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatalf("failed to write dummy executable: %v", err)
	}

	t.Setenv("PATH", tmpDir)

	path, err := findDaemonExecutable()
	if err != nil {
		t.Fatalf("findDaemonExecutable failed when in PATH: %v", err)
	}
	if path != exePath {
		t.Errorf("expected path %s, got %s", exePath, path)
	}
}
