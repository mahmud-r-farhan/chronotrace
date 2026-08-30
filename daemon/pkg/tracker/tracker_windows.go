//go:build windows

package tracker

import (
	"fmt"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// windowsTracker implements Tracker using Win32 APIs via golang.org/x/sys/windows.
type windowsTracker struct{}

func newPlatformTracker() Tracker {
	return &windowsTracker{}
}

// GetActiveWindow retrieves information about the current foreground window.
// It gracefully handles elevated/system processes that deny access.
func (t *windowsTracker) GetActiveWindow() (*ActiveWindowInfo, error) {
	hwnd := windows.GetForegroundWindow()
	if hwnd == 0 {
		return nil, fmt.Errorf("no foreground window")
	}

	// Get window title.
	title := getWindowTitle(hwnd)

	// Get PID from window handle.
	var pid uint32
	_, _ = windows.GetWindowThreadProcessId(hwnd, &pid)
	if pid == 0 {
		return &ActiveWindowInfo{
			AppName:     "Unknown",
			WindowTitle: title,
			PID:         0,
		}, nil
	}

	// Open process with limited-information access (no admin required for most apps).
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		// Return partial info — some system processes deny even limited access.
		return &ActiveWindowInfo{
			AppName:     "System",
			WindowTitle: title,
			PID:         pid,
		}, nil
	}
	defer windows.CloseHandle(handle)

	// Query the full path of the executable.
	var size uint32 = windows.MAX_PATH
	buf := make([]uint16, size)
	err = windows.QueryFullProcessImageName(handle, 0, &buf[0], &size)
	if err != nil {
		return &ActiveWindowInfo{
			AppName:     "Unknown",
			WindowTitle: title,
			PID:         pid,
		}, nil
	}

	exePath := windows.UTF16ToString(buf[:size])
	appName := exeToAppName(exePath)

	return &ActiveWindowInfo{
		AppName:     appName,
		WindowTitle: title,
		ExePath:     exePath,
		PID:         pid,
	}, nil
}

// getWindowTitle safely retrieves the title text of a window handle.
func getWindowTitle(hwnd windows.HWND) string {
	// GetWindowTextLength + GetWindowText via raw syscall (not wrapped in x/sys).
	user32 := windows.NewLazySystemDLL("user32.dll")
	getWindowTextW := user32.NewProc("GetWindowTextW")
	getWindowTextLengthW := user32.NewProc("GetWindowTextLengthW")

	length, _, _ := getWindowTextLengthW.Call(uintptr(hwnd))
	if length == 0 {
		return ""
	}
	buf := make([]uint16, length+1)
	n, _, _ := getWindowTextW.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(length+1),
	)
	if n == 0 {
		return ""
	}
	return syscall.UTF16ToString(buf)
}

// exeToAppName converts a full executable path to a friendly app name.
// e.g. "C:\Program Files\Google\Chrome\Application\chrome.exe" → "chrome"
func exeToAppName(exePath string) string {
	base := filepath.Base(exePath)
	// Strip .exe extension.
	name := strings.TrimSuffix(base, filepath.Ext(base))
	// Apply well-known display name overrides.
	if friendly, ok := knownApps[strings.ToLower(name)]; ok {
		return friendly
	}
	return name
}

// knownApps maps executable names to human-friendly display names.
var knownApps = map[string]string{
	"chrome":          "Google Chrome",
	"msedge":          "Microsoft Edge",
	"firefox":         "Firefox",
	"code":            "VS Code",
	"windowsterminal": "Windows Terminal",
	"explorer":        "File Explorer",
	"notepad":         "Notepad",
	"slack":           "Slack",
	"discord":         "Discord",
	"zoom":            "Zoom",
	"outlook":         "Microsoft Outlook",
	"winword":         "Microsoft Word",
	"excel":           "Microsoft Excel",
	"powerpnt":        "Microsoft PowerPoint",
	"teams":           "Microsoft Teams",
	"rider64":         "JetBrains Rider",
	"idea64":          "IntelliJ IDEA",
	"pycharm64":       "PyCharm",
	"goland64":        "GoLand",
	"spotify":         "Spotify",
	"obsidian":        "Obsidian",
	"figma":           "Figma",
	"postman":         "Postman",
}
