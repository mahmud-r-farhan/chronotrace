<div align="center">

<img src="https://img.shields.io/github/v/release/mahmud-r-farhan/chronotrace?style=for-the-badge&color=8b5cf6" alt="Latest Release"/>
<img src="https://img.shields.io/github/actions/workflow/status/mahmud-r-farhan/chronotrace/build.yml?style=for-the-badge&color=10b981" alt="Build Status"/>
<img src="https://img.shields.io/badge/license-MIT-blue?style=for-the-badge" alt="License"/>
<img src="https://img.shields.io/badge/RAM-<15MB-purple?style=for-the-badge" alt="RAM Usage"/>
<img src="https://img.shields.io/badge/platforms-Windows%20%7C%20macOS%20%7C%20Linux-informational?style=for-the-badge" alt="Platforms"/>

# ⏱ ChronoTrace

**Privacy-first, ultra-lightweight screen time & app usage tracker.**  
Open-source, cross-platform, zero telemetry. Your data never leaves your device.

</div>

---

## ✨ Features

- 📊 **App Usage Tracking** — Tracks which applications you use and for how long
- 🕐 **Hourly Timeline** — Visual heatmap of your activity throughout the day
- 📅 **Day / Week / Month** — Aggregated views with historical trends
- 🏃 **Ultra-Lightweight Daemon** — < 15 MB RAM, ~0% CPU at idle
- 🔒 **100% Local & Private** — All data stored in a local SQLite database
- 🖥 **Decoupled Architecture** — Daemon runs silently in the background; GUI is optional
- 🚀 **Auto-starts with OS** — Optional native OS autostart (Registry / systemd / LaunchAgent)
- 🌍 **Cross-Platform** — Windows, macOS (Intel & Apple Silicon), Linux (x64 & arm64)

---

## 🏗 Architecture

```
chronotrace-daemon  ──► SQLite DB  ──► REST API (127.0.0.1:42069)
                                                │
chronotrace-gui (Wails)  ◄──────────── HTTP    ┘
```

| Component | Tech | RAM | Startup |
|---|---|---|---|
| `chronotrace-daemon` | Pure Go, no CGO | < 15 MB | Auto (OS login) |
| `chronotrace-gui` | Wails v2 + HTML/CSS/JS | ~80 MB | Manual only |

---

## 📦 1-Click Installation (Recommended)

Download the complete package for your platform from [Releases](https://github.com/mahmud-r-farhan/chronotrace/releases):

| Platform | Package | Setup Method |
|---|---|---|
| 🪟 **Windows (x64)** | [`ChronoTrace-Windows-x64-Setup.zip`](https://github.com/mahmud-r-farhan/chronotrace/releases) | Extract & double-click `install-windows.bat` |
| 🐧 **Linux (x64)** | [`ChronoTrace-Linux-x64.tar.gz`](https://github.com/mahmud-r-farhan/chronotrace/releases) | Extract & run `./install-linux.sh` |
| 🍎 **macOS (Universal)** | [`ChronoTrace-macOS-Universal.zip`](https://github.com/mahmud-r-farhan/chronotrace/releases) | Extract & run `./install-macos.sh` |

*The installer automatically registers the lightweight background tracking daemon to start with your operating system, creates desktop/start menu shortcuts, and launches the application.*

---

### Standalone Headless Daemon Binaries

For servers or users who only want headless tracking via the REST API:
- **Windows**: `chronotrace-daemon-windows-amd64.exe`
- **Linux x64**: `chronotrace-daemon-linux-amd64`
- **Linux arm64**: `chronotrace-daemon-linux-arm64`
- **macOS Intel**: `chronotrace-daemon-darwin-amd64`
- **macOS Apple Silicon**: `chronotrace-daemon-darwin-arm64`

### Option 2 — Build from Source

**Prerequisites:** Go 1.22+, Node.js 18+ (for GUI only), Wails v2

```bash
git clone https://github.com/mahmud-r-farhan/chronotrace.git
cd chronotrace

# Install dependencies
go mod tidy
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Build daemon
make daemon

# Build GUI (requires Wails)
make gui
```

---

## 🚀 Usage

### Starting the Daemon

```bash
# Run the daemon (foreground for testing)
./build/chronotrace-daemon

# Install as OS autostart (runs on every login)
./build/chronotrace-daemon --autostart-install

# Remove autostart
./build/chronotrace-daemon --autostart-remove

# Check version
./build/chronotrace-daemon --version
```

### Opening the GUI

Simply launch `ChronoTrace` from your Applications folder, Start Menu, or:

```bash
./build/bin/ChronoTrace
```

> **Note:** The GUI connects to the daemon via `http://127.0.0.1:42069`. Make sure the daemon is running first.

### REST API (for power users)

The daemon exposes a simple HTTP API:

```bash
# Health check
curl http://127.0.0.1:42069/api/v1/status

# Today's usage
curl http://127.0.0.1:42069/api/v1/usage/today

# This week
curl http://127.0.0.1:42069/api/v1/usage/week

# Hourly timeline for a date
curl "http://127.0.0.1:42069/api/v1/usage/timeline?date=2026-08-31"
```

---

## 🗄 Data Storage

Data is stored locally at:

| Platform | Location |
|---|---|
| Windows | `%APPDATA%\ChronoTrace\data.db` |
| macOS | `~/Library/Application Support/ChronoTrace/data.db` |
| Linux | `~/.local/share/chronotrace/data.db` |

SQLite WAL mode is used for maximum performance with minimal write I/O.

---

## 🛠 Development

```bash
# Run daemon in foreground
make run-daemon

# Wails hot-reload dev server
make dev

# Cross-compile daemon for all platforms
make release

# Run go vet
make vet
```

---

## 🔒 Privacy

- ✅ No network requests ever made to external servers
- ✅ No analytics, no telemetry, no crash reporting
- ✅ IPC is bound strictly to `127.0.0.1` (loopback only)
- ✅ All data lives in a local SQLite file you fully own
- ✅ Open-source — audit every line of code

---

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) and open a PR.

1. Fork the repo
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

[MIT License](LICENSE) — Copyright © 2026 ChronoTrace Contributors
