package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/mahmud-r-farhan/chronotrace/pkg/autostart"
	"github.com/mahmud-r-farhan/chronotrace/pkg/ipc"
	"github.com/mahmud-r-farhan/chronotrace/pkg/storage"
)

// App is the Go backend bound to the Wails frontend via JS bindings.
// All public methods are automatically exposed to the frontend as JS functions.
type App struct {
	ctx    context.Context
	client *ipc.Client
}

// NewApp creates a new App instance.
func NewApp() *App {
	return &App{
		client: ipc.NewClient(ipc.DefaultAddr),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	log.Println("[gui] started")

	// On startup, ensure background daemon is running.
	go func() {
		time.Sleep(500 * time.Millisecond)
		if _, err := a.client.GetStatus(); err != nil {
			log.Println("[gui] daemon not detected on startup, attempting auto-spawn...")
			if err := a.StartDaemon(); err != nil {
				log.Printf("[gui] failed to auto-spawn daemon: %v", err)
			}
		}
	}()
}

func (a *App) shutdown(ctx context.Context) {
	log.Println("[gui] shutdown")
}

// findDaemonExecutable locates the chronotrace-daemon binary.
func findDaemonExecutable() (string, error) {
	daemonName := "chronotrace-daemon"
	if runtime.GOOS == "windows" {
		daemonName = "chronotrace-daemon.exe"
	}

	// 1. Check directory of current GUI executable
	if exePath, err := os.Executable(); err == nil {
		dir := filepath.Dir(exePath)
		candidate := filepath.Join(dir, daemonName)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
		// Also check parent build folder
		candidate = filepath.Join(dir, "..", daemonName)
		if _, err := os.Stat(candidate); err == nil {
			return filepath.Clean(candidate), nil
		}
	}

	// 2. Check current working directory and build/ folder
	if candidate, err := filepath.Abs(daemonName); err == nil {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	if candidate, err := filepath.Abs(filepath.Join("build", daemonName)); err == nil {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	// 3. Check standard installation directory (%LOCALAPPDATA%\Programs\ChronoTrace)
	if runtime.GOOS == "windows" {
		if localApp := os.Getenv("LOCALAPPDATA"); localApp != "" {
			installed := filepath.Join(localApp, "Programs", "ChronoTrace", daemonName)
			if _, err := os.Stat(installed); err == nil {
				return installed, nil
			}
		}
	}

	// 4. Check PATH
	if path, err := exec.LookPath(daemonName); err == nil {
		return path, nil
	}

	return "", fmt.Errorf("chronotrace-daemon binary not found")
}

// StartDaemon starts the background tracking service if not running.
func (a *App) StartDaemon() error {
	daemonPath, err := findDaemonExecutable()
	if err != nil {
		return err
	}

	cmd := exec.Command(daemonPath)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start daemon: %w", err)
	}

	// Wait up to 3s for daemon to become ready
	for i := 0; i < 6; i++ {
		time.Sleep(500 * time.Millisecond)
		if _, err := a.client.GetStatus(); err == nil {
			return nil
		}
	}

	return nil
}

// --- Status ---

// GetStatus returns the daemon connection status.
func (a *App) GetStatus() map[string]interface{} {
	s, err := a.client.GetStatus()
	if err != nil {
		return map[string]interface{}{
			"connected": false,
			"error":     err.Error(),
		}
	}
	return map[string]interface{}{
		"connected":      true,
		"version":        s.Version,
		"uptime_seconds": s.UptimeSec,
		"addr":           s.Addr,
	}
}

// --- Usage data ---

// GetUsageToday returns today's per-app usage sorted by total time.
func (a *App) GetUsageToday() ([]storage.AppUsage, error) {
	return a.client.GetUsageToday()
}

// GetUsageWeek returns the past 7 days of per-app usage.
func (a *App) GetUsageWeek() ([]storage.AppUsage, error) {
	return a.client.GetUsageWeek()
}

// GetUsageMonth returns the past 30 days of per-app usage.
func (a *App) GetUsageMonth() ([]storage.AppUsage, error) {
	return a.client.GetUsageMonth()
}

// GetTimeline returns the hourly timeline for the given date (YYYY-MM-DD).
func (a *App) GetTimeline(date string) ([]storage.TimelineSlot, error) {
	return a.client.GetTimeline(date)
}

// GetSummary returns a day summary for the given date.
func (a *App) GetSummary(date string) (*storage.DaySummary, error) {
	return a.client.GetSummary(date)
}

// --- Autostart ---

// EnableAutostart registers the daemon in OS autostart.
func (a *App) EnableAutostart() error {
	mgr := autostart.New()
	daemonPath, err := findDaemonExecutable()
	if err != nil {
		daemonPath = "chronotrace-daemon"
	}
	return mgr.Enable(daemonPath)
}

// DisableAutostart removes the daemon from OS autostart.
func (a *App) DisableAutostart() error {
	return autostart.New().Disable()
}

// IsAutostartEnabled returns whether the daemon is set to auto-start.
func (a *App) IsAutostartEnabled() (bool, error) {
	return autostart.New().IsEnabled()
}
