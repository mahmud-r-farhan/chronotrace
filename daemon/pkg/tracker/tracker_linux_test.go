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

func TestParseGDBusEval(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		wantOK   bool
		wantBody string
	}{
		{"simple tuple", "(true, 'Firefox|firefox')\n", true, "Firefox|firefox"},
		// The old strings.Trim cutset-based parsing corrupted payloads whose
		// first/last characters appeared in "(true, ')\n" — e.g. a class name
		// ending in "e" lost that final character.
		{"class ending in cutset chars", "(true, 'Terminal|google-chrome')", true, "Terminal|google-chrome"},
		{"title starting with cutset chars", "(true, 'true crime|vlc')", true, "true crime|vlc"},
		{"double-quoted payload", `(true, "it's here|xterm")`, true, "it's here|xterm"},
		{"failed eval", "(false, '')", false, ""},
		{"empty success payload", "(true, '')", false, ""},
		{"dbus error", "Error: GDBus.Error:org.freedesktop.DBus.Error.NotSupported: Method Eval not available", false, ""},
		{"empty output", "", false, ""},
	}

	for _, tt := range tests {
		body, ok := parseGDBusEval(tt.raw)
		if ok != tt.wantOK || body != tt.wantBody {
			t.Errorf("%s: parseGDBusEval(%q) = (%q, %v), want (%q, %v)",
				tt.name, tt.raw, body, ok, tt.wantBody, tt.wantOK)
		}
	}
}
