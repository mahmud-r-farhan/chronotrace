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
