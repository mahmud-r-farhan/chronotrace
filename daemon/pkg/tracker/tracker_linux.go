//go:build linux

package tracker

import (
	"os/exec"
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
	// Get active window ID
	widOut, err := exec.Command("xdotool", "getactivewindow").Output()
	if err != nil {
		return &ActiveWindowInfo{AppName: "Unknown"}, nil
	}
	winID := strings.TrimSpace(string(widOut))

	// Get window name and class in one call
	nameOut, _ := exec.Command("xdotool", "getactivewindow", "getwindowname").Output()
	title := strings.TrimSpace(string(nameOut))

	// Get WM_CLASS to identify app name
	classOut, _ := exec.Command("xprop", "-id", winID, "WM_CLASS").Output()
	appName := parseWMClass(string(classOut))
	if appName == "" {
		appName = title
	}

	// Try to get PID and exe path
	pidOut, _ := exec.Command("xdotool", "getactivewindow", "getwindowpid").Output()
	pidStr := strings.TrimSpace(string(pidOut))

	exePath := ""
	if pidStr != "" {
		exeOut, _ := exec.Command("readlink", "-f", "/proc/"+pidStr+"/exe").Output()
		exePath = strings.TrimSpace(string(exeOut))
	}

	return &ActiveWindowInfo{
		AppName:     appName,
		WindowTitle: title,
		ExePath:     exePath,
	}, nil
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

	// Parse the response — format: (true, '<title>|<class>\n')
	raw := string(out)
	raw = strings.Trim(raw, "(true, ')\n")
	parts := strings.SplitN(raw, "|", 2)

	info := &ActiveWindowInfo{}
	if len(parts) >= 1 {
		info.WindowTitle = strings.TrimSpace(parts[0])
		info.AppName = info.WindowTitle
	}
	if len(parts) >= 2 {
		cls := strings.TrimSpace(parts[1])
		if cls != "" && cls != "unknown" {
			info.AppName = cls
		}
	}
	return info, nil
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
