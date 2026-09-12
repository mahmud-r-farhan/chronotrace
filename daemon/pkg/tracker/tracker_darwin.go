//go:build darwin

package tracker

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// darwinTracker uses AppleScript to get the frontmost application.
type darwinTracker struct{}

func newPlatformTracker() Tracker {
	return &darwinTracker{}
}

// activeWindowScript returns "name|path|title" for the frontmost application
// in a single osascript invocation. Spawning osascript is comparatively
// expensive, so one call per poll (instead of two) matters at 2-3 s intervals.
const activeWindowScript = `tell application "System Events"
	set frontApp to first application process whose frontmost is true
	set appName to name of frontApp
	set appPath to ""
	try
		set appPath to POSIX path of (application file of frontApp as alias)
	end try
	set winTitle to ""
	try
		if exists (front window of frontApp) then
			set winTitle to title of front window of frontApp
		end if
	end try
	return appName & "|" & appPath & "|" & winTitle
end tell`

// GetActiveWindow retrieves the active application on macOS using AppleScript.
// AppleScript is used instead of full CGO to keep CGO_ENABLED=0 compatible builds.
func (t *darwinTracker) GetActiveWindow() (*ActiveWindowInfo, error) {
	out, err := runAppleScript(activeWindowScript)
	if err != nil {
		// Fallback: simpler script without path/title (e.g. when permissions
		// restrict what System Events may query).
		name, ferr := runAppleScript(`tell application "System Events" to get name of first application process whose frontmost is true`)
		if ferr != nil {
			return &ActiveWindowInfo{AppName: "Unknown"}, nil
		}
		return &ActiveWindowInfo{AppName: strings.TrimSpace(name)}, nil
	}

	parts := strings.SplitN(strings.TrimSpace(out), "|", 3)
	info := &ActiveWindowInfo{}
	if len(parts) > 0 {
		info.AppName = strings.TrimSpace(parts[0])
	}
	if len(parts) > 1 {
		info.ExePath = strings.TrimSpace(parts[1])
	}
	if len(parts) > 2 {
		// The title may itself contain "|", so everything after the second
		// separator belongs to it.
		info.WindowTitle = strings.TrimSpace(parts[2])
	}
	// Refine app name from the bundle path when System Events gave none.
	if info.AppName == "" && info.ExePath != "" {
		info.AppName = strings.TrimSuffix(filepath.Base(info.ExePath), ".app")
	}

	return info, nil
}

// runAppleScript executes an AppleScript via osascript.
func runAppleScript(script string) (string, error) {
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
