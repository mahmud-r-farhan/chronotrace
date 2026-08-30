BINARY_DAEMON := chronotrace-daemon
BINARY_GUI    := ChronoTrace
VERSION       ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS       := -s -w -X main.version=$(VERSION)
BUILD_DIR     := build

.PHONY: all daemon gui run-daemon dev clean install-tools release help

all: daemon gui

## daemon: Build the background daemon for the current OS (headless on Windows)
daemon:
	@mkdir -p $(BUILD_DIR)
ifeq ($(OS),Windows_NT)
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS) -H windowsgui" -o $(BUILD_DIR)/$(BINARY_DAEMON).exe ./daemon/cmd/chronotrace-daemon
else
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_DAEMON) ./daemon/cmd/chronotrace-daemon
endif
	@echo "✓ Daemon built → $(BUILD_DIR)/$(BINARY_DAEMON)"

## gui: Build the Wails GUI for the current OS
gui:
	@cd gui/frontend && npm install --silent && npm run build
ifeq ($(OS),Windows_NT)
	@cd gui && go build -tags "desktop,production" -ldflags="$(LDFLAGS) -H windowsgui" -o ../$(BUILD_DIR)/$(BINARY_GUI).exe .
else
	@cd gui && go build -tags "desktop,production" -ldflags="$(LDFLAGS)" -o ../$(BUILD_DIR)/$(BINARY_GUI) .
endif
	@echo "✓ GUI built → $(BUILD_DIR)/$(BINARY_GUI)"

## run-daemon: Run the daemon directly (foreground)
run-daemon:
	go run ./daemon/cmd/chronotrace-daemon

## dev: Start Wails dev server with hot reload
dev:
	cd gui && wails dev

## autostart-install: Register daemon in OS autostart
autostart-install: daemon
	$(BUILD_DIR)/$(BINARY_DAEMON) --autostart-install

## autostart-remove: Remove daemon from OS autostart
autostart-remove:
	$(BUILD_DIR)/$(BINARY_DAEMON) --autostart-remove

## release: Cross-compile daemon for all 5 targets
release:
	@mkdir -p $(BUILD_DIR)
	@echo "Building daemon for all platforms..."
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="$(LDFLAGS) -H windowsgui" -o $(BUILD_DIR)/$(BINARY_DAEMON)-windows-amd64.exe ./daemon/cmd/chronotrace-daemon
	GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_DAEMON)-linux-amd64    ./daemon/cmd/chronotrace-daemon
	GOOS=linux   GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_DAEMON)-linux-arm64    ./daemon/cmd/chronotrace-daemon
	GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_DAEMON)-darwin-amd64   ./daemon/cmd/chronotrace-daemon
	GOOS=darwin  GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_DAEMON)-darwin-arm64   ./daemon/cmd/chronotrace-daemon
	@echo "✓ All daemon binaries built in $(BUILD_DIR)/"

## vet: Run go vet
vet:
	cd daemon && go vet ./...
	cd gui && go vet ./...

## clean: Remove build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -rf gui/frontend/dist gui/frontend/node_modules gui/build

## install-tools: Install Wails CLI
install-tools:
	go install github.com/wailsapp/wails/v2/cmd/wails@latest

## help: Show this help
help:
	@grep -E '^## ' Makefile | sed 's/## /  /'
