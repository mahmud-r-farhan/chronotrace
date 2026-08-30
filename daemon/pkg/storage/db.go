// Package storage manages the SQLite database for ChronoTrace.
// It uses github.com/glebarez/go-sqlite (pure Go, no CGO) for
// cross-platform compatibility.
package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

const (
	driverName    = "sqlite"
	batchInterval = 45 * time.Second // flush in-memory buffer every 45 s
)

// UsageRecord is a single in-memory usage event before it is flushed to disk.
type UsageRecord struct {
	AppName     string
	WindowTitle string
	Duration    int64 // seconds
	RecordedAt  time.Time
}

// DB wraps an *sql.DB with an in-memory batch buffer for low I/O writes.
type DB struct {
	db     *sql.DB
	mu     sync.Mutex
	buffer []UsageRecord
	done   chan struct{}
	wg     sync.WaitGroup
}

// Open returns an initialised DB, creating the database file and schema if needed.
func Open() (*DB, error) {
	dir, err := dataDir()
	if err != nil {
		return nil, fmt.Errorf("storage: resolve data dir: %w", err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("storage: create data dir: %w", err)
	}

	dbPath := filepath.Join(dir, "data.db")
	sqlDB, err := sql.Open(driverName, dbPath)
	if err != nil {
		return nil, fmt.Errorf("storage: open db: %w", err)
	}

	// Single writer connection is optimal for WAL mode + low contention.
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(0)

	store := &DB{
		db:   sqlDB,
		done: make(chan struct{}),
	}

	if err := store.migrate(); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("storage: migration: %w", err)
	}

	store.wg.Add(1)
	go store.flushLoop()

	log.Printf("[storage] opened database at %s", dbPath)
	return store, nil
}

// migrate runs idempotent schema creation and enables WAL mode for performance.
func (d *DB) migrate() error {
	pragmas := []string{
		`PRAGMA journal_mode=WAL`,
		`PRAGMA synchronous=NORMAL`,
		`PRAGMA temp_store=MEMORY`,
		`PRAGMA cache_size=-8000`, // 8 MB page cache
	}
	for _, p := range pragmas {
		if _, err := d.db.Exec(p); err != nil {
			return fmt.Errorf("pragma %q: %w", p, err)
		}
	}

	schema := `
CREATE TABLE IF NOT EXISTS schema_version (
    version INTEGER PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS usage_logs (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    app_name         TEXT    NOT NULL,
    window_title     TEXT,
    duration_seconds INTEGER NOT NULL,
    recorded_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_recorded_at ON usage_logs(recorded_at);
CREATE INDEX IF NOT EXISTS idx_app_name    ON usage_logs(app_name);

INSERT OR IGNORE INTO schema_version (version) VALUES (1);
`
	if _, err := d.db.Exec(schema); err != nil {
		return fmt.Errorf("create schema: %w", err)
	}
	return nil
}

// Add enqueues a usage record into the in-memory buffer (thread-safe).
func (d *DB) Add(rec UsageRecord) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.buffer = append(d.buffer, rec)
}

// flushLoop periodically flushes the buffer to SQLite.
func (d *DB) flushLoop() {
	defer d.wg.Done()
	ticker := time.NewTicker(batchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := d.Flush(); err != nil {
				log.Printf("[storage] flush error: %v", err)
			}
		case <-d.done:
			// Final flush before shutdown.
			if err := d.Flush(); err != nil {
				log.Printf("[storage] final flush error: %v", err)
			}
			return
		}
	}
}

// Flush writes all buffered records to SQLite in a single transaction.
func (d *DB) Flush() error {
	d.mu.Lock()
	if len(d.buffer) == 0 {
		d.mu.Unlock()
		return nil
	}
	batch := d.buffer
	d.buffer = nil
	d.mu.Unlock()

	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	stmt, err := tx.Prepare(
		`INSERT INTO usage_logs (app_name, window_title, duration_seconds, recorded_at)
		 VALUES (?, ?, ?, ?)`,
	)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	for _, r := range batch {
		if _, err := stmt.Exec(r.AppName, r.WindowTitle, r.Duration, r.RecordedAt); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("insert: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	log.Printf("[storage] flushed %d records", len(batch))
	return nil
}

// Close flushes remaining data and closes the database cleanly.
func (d *DB) Close() error {
	close(d.done)
	d.wg.Wait()
	return d.db.Close()
}

// dataDir returns the OS-appropriate application data directory for ChronoTrace.
func dataDir() (string, error) {
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("%%APPDATA%% not set")
		}
		return filepath.Join(appData, "ChronoTrace"), nil

	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, "Library", "Application Support", "ChronoTrace"), nil

	default: // Linux and others
		if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
			return filepath.Join(xdg, "chronotrace"), nil
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".local", "share", "chronotrace"), nil
	}
}
