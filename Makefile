.PHONY: build build-server build-embed build-embed-windows build-embed-linux-arm64 build-embed-all build-web run clean

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -s -w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME)

SERVER_OUT := mcp_plat-console
WEB_DIR := web

build: build-web build-embed

build-server:
	@echo "==> building server version $(VERSION)"
	go build -ldflags "$(LDFLAGS)" -o $(SERVER_OUT) .

build-embed: build-web
	@echo "==> building standalone (embed web) version $(VERSION)"
	go build -tags embed -ldflags "$(LDFLAGS)" -o $(SERVER_OUT) .

build-embed-windows: build-web
	@echo "==> building standalone (embed web) windows-amd64 version $(VERSION)"
	GOOS=windows GOARCH=amd64 go build -tags embed -ldflags "$(LDFLAGS)" -o $(SERVER_OUT)-windows-amd64.exe .

build-embed-linux-arm64: build-web
	@echo "==> building standalone (embed web) linux-arm64 version $(VERSION)"
	GOOS=linux GOARCH=arm64 go build -tags embed -ldflags "$(LDFLAGS)" -o $(SERVER_OUT)-linux-arm64 .

build-embed-all: build-embed build-embed-windows build-embed-linux-arm64

build-web:
	@echo "==> building web"
	cd $(WEB_DIR) && npm ci && npm run build

run:
	go run .

dev:
	@echo "==> starting dev server on :8080"
	go run . &
	@echo "==> starting dev web on :5174"
	cd $(WEB_DIR) && npm run dev

clean:
	rm -f $(SERVER_OUT) $(SERVER_OUT)-windows-amd64.exe $(SERVER_OUT)-linux-arm64
	rm -rf $(WEB_DIR)/dist
