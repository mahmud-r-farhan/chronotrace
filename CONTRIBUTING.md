# Contributing to ChronoTrace

Thanks for your interest in improving ChronoTrace! This document explains how to
set up a development environment and what we look for in contributions.

## Getting Started

**Prerequisites:** Go 1.22+, Node.js 18+ (GUI frontend only), and the Wails v2
CLI if you want to build the desktop app.

```bash
git clone https://github.com/mahmud-r-farhan/chronotrace.git
cd chronotrace

make install-tools   # installs the Wails CLI
make daemon          # builds the headless daemon into ./build
make gui             # builds the desktop GUI (runs npm install/build first)
```

### Useful Make Targets

| Command     | Description                                        |
|-------------|----------------------------------------------------|
| `make daemon` | Build the background daemon for the current OS   |
| `make gui`    | Build the Wails GUI                              |
| `make dev`    | Run the GUI with hot reload (`wails dev`)        |
| `make run-daemon` | Run the daemon in the foreground               |
| `make test`   | Run all Go tests (daemon + GUI modules)          |
| `make vet`    | Run `go vet` on both modules                     |
| `make fmt`    | Fail if any Go source is not gofmt-formatted     |
| `make release`| Cross-compile the daemon for all platforms       |

## Development Workflow

1. Fork the repository and create a feature branch from `main`
   (`git checkout -b feature/amazing-feature`).
2. Make your changes. Please keep the daemon lightweight — the project's core
   promise is < 15 MB RAM and ~0% idle CPU.
3. Verify your work before opening a pull request:

   ```bash
   make fmt && make vet && make test
   ```

4. Commit using [Conventional Commits](https://www.conventionalcommits.org/)
   where possible (`feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`).
5. Push your branch and open a Pull Request against `main`.

## Architecture Notes

- **`daemon/`** — pure Go, `CGO_ENABLED=0`. Platform-specific code lives in
  `_windows.go` / `_darwin.go` / `_linux.go` files behind build tags.
  - `pkg/tracker` — active-window detection per OS.
  - `pkg/storage` — SQLite (pure Go driver) with a batched write buffer.
    Timestamps are stored as **local** wall-clock strings; queries must use
    SQLite's `localtime` modifier when comparing against `now`.
  - `pkg/ipc` — loopback-only REST API on `127.0.0.1:42069`.
  - `pkg/autostart` — Registry (Windows) / LaunchAgent (macOS) / systemd user
    unit (Linux) registration.
- **`gui/`** — Wails v2 app. The Go backend (`app.go`) proxies the daemon's
  REST API; the frontend is framework-free JS in `gui/frontend/src`.
- **`website/`** — Next.js marketing site (independent of the app).

## Ground Rules

- **Privacy is non-negotiable.** No telemetry, no external network calls from
  the daemon or GUI. Data never leaves the user's device.
- Keep dependencies minimal; the daemon must stay CGO-free.
- Escape untrusted strings (window titles, app names) before injecting them
  into HTML in the frontend.
- Add or update tests for behaviour changes — CI runs `gofmt`, `go vet`, and
  `go test` on every push.

## Reporting Issues

Please include your OS and version, ChronoTrace version
(`chronotrace-daemon --version`), and steps to reproduce. For daemon problems,
logs are written to the system journal (Linux), Console (macOS), or the
service log (Windows).

## License

By contributing, you agree that your contributions will be licensed under the
[MIT License](LICENSE).
