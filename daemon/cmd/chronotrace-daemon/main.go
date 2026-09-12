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
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/mahmud-r-farhan/chronotrace/pkg/autostart"
	"github.com/mahmud-r-farhan/chronotrace/pkg/ipc"
	"github.com/mahmud-r-farhan/chronotrace/pkg/storage"
	"github.com/mahmud-r-farhan/chronotrace/pkg/tracker"
)

// version is stamped at build time via -ldflags "-X main.version=...".
// It must remain a var: the linker cannot override constants.
var version = "0.1.0"

const (
	appName   = "ChronoTrace Daemon"
	minPollMs = 2000
	maxPollMs = 3000

	// maxChunkSeconds limits how much time a single buffered usage record may
	// cover. Long sessions in the same app are split into chunks so the hourly
	// timeline attributes time to the correct hour instead of piling a whole
	// session into the hour in which it ended.
	maxChunkSeconds = 60

	// maxStaleSeconds is the largest gap between two successful polls that is
	// still trusted as continuous usage. Bigger gaps mean the loop was starved
	// (system suspend, locked screen, tracker errors) — that time cannot be
	// attributed to any app and is dropped instead of inflating the last one.
	maxStaleSeconds = 2 * maxChunkSeconds
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
		exe, err = filepath.Abs(exe)
		if err != nil {
			log.Fatalf("[daemon] cannot resolve absolute executable path: %v", err)
		}
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
// Time spent in the same app is aggregated and buffered in chunks of at most
// maxChunkSeconds, so hourly statistics stay accurate even for long sessions.
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
			// Flush any last pending duration on shutdown (if it is recent
			// enough to be trustworthy).
			if lastApp != "" && !lastSeen.IsZero() {
				elapsed := int64(time.Since(lastSeen).Seconds())
				if elapsed > 0 && elapsed <= maxStaleSeconds {
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

		if elapsed > maxStaleSeconds {
			// The loop was starved since the last successful poll (suspend,
			// locked screen, tracker outage). Drop the unverifiable gap and
			// restart tracking from the current window.
			lastApp = info.AppName
			lastTitle = info.WindowTitle
			lastSeen = now
			continue
		}

		if info.AppName != lastApp {
			// App changed — record the time spent on the previous app.
			if elapsed > 0 {
				db.Add(storage.UsageRecord{
					AppName:     lastApp,
					WindowTitle: lastTitle,
					Duration:    elapsed,
					RecordedAt:  now,
				})
			}
			lastApp = info.AppName
			lastSeen = now
		} else if elapsed >= maxChunkSeconds {
			// Same app for a while — emit a chunk so long sessions spread
			// across the hours they actually cover.
			db.Add(storage.UsageRecord{
				AppName:     lastApp,
				WindowTitle: info.WindowTitle,
				Duration:    elapsed,
				RecordedAt:  now,
			})
			lastSeen = now
		}
		// Track the latest title in case it changed (e.g. browser tab).
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

// acquireLock atomically creates a lock file containing the current PID.
// It returns an error if a live process already holds the lock; stale locks
// left behind by crashed processes are detected and reclaimed.
func acquireLock(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create lock dir: %w", err)
	}

	for attempt := 0; attempt < 2; attempt++ {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err == nil {
			_, _ = fmt.Fprintf(f, "%d\n", os.Getpid())
			return f.Close()
		}
		if !errors.Is(err, os.ErrExist) {
			return err
		}

		// Lock file exists — check whether the owning process is still alive.
		if data, rerr := os.ReadFile(path); rerr == nil {
			var oldPID int
			if _, serr := fmt.Sscan(strings.TrimSpace(string(data)), &oldPID); serr == nil &&
				oldPID > 0 && oldPID != os.Getpid() && isProcessAlive(oldPID) {
				return fmt.Errorf("process %d is already running", oldPID)
			}
		}

		// Stale or unreadable lock — reclaim it and retry once.
		if rerr := os.Remove(path); rerr != nil && !errors.Is(rerr, os.ErrNotExist) {
			return rerr
		}
	}
	return fmt.Errorf("could not acquire lock at %s", path)
}

// releaseLock removes the lock file.
func releaseLock(path string) {
	_ = os.Remove(path)
}
