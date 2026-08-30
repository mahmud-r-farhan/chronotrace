#!/usr/bin/env bash
set -e

APP_DEST="/Applications/ChronoTrace.app"
DAEMON_DIR="$HOME/Library/Application Support/ChronoTrace/bin"
PLIST_PATH="$HOME/Library/LaunchAgents/com.chronotrace.daemon.plist"

echo "Installing ChronoTrace for macOS..."

mkdir -p "$DAEMON_DIR" "$HOME/Library/LaunchAgents"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Copy daemon
cp -f "$SCRIPT_DIR/chronotrace-daemon" "$DAEMON_DIR/chronotrace-daemon" 2>/dev/null || cp -f "$SCRIPT_DIR/../build/chronotrace-daemon-darwin-arm64" "$DAEMON_DIR/chronotrace-daemon"
chmod +x "$DAEMON_DIR/chronotrace-daemon"

# Copy App if present
if [ -d "$SCRIPT_DIR/ChronoTrace.app" ]; then
    rm -rf "$APP_DEST"
    cp -R "$SCRIPT_DIR/ChronoTrace.app" "/Applications/"
fi

# Create launchd plist
cat <<EOF > "$PLIST_PATH"
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.chronotrace.daemon</string>
    <key>ProgramArguments</key>
    <array>
        <string>$DAEMON_DIR/chronotrace-daemon</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardErrorPath</key>
    <string>/tmp/chronotrace-daemon.err</string>
    <key>StandardOutPath</key>
    <string>/tmp/chronotrace-daemon.log</string>
</dict>
</plist>
EOF

launchctl unload "$PLIST_PATH" 2>/dev/null || true
launchctl load "$PLIST_PATH"

echo "✓ ChronoTrace LaunchAgent loaded and running"
