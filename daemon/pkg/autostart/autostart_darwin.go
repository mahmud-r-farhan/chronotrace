//go:build darwin

package autostart

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

const plistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
  "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.chronotrace.daemon</string>
    <key>ProgramArguments</key>
    <array>
        <string>{{.ExePath}}</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <false/>
    <key>StandardErrorPath</key>
    <string>{{.LogPath}}</string>
    <key>StandardOutPath</key>
    <string>{{.LogPath}}</string>
</dict>
</plist>
`

type darwinManager struct {
	plistPath string
}

func newPlatformManager() Manager {
	home, _ := os.UserHomeDir()
	return &darwinManager{
		plistPath: filepath.Join(home, "Library", "LaunchAgents", "com.chronotrace.daemon.plist"),
	}
}

func (m *darwinManager) Enable(executablePath string) error {
	if err := os.MkdirAll(filepath.Dir(m.plistPath), 0o755); err != nil {
		return fmt.Errorf("autostart: create LaunchAgents dir: %w", err)
	}

	home, _ := os.UserHomeDir()
	data := struct {
		ExePath string
		LogPath string
	}{
		ExePath: executablePath,
		LogPath: filepath.Join(home, "Library", "Logs", "ChronoTrace", "daemon.log"),
	}

	f, err := os.Create(m.plistPath)
	if err != nil {
		return fmt.Errorf("autostart: create plist: %w", err)
	}
	defer f.Close()

	tmpl, _ := template.New("plist").Parse(plistTemplate)
	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("autostart: write plist: %w", err)
	}
	return nil
}

func (m *darwinManager) Disable() error {
	if err := os.Remove(m.plistPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("autostart: remove plist: %w", err)
	}
	return nil
}

func (m *darwinManager) IsEnabled() (bool, error) {
	_, err := os.Stat(m.plistPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}
