#!/usr/bin/env bash
set -e

INSTALL_BIN="$HOME/.local/bin"
DESKTOP_DIR="$HOME/.local/share/applications"
ICON_DIR="$HOME/.local/share/icons/hicolor/256x256/apps"
SYSTEMD_DIR="$HOME/.config/systemd/user"

echo "Installing ChronoTrace for Linux..."

mkdir -p "$INSTALL_BIN" "$DESKTOP_DIR" "$ICON_DIR" "$SYSTEMD_DIR"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Copy binaries
cp -f "$SCRIPT_DIR/chronotrace-daemon" "$INSTALL_BIN/chronotrace-daemon" 2>/dev/null || cp -f "$SCRIPT_DIR/../build/chronotrace-daemon-linux-amd64" "$INSTALL_BIN/chronotrace-daemon"
cp -f "$SCRIPT_DIR/ChronoTrace" "$INSTALL_BIN/ChronoTrace" 2>/dev/null || cp -f "$SCRIPT_DIR/../build/ChronoTrace" "$INSTALL_BIN/ChronoTrace" 2>/dev/null || true

chmod +x "$INSTALL_BIN/chronotrace-daemon"
if [ -f "$INSTALL_BIN/ChronoTrace" ]; then
    chmod +x "$INSTALL_BIN/ChronoTrace"
fi

# Copy icon
if [ -f "$SCRIPT_DIR/icon.png" ]; then
    cp -f "$SCRIPT_DIR/icon.png" "$ICON_DIR/chronotrace.png"
elif [ -f "$SCRIPT_DIR/../gui/build/linux/icon.png" ]; then
    cp -f "$SCRIPT_DIR/../gui/build/linux/icon.png" "$ICON_DIR/chronotrace.png"
fi

# Create Desktop entry
cat <<EOF > "$DESKTOP_DIR/chronotrace.desktop"
[Desktop Entry]
Name=ChronoTrace
Comment=Ultra-Lightweight Privacy-First Screen Time Tracker
Exec=$INSTALL_BIN/ChronoTrace
Icon=chronotrace
Terminal=false
Type=Application
Categories=Utility;Clock;
StartupNotify=true
EOF

# Create Systemd User Service
cat <<EOF > "$SYSTEMD_DIR/chronotrace-daemon.service"
[Unit]
Description=ChronoTrace Background Tracker Service
After=default.target

[Service]
Type=simple
ExecStart=$INSTALL_BIN/chronotrace-daemon
Restart=on-failure
RestartSec=5

[Install]
WantedBy=default.target
EOF

# Reload & start systemd service
if command -v systemctl >/dev/null 2>&1; then
    systemctl --user daemon-reload
    systemctl --user enable --now chronotrace-daemon.service
    echo "✓ Systemd user service enabled and started"
fi

echo "✓ ChronoTrace installed successfully to $INSTALL_BIN"
