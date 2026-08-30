package ipc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mahmud-r-farhan/chronotrace/pkg/storage"
)

// Client is a lightweight HTTP client for the daemon IPC API.
// Used by the GUI to communicate with the running daemon.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a client connected to the given address.
func NewClient(addr string) *Client {
	if addr == "" {
		addr = DefaultAddr
	}
	return &Client{
		baseURL: "http://" + addr,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Ping checks if the daemon is running.
func (c *Client) Ping() error {
	resp, err := c.httpClient.Get(c.baseURL + apiPrefix + "/status")
	if err != nil {
		return fmt.Errorf("daemon not reachable: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("daemon returned status %d", resp.StatusCode)
	}
	return nil
}

// GetStatus returns the daemon status.
func (c *Client) GetStatus() (*StatusResponse, error) {
	var s StatusResponse
	return &s, c.get(apiPrefix+"/status", &s)
}

// GetUsageToday returns today's per-app usage.
func (c *Client) GetUsageToday() ([]storage.AppUsage, error) {
	var apps []storage.AppUsage
	return apps, c.get(apiPrefix+"/usage/today", &apps)
}

// GetUsageWeek returns the past 7 days of per-app usage.
func (c *Client) GetUsageWeek() ([]storage.AppUsage, error) {
	var apps []storage.AppUsage
	return apps, c.get(apiPrefix+"/usage/week", &apps)
}

// GetUsageMonth returns the past 30 days of per-app usage.
func (c *Client) GetUsageMonth() ([]storage.AppUsage, error) {
	var apps []storage.AppUsage
	return apps, c.get(apiPrefix+"/usage/month", &apps)
}

// GetTimeline returns the hourly timeline for a date (YYYY-MM-DD), or today if empty.
func (c *Client) GetTimeline(date string) ([]storage.TimelineSlot, error) {
	var slots []storage.TimelineSlot
	path := apiPrefix + "/usage/timeline"
	if date != "" {
		path += "?date=" + date
	}
	return slots, c.get(path, &slots)
}

// GetSummary returns a day summary for the given date.
func (c *Client) GetSummary(date string) (*storage.DaySummary, error) {
	var s storage.DaySummary
	path := apiPrefix + "/usage/summary"
	if date != "" {
		path += "?date=" + date
	}
	return &s, c.get(path, &s)
}

func (c *Client) get(path string, out interface{}) error {
	resp, err := c.httpClient.Get(c.baseURL + path)
	if err != nil {
		return fmt.Errorf("GET %s: %w", path, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s: status %d", path, resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
