//go:build linux

package autostart

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const systemdUnitTemplate = `[Unit]
Description=ChronoTrace Background Daemon
After=graphical-session.target
PartOf=graphical-session.target

[Service]
Type=simple
ExecStart="{{.ExePath}}"
Restart=on-failure
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=graphical-session.target
`

type linuxManager struct {
	unitPath string
}

func newPlatformManager() Manager {
	cfgHome := os.Getenv("XDG_CONFIG_HOME")
	if cfgHome == "" {
		home, _ := os.UserHomeDir()
		cfgHome = filepath.Join(home, ".config")
	}
	return &linuxManager{
		unitPath: filepath.Join(cfgHome, "systemd", "user", "chronotrace-daemon.service"),
	}
}

// Enable writes the systemd user unit and enables the service.
// The executable path must be absolute — systemd rejects relative ExecStart paths.
func (m *linuxManager) Enable(executablePath string) error {
	if !filepath.IsAbs(executablePath) {
		return fmt.Errorf("autostart: executable path must be absolute, got %q", executablePath)
	}
	if err := os.MkdirAll(filepath.Dir(m.unitPath), 0o755); err != nil {
		return fmt.Errorf("autostart: create systemd user dir: %w", err)
	}

	tmpl, err := template.New("unit").Parse(systemdUnitTemplate)
	if err != nil {
		return fmt.Errorf("autostart: parse unit template: %w", err)
	}

	f, err := os.Create(m.unitPath)
	if err != nil {
		return fmt.Errorf("autostart: create unit file: %w", err)
	}

	data := struct{ ExePath string }{ExePath: executablePath}
	if err := tmpl.Execute(f, data); err != nil {
		_ = f.Close()
		return fmt.Errorf("autostart: write unit file: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("autostart: close unit file: %w", err)
	}

	// Enable the service. These commands need a running systemd user session
	// (D-Bus); when one is not available the unit file remains installed and
	// can be enabled later, so report that clearly.
	if err := exec.Command("systemctl", "--user", "daemon-reload").Run(); err != nil {
		return fmt.Errorf("autostart: unit installed but systemctl daemon-reload failed (is a systemd user session available?): %w", err)
	}
	if err := exec.Command("systemctl", "--user", "enable", "chronotrace-daemon.service").Run(); err != nil {
		return fmt.Errorf("autostart: unit installed but systemctl enable failed: %w", err)
	}
	return nil
}

func (m *linuxManager) Disable() error {
	_ = exec.Command("systemctl", "--user", "disable", "--now", "chronotrace-daemon.service").Run()
	if err := os.Remove(m.unitPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("autostart: remove unit file: %w", err)
	}
	_ = exec.Command("systemctl", "--user", "daemon-reload").Run()
	return nil
}

func (m *linuxManager) IsEnabled() (bool, error) {
	_, err := os.Stat(m.unitPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}
