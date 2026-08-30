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

// GetActiveWindow retrieves the active application on macOS using AppleScript.
// AppleScript is used instead of full CGO to keep CGO_ENABLED=0 compatible builds.
func (t *darwinTracker) GetActiveWindow() (*ActiveWindowInfo, error) {
	// Get frontmost app name
	appOut, err := runAppleScript(`tell application "System Events"
		set frontApp to first application process whose frontmost is true
		set appName to name of frontApp
		set appPath to POSIX path of (application file of frontApp as alias)
		return appName & "|" & appPath
	end tell`)
	if err != nil {
		// Fallback: simpler script without path
		appOut, err = runAppleScript(`tell application "System Events" to get name of first application process whose frontmost is true`)
		if err != nil {
			return &ActiveWindowInfo{AppName: "Unknown"}, nil
		}
		return &ActiveWindowInfo{AppName: strings.TrimSpace(appOut)}, nil
	}

	parts := strings.SplitN(strings.TrimSpace(appOut), "|", 2)
	info := &ActiveWindowInfo{}
	if len(parts) >= 1 {
		info.AppName = parts[0]
	}
	if len(parts) >= 2 {
		info.ExePath = parts[1]
		// Refine app name from path if available
		if info.AppName == "" {
			info.AppName = strings.TrimSuffix(filepath.Base(info.ExePath), ".app")
		}
	}

	// Get window title of frontmost window
	titleOut, err := runAppleScript(`tell application "System Events"
		set frontApp to first application process whose frontmost is true
		if exists (front window of frontApp) then
			return title of front window of frontApp
		end if
		return ""
	end tell`)
	if err == nil {
		info.WindowTitle = strings.TrimSpace(titleOut)
	}

	return info, nil
}

// runAppleScript executes an AppleScript one-liner via osascript.
func runAppleScript(script string) (string, error) {
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
