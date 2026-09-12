package ipc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	if err := c.get(apiPrefix+"/status", &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// GetUsageToday returns today's per-app usage.
func (c *Client) GetUsageToday() ([]storage.AppUsage, error) {
	var apps []storage.AppUsage
	err := c.get(apiPrefix+"/usage/today", &apps)
	return apps, err
}

// GetUsageWeek returns the past 7 days of per-app usage.
func (c *Client) GetUsageWeek() ([]storage.AppUsage, error) {
	var apps []storage.AppUsage
	err := c.get(apiPrefix+"/usage/week", &apps)
	return apps, err
}

// GetUsageMonth returns the past 30 days of per-app usage.
func (c *Client) GetUsageMonth() ([]storage.AppUsage, error) {
	var apps []storage.AppUsage
	err := c.get(apiPrefix+"/usage/month", &apps)
	return apps, err
}

// GetTimeline returns the hourly timeline for a date (YYYY-MM-DD), or today if empty.
func (c *Client) GetTimeline(date string) ([]storage.TimelineSlot, error) {
	var slots []storage.TimelineSlot
	err := c.get(apiPrefix+"/usage/timeline"+dateQuery(date), &slots)
	return slots, err
}

// GetSummary returns a day summary for the given date.
func (c *Client) GetSummary(date string) (*storage.DaySummary, error) {
	var s storage.DaySummary
	err := c.get(apiPrefix+"/usage/summary"+dateQuery(date), &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// dateQuery builds an escaped "?date=..." query string, or "" when date is empty.
func dateQuery(date string) string {
	if date == "" {
		return ""
	}
	return "?date=" + url.QueryEscape(date)
}

func (c *Client) get(path string, out interface{}) error {
	resp, err := c.httpClient.Get(c.baseURL + path)
	if err != nil {
		return fmt.Errorf("GET %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Surface the daemon's JSON error payload when there is one.
		var body struct {
			Error string `json:"error"`
		}
		if json.NewDecoder(resp.Body).Decode(&body) == nil && body.Error != "" {
			return fmt.Errorf("GET %s: status %d: %s", path, resp.StatusCode, body.Error)
		}
		return fmt.Errorf("GET %s: unexpected status %d", path, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("GET %s: decode response: %w", path, err)
	}
	return nil
}
