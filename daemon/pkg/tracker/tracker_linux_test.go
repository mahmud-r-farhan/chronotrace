package tracker

import (
	"testing"
)

func TestNewTracker(t *testing.T) {
	tr := New()
	if tr == nil {
		t.Fatal("New() returned nil Tracker")
	}

	info, err := tr.GetActiveWindow()
	if err != nil {
		t.Logf("GetActiveWindow returned error (expected in headless environment): %v", err)
	} else {
		t.Logf("GetActiveWindow returned app: %q, title: %q", info.AppName, info.WindowTitle)
	}
}

func TestParseWMClass(t *testing.T) {
	tests := []struct {
		raw      string
		expected string
	}{
		{`WM_CLASS(STRING) = "google-chrome", "Google-chrome"`, "Google-chrome"},
		{`WM_CLASS(STRING) = "code", "Code"`, "Code"},
		{`invalid format`, ""},
	}

	for _, tt := range tests {
		got := parseWMClass(tt.raw)
		if got != tt.expected {
			t.Errorf("parseWMClass(%q) = %q, want %q", tt.raw, got, tt.expected)
		}
	}
}
