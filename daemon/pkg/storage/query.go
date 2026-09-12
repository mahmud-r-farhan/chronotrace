package storage

import (
	"database/sql"
	"fmt"
	"time"
)

// sessionGapSeconds is the idle gap after which a row for the same app counts
// as a new session instead of a continuation chunk of a long one. The daemon
// splits long sessions into ~60 s chunks (see the polling loop), so genuine
// chunks are seconds apart while real session breaks are minutes apart.
const sessionGapSeconds = 90

// AppUsage represents aggregated usage for a single application.
type AppUsage struct {
	AppName       string    `json:"app_name"`
	TotalSeconds  int64     `json:"total_seconds"`
	SessionCount  int       `json:"session_count"`
	LastSeen      time.Time `json:"last_seen"`
	FormattedTime string    `json:"formatted_time"`
}

// TimelineSlot represents usage for one hour of the day.
type TimelineSlot struct {
	Hour         int   `json:"hour"`
	TotalSeconds int64 `json:"total_seconds"`
}

// DaySummary is a high-level overview for a single day.
type DaySummary struct {
	Date         string     `json:"date"`
	TotalSeconds int64      `json:"total_seconds"`
	AppCount     int        `json:"app_count"`
	TopApps      []AppUsage `json:"top_apps"`
}

// GetUsageByDay returns aggregated per-app usage for the given date (YYYY-MM-DD).
// Passing an empty string defaults to today.
func (d *DB) GetUsageByDay(date string) ([]AppUsage, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	query := aggregationQuery(`date(recorded_at) = ?`)
	return d.queryAppUsage(query, sessionGapSeconds, date)
}

// GetUsageByWeek returns aggregated per-app usage for the last 7 rolling days.
// Timestamps are stored as local wall-clock time, so the window boundary must
// be computed in localtime rather than SQLite's default UTC.
func (d *DB) GetUsageByWeek() ([]AppUsage, error) {
	query := aggregationQuery(`recorded_at >= datetime('now', 'localtime', '-7 days')`)
	return d.queryAppUsage(query, sessionGapSeconds)
}

// GetUsageByMonth returns aggregated per-app usage for the last 30 rolling days.
func (d *DB) GetUsageByMonth() ([]AppUsage, error) {
	query := aggregationQuery(`recorded_at >= datetime('now', 'localtime', '-30 days')`)
	return d.queryAppUsage(query, sessionGapSeconds)
}

// aggregationQuery builds the per-app aggregation SELECT for the given WHERE
// clause. The inner query flags rows that begin a new session using a window
// function: the first row of an app is always a new session, and later rows
// only count when the idle gap between the end of the previous record and the
// start of this record (elapsed time between them minus this row's duration)
// exceeds the sessionGapSeconds bind parameter.
func aggregationQuery(where string) string {
	return `
		SELECT   app_name,
		         SUM(duration_seconds) AS total_seconds,
		         SUM(new_session)      AS session_count,
		         MAX(recorded_at)      AS last_seen
		FROM (
		    SELECT app_name,
		           duration_seconds,
		           recorded_at,
		           CASE
		               WHEN LAG(recorded_at) OVER w IS NULL THEN 1
		               WHEN (julianday(recorded_at) - julianday(LAG(recorded_at) OVER w)) * 86400.0
		                    - duration_seconds > ? THEN 1
		               ELSE 0
		           END AS new_session
		    FROM usage_logs
		    WHERE ` + where + `
		    WINDOW w AS (PARTITION BY app_name ORDER BY recorded_at, id)
		)
		GROUP BY app_name
		ORDER BY total_seconds DESC
	`
}

// GetTimeline returns an hourly breakdown for the given date.
func (d *DB) GetTimeline(date string) ([]TimelineSlot, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	query := `
		SELECT   CAST(strftime('%H', recorded_at) AS INTEGER) AS hour,
		         SUM(duration_seconds) AS total_seconds
		FROM     usage_logs
		WHERE    date(recorded_at) = ?
		GROUP BY hour
		ORDER BY hour
	`
	rows, err := d.db.Query(query, date)
	if err != nil {
		return nil, fmt.Errorf("GetTimeline query: %w", err)
	}
	defer rows.Close()

	// Pre-fill all 24 hours with 0.
	slots := make([]TimelineSlot, 24)
	for i := range slots {
		slots[i].Hour = i
	}

	for rows.Next() {
		var hour int
		var secs int64
		if err := rows.Scan(&hour, &secs); err != nil {
			return nil, fmt.Errorf("GetTimeline scan: %w", err)
		}
		if hour >= 0 && hour < 24 {
			slots[hour].TotalSeconds = secs
		}
	}
	return slots, rows.Err()
}

// GetDaySummary returns a DaySummary for the given date.
func (d *DB) GetDaySummary(date string) (*DaySummary, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	apps, err := d.GetUsageByDay(date)
	if err != nil {
		return nil, err
	}

	var total int64
	for _, a := range apps {
		total += a.TotalSeconds
	}

	topN := apps
	if len(topN) > 5 {
		topN = topN[:5]
	}

	return &DaySummary{
		Date:         date,
		TotalSeconds: total,
		AppCount:     len(apps),
		TopApps:      topN,
	}, nil
}

// queryAppUsage is a helper that runs an AppUsage SELECT and scans the results.
func (d *DB) queryAppUsage(query string, args ...interface{}) ([]AppUsage, error) {
	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("queryAppUsage: %w", err)
	}
	defer rows.Close()

	var results []AppUsage
	for rows.Next() {
		var a AppUsage
		var lastSeenStr sql.NullString
		if err := rows.Scan(&a.AppName, &a.TotalSeconds, &a.SessionCount, &lastSeenStr); err != nil {
			return nil, fmt.Errorf("queryAppUsage scan: %w", err)
		}
		if lastSeenStr.Valid {
			a.LastSeen = parseLocalTime(lastSeenStr.String)
		}
		a.FormattedTime = formatDuration(a.TotalSeconds)
		results = append(results, a)
	}
	return results, rows.Err()
}

// parseLocalTime parses a timestamp string stored as local wall-clock time.
// Plain time.Parse would mislabel it as UTC, shifting the value when the JSON
// response is rendered in the GUI.
func parseLocalTime(s string) time.Time {
	layouts := []string{"2006-01-02 15:04:05", time.RFC3339}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return t
		}
	}
	return time.Time{}
}

// formatDuration converts seconds into a human-readable string like "2h 30m".
func formatDuration(seconds int64) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	m := seconds / 60
	if m < 60 {
		return fmt.Sprintf("%dm", m)
	}
	h := m / 60
	rem := m % 60
	if rem == 0 {
		return fmt.Sprintf("%dh", h)
	}
	return fmt.Sprintf("%dh %dm", h, rem)
}
