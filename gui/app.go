package main

import (
	"context"
	"log"

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
}

func (a *App) shutdown(ctx context.Context) {
	log.Println("[gui] shutdown")
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
		"connected":       true,
		"version":         s.Version,
		"uptime_seconds":  s.UptimeSec,
		"addr":            s.Addr,
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
// Pass empty string for today.
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
	// The daemon executable is expected to be next to the GUI or on PATH.
	// The user should run the daemon with --autostart-install for first-time setup.
	return mgr.Enable("chronotrace-daemon")
}

// DisableAutostart removes the daemon from OS autostart.
func (a *App) DisableAutostart() error {
	return autostart.New().Disable()
}

// IsAutostartEnabled returns whether the daemon is set to auto-start.
func (a *App) IsAutostartEnabled() (bool, error) {
	return autostart.New().IsEnabled()
}
