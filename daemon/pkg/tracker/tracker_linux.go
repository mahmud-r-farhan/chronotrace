//go:build linux

package tracker

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// linuxTracker retrieves the active window on Linux via xdotool (X11) or
// DBus/gdbus fallback for GNOME on Wayland.
type linuxTracker struct {
	method string // "xdotool", "gdbus", or "unknown"
}

func newPlatformTracker() Tracker {
	t := &linuxTracker{}
	t.detectMethod()
	return t
}

// detectMethod probes available tools and chooses the best strategy.
func (t *linuxTracker) detectMethod() {
	if _, err := exec.LookPath("xdotool"); err == nil {
		t.method = "xdotool"
		return
	}
	if _, err := exec.LookPath("gdbus"); err == nil {
		t.method = "gdbus"
		return
	}
	t.method = "unknown"
}

// GetActiveWindow retrieves the active window info on Linux.
func (t *linuxTracker) GetActiveWindow() (*ActiveWindowInfo, error) {
	switch t.method {
	case "xdotool":
		return t.getViaXdotool()
	case "gdbus":
		return t.getViaGDBus()
	default:
		return &ActiveWindowInfo{AppName: "Unknown (no xdotool/gdbus)"}, nil
	}
}

func (t *linuxTracker) getViaXdotool() (*ActiveWindowInfo, error) {
	// Get active window ID.
	widOut, err := exec.Command("xdotool", "getactivewindow").Output()
	if err != nil {
		return &ActiveWindowInfo{AppName: "Unknown"}, nil
	}
	winID := strings.TrimSpace(string(widOut))
	if winID == "" || winID == "0" {
		return &ActiveWindowInfo{AppName: "Unknown"}, nil
	}

	// Get window name.
	nameOut, _ := exec.Command("xdotool", "getwindowname", winID).Output()
	title := strings.TrimSpace(string(nameOut))

	// Identify the app via WM_CLASS: prefer xdotool's own lookup, fall back
	// to xprop when xdotool cannot resolve the class.
	appName := ""
	if classOut, cerr := exec.Command("xdotool", "getwindowclassname", winID).Output(); cerr == nil {
		appName = strings.TrimSpace(string(classOut))
	} else if xpropOut, perr := exec.Command("xprop", "-id", winID, "WM_CLASS").Output(); perr == nil {
		appName = parseWMClass(string(xpropOut))
	}
	if appName == "" {
		appName = title
	}

	info := &ActiveWindowInfo{AppName: appName, WindowTitle: title}

	// PID and executable path (best-effort: not every window sets _NET_WM_PID).
	if pidOut, perr := exec.Command("xdotool", "getwindowpid", winID).Output(); perr == nil {
		if pid, aerr := strconv.Atoi(strings.TrimSpace(string(pidOut))); aerr == nil && pid > 0 {
			info.PID = uint32(pid)
			if exe, lerr := os.Readlink(fmt.Sprintf("/proc/%d/exe", pid)); lerr == nil {
				info.ExePath = exe
			}
		}
	}

	return info, nil
}

func (t *linuxTracker) getViaGDBus() (*ActiveWindowInfo, error) {
	// GNOME Shell via DBus — works on Wayland
	out, err := exec.Command("gdbus", "call",
		"--session",
		"--dest", "org.gnome.Shell",
		"--object-path", "/org/gnome/Shell",
		"--method", "org.gnome.Shell.Eval",
		`global.display.focus_window ? global.display.focus_window.get_title() + "|" + global.display.focus_window.get_wm_class() : "unknown|unknown"`,
	).Output()
	if err != nil {
		return &ActiveWindowInfo{AppName: "Unknown"}, nil
	}

	// Parse the response — format: (true, '<title>|<class>')
	payload, ok := parseGDBusEval(string(out))
	if !ok {
		return &ActiveWindowInfo{AppName: "Unknown"}, nil
	}

	title, class := payload, ""
	if idx := strings.Index(payload, "|"); idx >= 0 {
		title, class = payload[:idx], payload[idx+1:]
	}

	info := &ActiveWindowInfo{WindowTitle: strings.TrimSpace(title)}
	info.AppName = info.WindowTitle
	if class = strings.TrimSpace(class); class != "" && class != "unknown" {
		info.AppName = class
	}
	return info, nil
}

// parseGDBusEval extracts the string result from a gdbus Eval response such as
// "(true, 'title|class')" and reports whether the evaluation succeeded.
// gdbus prints GVariant strings in single quotes; payloads that themselves
// contain an apostrophe are printed double-quoted instead.
func parseGDBusEval(raw string) (string, bool) {
	raw = strings.TrimSpace(raw)
	if !strings.HasPrefix(raw, "(true") {
		return "", false
	}
	for _, quote := range []byte{'\'', '"'} {
		start := strings.IndexByte(raw, quote)
		end := strings.LastIndexByte(raw, quote)
		// end > start+1 requires a non-empty payload between the quotes;
		// an empty result ('' or "") carries no window information.
		if start >= 0 && end > start+1 {
			return raw[start+1 : end], true
		}
	}
	return "", false
}

// parseWMClass extracts the application name from xprop WM_CLASS output.
// Format: WM_CLASS(STRING) = "instance", "ClassName"
func parseWMClass(raw string) string {
	idx := strings.Index(raw, "= ")
	if idx < 0 {
		return ""
	}
	parts := strings.Split(raw[idx+2:], ",")
	if len(parts) < 2 {
		return ""
	}
	// Use the class name (second field) — it's more human-readable.
	cls := strings.Trim(strings.TrimSpace(parts[1]), `"`)
	return cls
}
