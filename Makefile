.PHONY: build build-server build-embed build-web run clean

VERSION := $(shell date -u '+%Y%m%d%H%M%S')
LDFLAGS := -s -w -X main.Version=$(VERSION) -X main.BuildTime=$(VERSION)

SERVER_OUT := mcp_plat-console
WEB_DIR := web

build: build-web build-embed

build-server:
	@echo "==> building server version $(VERSION)"
	go build -ldflags "$(LDFLAGS)" -o $(SERVER_OUT) .

build-embed: build-web
	@echo "==> building standalone (embed web) version $(VERSION)"
	go build -tags embed -ldflags "$(LDFLAGS)" -o $(SERVER_OUT) .

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
	rm -f $(SERVER_OUT)
	rm -rf $(WEB_DIR)/dist
