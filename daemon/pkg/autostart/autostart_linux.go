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
ExecStart={{.ExePath}}
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

func (m *linuxManager) Enable(executablePath string) error {
	if err := os.MkdirAll(filepath.Dir(m.unitPath), 0o755); err != nil {
		return fmt.Errorf("autostart: create systemd user dir: %w", err)
	}

	f, err := os.Create(m.unitPath)
	if err != nil {
		return fmt.Errorf("autostart: create unit file: %w", err)
	}
	defer f.Close()

	data := struct{ ExePath string }{ExePath: executablePath}
	tmpl, _ := template.New("unit").Parse(systemdUnitTemplate)
	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("autostart: write unit file: %w", err)
	}

	// Enable and start the service.
	if err := exec.Command("systemctl", "--user", "daemon-reload").Run(); err != nil {
		return fmt.Errorf("autostart: daemon-reload: %w", err)
	}
	if err := exec.Command("systemctl", "--user", "enable", "chronotrace-daemon.service").Run(); err != nil {
		return fmt.Errorf("autostart: enable service: %w", err)
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
