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

// user32 procs are resolved once; GetActiveWindow runs every 2-3 seconds, so
// re-creating the lazy DLL bindings on each poll would be wasteful.
var (
	user32                = windows.NewLazySystemDLL("user32.dll")
	procGetWindowTextW    = user32.NewProc("GetWindowTextW")
	procGetWindowTextLenW = user32.NewProc("GetWindowTextLengthW")
)

// maxPathLen is the Windows long-path limit. MAX_PATH (260) would fail for
// executables installed under deeply nested directories.
const maxPathLen = 32768

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

	// Query the full path of the executable (long-path safe buffer).
	var size uint32 = maxPathLen
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
	length, _, _ := procGetWindowTextLenW.Call(uintptr(hwnd))
	if length == 0 {
		return ""
	}
	buf := make([]uint16, length+1)
	n, _, _ := procGetWindowTextW.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(length+1),
	)
	if n == 0 {
		return ""
	}
	return syscall.UTF16ToString(buf[:n])
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
