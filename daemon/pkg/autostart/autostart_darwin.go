//go:build darwin

package autostart

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

// launchAgentLabel is the launchd service identifier for the daemon.
const launchAgentLabel = "com.chronotrace.daemon"

const plistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
  "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>{{.Label}}</string>
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
		plistPath: filepath.Join(home, "Library", "LaunchAgents", launchAgentLabel+".plist"),
	}
}

// Enable writes the LaunchAgent plist and (best-effort) loads it immediately
// so autostart takes effect without requiring a re-login.
func (m *darwinManager) Enable(executablePath string) error {
	if err := os.MkdirAll(filepath.Dir(m.plistPath), 0o755); err != nil {
		return fmt.Errorf("autostart: create LaunchAgents dir: %w", err)
	}

	home, _ := os.UserHomeDir()
	logDir := filepath.Join(home, "Library", "Logs", "ChronoTrace")
	// launchd does not create missing directories for Standard*Path.
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return fmt.Errorf("autostart: create log dir: %w", err)
	}

	data := struct {
		Label   string
		ExePath string
		LogPath string
	}{
		Label:   launchAgentLabel,
		ExePath: xmlEscape(executablePath),
		LogPath: xmlEscape(filepath.Join(logDir, "daemon.log")),
	}

	tmpl, err := template.New("plist").Parse(plistTemplate)
	if err != nil {
		return fmt.Errorf("autostart: parse plist template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("autostart: render plist: %w", err)
	}
	if err := os.WriteFile(m.plistPath, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("autostart: write plist: %w", err)
	}

	// Best-effort: (re)load the agent for the current GUI session. Failure is
	// not fatal — launchd picks the plist up at the next login regardless.
	domain := fmt.Sprintf("gui/%d", os.Getuid())
	_ = exec.Command("launchctl", "bootout", domain+"/"+launchAgentLabel).Run()
	_ = exec.Command("launchctl", "bootstrap", domain, m.plistPath).Run()
	return nil
}

// Disable unloads the LaunchAgent (best-effort) and removes its plist.
func (m *darwinManager) Disable() error {
	domain := fmt.Sprintf("gui/%d", os.Getuid())
	_ = exec.Command("launchctl", "bootout", domain+"/"+launchAgentLabel).Run()
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

// xmlEscape escapes values embedded in the plist XML (&, <, >, quotes) so
// paths containing those characters cannot produce a malformed plist.
func xmlEscape(s string) string {
	var buf bytes.Buffer
	_ = xml.EscapeText(&buf, []byte(s))
	return buf.String()
}
