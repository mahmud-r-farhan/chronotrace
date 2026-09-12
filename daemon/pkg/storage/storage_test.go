package storage

import (
	"path/filepath"
	"testing"
	"time"
)

func TestStorage(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmpDir)

	db, err := Open()
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer db.Close()

	now := time.Now()
	todayStr := now.Format("2006-01-02")

	// Add test records
	db.Add(UsageRecord{
		AppName:     "VSCode",
		WindowTitle: "main.go - chronotrace",
		Duration:    120,
		RecordedAt:  now,
	})
	db.Add(UsageRecord{
		AppName:     "Chrome",
		WindowTitle: "GitHub",
		Duration:    60,
		RecordedAt:  now,
	})
	db.Add(UsageRecord{
		AppName:     "VSCode",
		WindowTitle: "db.go - chronotrace",
		Duration:    180,
		RecordedAt:  now,
	})

	if err := db.Flush(); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	// Test GetUsageByDay
	apps, err := db.GetUsageByDay(todayStr)
	if err != nil {
		t.Fatalf("GetUsageByDay failed: %v", err)
	}
	if len(apps) != 2 {
		t.Fatalf("expected 2 distinct apps, got %d", len(apps))
	}
	if apps[0].AppName != "VSCode" || apps[0].TotalSeconds != 300 {
		t.Errorf("unexpected top app: %+v", apps[0])
	}
	if apps[1].AppName != "Chrome" || apps[1].TotalSeconds != 60 {
		t.Errorf("unexpected second app: %+v", apps[1])
	}

	// Test GetUsageByWeek
	weekApps, err := db.GetUsageByWeek()
	if err != nil {
		t.Fatalf("GetUsageByWeek failed: %v", err)
	}
	if len(weekApps) < 2 {
		t.Errorf("expected at least 2 apps in week usage, got %d", len(weekApps))
	}

	// Test GetUsageByMonth
	monthApps, err := db.GetUsageByMonth()
	if err != nil {
		t.Fatalf("GetUsageByMonth failed: %v", err)
	}
	if len(monthApps) < 2 {
		t.Errorf("expected at least 2 apps in month usage, got %d", len(monthApps))
	}

	// Test GetTimeline
	timeline, err := db.GetTimeline(todayStr)
	if err != nil {
		t.Fatalf("GetTimeline failed: %v", err)
	}
	if len(timeline) != 24 {
		t.Fatalf("expected 24 timeline slots, got %d", len(timeline))
	}
	hour := now.Hour()
	if timeline[hour].TotalSeconds != 360 {
		t.Errorf("expected 360s in hour %d, got %d", hour, timeline[hour].TotalSeconds)
	}

	// Test GetDaySummary
	summary, err := db.GetDaySummary(todayStr)
	if err != nil {
		t.Fatalf("GetDaySummary failed: %v", err)
	}
	if summary.TotalSeconds != 360 {
		t.Errorf("expected total 360s, got %d", summary.TotalSeconds)
	}
	if summary.AppCount != 2 {
		t.Errorf("expected app count 2, got %d", summary.AppCount)
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		secs     int64
		expected string
	}{
		{30, "30s"},
		{60, "1m"},
		{90, "1m"},
		{3600, "1h"},
		{3660, "1h 1m"},
	}

	for _, tt := range tests {
		got := formatDuration(tt.secs)
		if got != tt.expected {
			t.Errorf("formatDuration(%d) = %q, want %q", tt.secs, got, tt.expected)
		}
	}
}

// TestSessionCounting verifies that long sessions stored as consecutive chunks
// count as a single session, while chunks separated by a long idle gap count
// as separate sessions.
func TestSessionCounting(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmpDir)

	db, err := Open()
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer db.Close()

	now := time.Now()
	base := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, time.Local)
	day := base.Format("2006-01-02")

	// One long session: three consecutive ~60 s chunks (as the daemon emits).
	db.Add(UsageRecord{AppName: "VSCode", Duration: 60, RecordedAt: base.Add(60 * time.Second)})
	db.Add(UsageRecord{AppName: "VSCode", Duration: 60, RecordedAt: base.Add(121 * time.Second)})
	db.Add(UsageRecord{AppName: "VSCode", Duration: 60, RecordedAt: base.Add(182 * time.Second)})
	// Same app again after a 10-minute gap => a second session.
	db.Add(UsageRecord{AppName: "VSCode", Duration: 60, RecordedAt: base.Add(10*time.Minute + 242*time.Second)})
	// A different app with a single record.
	db.Add(UsageRecord{AppName: "Chrome", Duration: 90, RecordedAt: base.Add(300 * time.Second)})

	if err := db.Flush(); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	apps, err := db.GetUsageByDay(day)
	if err != nil {
		t.Fatalf("GetUsageByDay failed: %v", err)
	}
	if len(apps) != 2 {
		t.Fatalf("expected 2 apps, got %d: %+v", len(apps), apps)
	}

	byName := map[string]AppUsage{}
	for _, a := range apps {
		byName[a.AppName] = a
	}
	if got := byName["VSCode"]; got.TotalSeconds != 240 || got.SessionCount != 2 {
		t.Errorf("VSCode: got total=%d sessions=%d, want total=240 sessions=2", got.TotalSeconds, got.SessionCount)
	}
	if got := byName["Chrome"]; got.TotalSeconds != 90 || got.SessionCount != 1 {
		t.Errorf("Chrome: got total=%d sessions=%d, want total=90 sessions=1", got.TotalSeconds, got.SessionCount)
	}

	// LastSeen must be parsed as local time, not UTC.
	if !byName["VSCode"].LastSeen.Equal(base.Add(10*time.Minute + 242*time.Second)) {
		t.Errorf("VSCode LastSeen = %v, want %v", byName["VSCode"].LastSeen, base.Add(10*time.Minute+242*time.Second))
	}
}

// TestCloseIdempotent verifies Close can be called more than once safely.
func TestCloseIdempotent(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmpDir)

	db, err := Open()
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	db.Add(UsageRecord{AppName: "App", Duration: 5, RecordedAt: time.Now()})
	if err := db.Close(); err != nil {
		t.Fatalf("first Close failed: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("second Close should be a no-op, got: %v", err)
	}

	// The record added before Close must have been flushed to disk.
	db2, err := Open()
	if err != nil {
		t.Fatalf("reopen failed: %v", err)
	}
	defer db2.Close()

	apps, err := db2.GetUsageByDay("")
	if err != nil {
		t.Fatalf("GetUsageByDay failed: %v", err)
	}
	if len(apps) != 1 || apps[0].AppName != "App" || apps[0].TotalSeconds != 5 {
		t.Errorf("expected the buffered record to survive Close, got %+v", apps)
	}
}

func TestParseLocalTime(t *testing.T) {
	tm := parseLocalTime("2026-09-12 14:30:00")
	if tm.IsZero() {
		t.Fatal("expected a non-zero time")
	}
	if tm.Location() != time.Local {
		t.Errorf("expected local timezone, got %v", tm.Location())
	}
	if tm.Hour() != 14 || tm.Minute() != 30 {
		t.Errorf("unexpected time of day: %v", tm)
	}
	if !parseLocalTime("not a timestamp").IsZero() {
		t.Error("expected zero time for unparseable input")
	}
}

func TestDataDirOverride(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmpDir)

	dir, err := dataDir()
	if err != nil {
		t.Fatalf("dataDir error: %v", err)
	}
	if dir != filepath.Join(tmpDir, "chronotrace") {
		t.Errorf("unexpected dataDir: %s", dir)
	}
}
