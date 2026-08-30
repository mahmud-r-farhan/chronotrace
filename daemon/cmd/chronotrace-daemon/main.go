// chronotrace-daemon is the headless background service for ChronoTrace.
// It polls the active window every 2-3 seconds, buffers usage records in memory,
// batch-writes them to SQLite every 45 seconds, and exposes a local REST API
// on 127.0.0.1:42069 for the GUI to consume.
//
// Design goals:
//   - RAM usage < 15 MB
//   - CPU usage ~0% (polling at 2-3 s intervals)
//   - Never panic; always log and recover
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/mahmud-r-farhan/chronotrace/pkg/autostart"
	"github.com/mahmud-r-farhan/chronotrace/pkg/ipc"
	"github.com/mahmud-r-farhan/chronotrace/pkg/storage"
	"github.com/mahmud-r-farhan/chronotrace/pkg/tracker"
)

const (
	version    = "0.1.0"
	appName    = "ChronoTrace Daemon"
	minPollMs  = 2000
	maxPollMs  = 3000
)

func main() {
	// --- CLI flags ---
	addr := flag.String("addr", ipc.DefaultAddr, "IPC server listen address")
	installAutostart := flag.Bool("autostart-install", false, "Install daemon in OS autostart and exit")
	removeAutostart := flag.Bool("autostart-remove", false, "Remove daemon from OS autostart and exit")
	showVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("%s v%s (%s/%s)\n", appName, version, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	// --- Autostart management commands ---
	as := autostart.New()
	if *installAutostart {
		exe, err := os.Executable()
		if err != nil {
			log.Fatalf("[daemon] cannot determine executable path: %v", err)
		}
		exe, _ = filepath.Abs(exe)
		if err := as.Enable(exe); err != nil {
			log.Fatalf("[daemon] autostart install failed: %v", err)
		}
		log.Printf("[daemon] autostart installed for %s", exe)
		os.Exit(0)
	}
	if *removeAutostart {
		if err := as.Disable(); err != nil {
			log.Fatalf("[daemon] autostart remove failed: %v", err)
		}
		log.Println("[daemon] autostart removed")
		os.Exit(0)
	}

	// --- Single-instance lock ---
	lockPath := lockFilePath()
	if err := acquireLock(lockPath); err != nil {
		log.Fatalf("[daemon] another instance is already running: %v", err)
	}
	defer releaseLock(lockPath)

	log.Printf("[daemon] %s v%s starting on %s/%s", appName, version, runtime.GOOS, runtime.GOARCH)

	// --- Storage ---
	db, err := storage.Open()
	if err != nil {
		log.Fatalf("[daemon] storage init failed: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("[daemon] db close error: %v", err)
		}
	}()

	// --- IPC server ---
	ipcServer := ipc.New(*addr, db, version)
	go func() {
		if err := ipcServer.Start(); err != nil {
			log.Printf("[daemon] IPC server stopped: %v", err)
		}
	}()

	// --- Tracker ---
	t := tracker.New()

	// --- Graceful shutdown context ---
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Printf("[daemon] received signal %s, shutting down...", sig)
		ipcServer.Stop()
		cancel()
	}()

	// --- Main polling loop ---
	log.Println("[daemon] polling loop started")
	runPollingLoop(ctx, t, db)
	log.Println("[daemon] exited cleanly")
}

// runPollingLoop polls the active window at 2-3 second jittered intervals.
// It aggregates consecutive time spent in the same app before buffering.
func runPollingLoop(ctx context.Context, t tracker.Tracker, db *storage.DB) {
	var (
		lastApp   string
		lastTitle string
		lastSeen  time.Time
	)

	for {
		// Jittered sleep: 2000-3000 ms to avoid wall-clock alignment.
		jitter := time.Duration(minPollMs+rand.Intn(maxPollMs-minPollMs)) * time.Millisecond
		select {
		case <-ctx.Done():
			// Flush any last pending duration on shutdown.
			if lastApp != "" && !lastSeen.IsZero() {
				elapsed := int64(time.Since(lastSeen).Seconds())
				if elapsed > 0 {
					db.Add(storage.UsageRecord{
						AppName:     lastApp,
						WindowTitle: lastTitle,
						Duration:    elapsed,
						RecordedAt:  time.Now(),
					})
				}
			}
			return
		case <-time.After(jitter):
		}

		info, err := t.GetActiveWindow()
		if err != nil || info == nil {
			continue
		}
		if info.AppName == "" {
			info.AppName = "Unknown"
		}

		now := time.Now()

		if lastApp == "" {
			// First poll.
			lastApp = info.AppName
			lastTitle = info.WindowTitle
			lastSeen = now
			continue
		}

		elapsed := int64(now.Sub(lastSeen).Seconds())

		if info.AppName != lastApp {
			// App changed — record the time spent on the previous app.
			if elapsed > 0 && lastApp != "" {
				db.Add(storage.UsageRecord{
					AppName:     lastApp,
					WindowTitle: lastTitle,
					Duration:    elapsed,
					RecordedAt:  now,
				})
			}
			lastApp = info.AppName
			lastTitle = info.WindowTitle
			lastSeen = now
		}
		// Same app — update title in case it changed (e.g. browser tab).
		lastTitle = info.WindowTitle
	}
}

// lockFilePath returns the OS-appropriate path for the daemon lock file.
func lockFilePath() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("TEMP"), "chronotrace-daemon.lock")
	case "darwin":
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "Library", "Application Support", "ChronoTrace", "daemon.lock")
	default:
		if xdg := os.Getenv("XDG_RUNTIME_DIR"); xdg != "" {
			return filepath.Join(xdg, "chronotrace-daemon.lock")
		}
		return filepath.Join(os.TempDir(), "chronotrace-daemon.lock")
	}
}

// acquireLock creates a lock file containing the current PID.
// Returns an error if the lock file already exists and the process is still alive.
func acquireLock(path string) error {
	_ = os.MkdirAll(filepath.Dir(path), 0o755)

	// Check if existing lock is stale.
	if data, err := os.ReadFile(path); err == nil {
		var oldPID int
		if _, err := fmt.Sscan(string(data), &oldPID); err == nil {
			if isProcessAlive(oldPID) {
				return fmt.Errorf("process %d is already running", oldPID)
			}
			// Stale lock — remove it.
			_ = os.Remove(path)
		}
	}

	return os.WriteFile(path, []byte(fmt.Sprintf("%d", os.Getpid())), 0o644)
}

// releaseLock removes the lock file.
func releaseLock(path string) {
	_ = os.Remove(path)
}
