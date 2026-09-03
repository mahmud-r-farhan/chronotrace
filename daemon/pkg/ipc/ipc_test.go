package ipc

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mahmud-r-farhan/chronotrace/pkg/storage"
)

func TestIPCServerAndClient(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmpDir)

	db, err := storage.Open()
	if err != nil {
		t.Fatalf("storage.Open failed: %v", err)
	}
	defer db.Close()

	now := time.Now()
	db.Add(storage.UsageRecord{
		AppName:     "TestApp",
		WindowTitle: "Test Window",
		Duration:    100,
		RecordedAt:  now,
	})
	if err := db.Flush(); err != nil {
		t.Fatalf("db.Flush failed: %v", err)
	}

	server := New("127.0.0.1:0", db, "1.0.0-test")

	// Use httptest.Server wrapping the server handler
	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case apiPrefix + "/status":
			server.handleStatus(w, r)
		case apiPrefix + "/usage/today":
			server.handleToday(w, r)
		case apiPrefix + "/usage/week":
			server.handleWeek(w, r)
		case apiPrefix + "/usage/month":
			server.handleMonth(w, r)
		case apiPrefix + "/usage/timeline":
			server.handleTimeline(w, r)
		case apiPrefix + "/usage/summary":
			server.handleSummary(w, r)
		default:
			http.NotFound(w, r)
		}
	}))

	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Connect Client to test server
	// ts.URL looks like "http://127.0.0.1:PORT"
	client := &Client{
		baseURL:    ts.URL,
		httpClient: ts.Client(),
	}

	// 1. Ping
	if err := client.Ping(); err != nil {
		t.Errorf("client.Ping failed: %v", err)
	}

	// 2. GetStatus
	status, err := client.GetStatus()
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Status != "ok" || status.Version != "1.0.0-test" {
		t.Errorf("unexpected status: %+v", status)
	}

	// 3. GetUsageToday
	todayApps, err := client.GetUsageToday()
	if err != nil {
		t.Fatalf("GetUsageToday failed: %v", err)
	}
	if len(todayApps) != 1 || todayApps[0].AppName != "TestApp" {
		t.Errorf("unexpected todayApps: %+v", todayApps)
	}

	// 4. GetUsageWeek
	weekApps, err := client.GetUsageWeek()
	if err != nil {
		t.Fatalf("GetUsageWeek failed: %v", err)
	}
	if len(weekApps) != 1 {
		t.Errorf("unexpected weekApps count: %d", len(weekApps))
	}

	// 5. GetUsageMonth
	monthApps, err := client.GetUsageMonth()
	if err != nil {
		t.Fatalf("GetUsageMonth failed: %v", err)
	}
	if len(monthApps) != 1 {
		t.Errorf("unexpected monthApps count: %d", len(monthApps))
	}

	// 6. GetTimeline
	timeline, err := client.GetTimeline(now.Format("2006-01-02"))
	if err != nil {
		t.Fatalf("GetTimeline failed: %v", err)
	}
	if len(timeline) != 24 {
		t.Errorf("expected 24 slots, got %d", len(timeline))
	}

	// 7. GetSummary
	summary, err := client.GetSummary(now.Format("2006-01-02"))
	if err != nil {
		t.Fatalf("GetSummary failed: %v", err)
	}
	if summary.TotalSeconds != 100 || summary.AppCount != 1 {
		t.Errorf("unexpected summary: %+v", summary)
	}

	// 8. CORS OPTIONS check
	req, _ := http.NewRequest(http.MethodOptions, ts.URL+apiPrefix+"/status", nil)
	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("OPTIONS request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204 No Content for OPTIONS, got %d", resp.StatusCode)
	}
	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("missing CORS header")
	}
}
