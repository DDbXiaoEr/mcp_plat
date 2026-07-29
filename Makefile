.PHONY: build build-server build-embed build-embed-windows build-embed-linux-arm64 build-embed-all build-web run clean build-tool datagen build-accesskey-auth-server build-accesskey-auth-server-linux-amd64 build-apisix-runner proto releasebuild

VERSION := $(shell git describe --tags --always 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -s -w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME)

BUILD_DIR := build
BIN_DIR := $(BUILD_DIR)/bin
SERVER_OUT := $(BIN_DIR)/mcp_plat-console
EMBED_OUT := $(BIN_DIR)/mcp_plat
ACCESSKEY_SERVER_OUT := $(BIN_DIR)/accesskey-auth-server
APISIX_RUNNER_OUT := $(BIN_DIR)/apisix-go-runner
TOOLS_DIR := $(BIN_DIR)/tools
ARCH ?= amd64
WEB_DIR := web

build: build-web build-embed

build-server:
	@echo "==> building server version $(VERSION)"
	@mkdir -p $(BIN_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(SERVER_OUT) .

build-accesskey-auth-server: proto
	@echo "==> building accesskey gRPC server (current platform)"
	@mkdir -p $(BIN_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(ACCESSKEY_SERVER_OUT) ./cmd/accesskey-auth-server/

build-accesskey-auth-server-linux-amd64: proto
	@echo "==> building accesskey gRPC server (linux/amd64)"
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(ACCESSKEY_SERVER_OUT)-linux-amd64 ./cmd/accesskey-auth-server/

build-apisix-runner: proto
	@echo "==> building APISIX go plugin runner (linux/$(ARCH))"
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=$(ARCH) go build -ldflags "$(LDFLAGS)" -o $(APISIX_RUNNER_OUT)-linux-$(ARCH) ./cmd/apisix-runner/

proto:
	@echo "==> generating protobuf code"
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative plugin/accesskey.proto

build-embed: build-web
	@echo "==> building standalone (embed web) version $(VERSION)"
	@mkdir -p $(BIN_DIR)
	go build -tags embed -ldflags "$(LDFLAGS)" -o $(EMBED_OUT) .

build-embed-windows: build-web
	@echo "==> building standalone (embed web) windows-amd64 version $(VERSION)"
	@mkdir -p $(BIN_DIR)
	GOOS=windows GOARCH=amd64 go build -tags embed -ldflags "$(LDFLAGS)" -o $(EMBED_OUT)-windows-amd64.exe .

build-embed-linux-arm64: build-web
	@echo "==> building standalone (embed web) linux-arm64 version $(VERSION)"
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=arm64 go build -tags embed -ldflags "$(LDFLAGS)" -o $(EMBED_OUT)-linux-arm64 .

build-embed-all: build-embed build-embed-windows build-embed-linux-arm64

releasebuild: build-web
	@echo "==> building release (embed web, gin release mode) version $(VERSION)"
	@mkdir -p $(BIN_DIR)
	go build -tags embed -ldflags "$(LDFLAGS) -X 'github.com/gin-gonic/gin.mode=release'" -o $(EMBED_OUT) .

build-web:
	@echo "==> building web"
	$(MAKE) -C $(WEB_DIR) build

run:
	go run .

dev:
	@echo "==> starting dev server on :8080"
	go run . &
	@echo "==> starting dev web on :5174"
	cd $(WEB_DIR) && npm run dev

build-tool:
	@echo "==> building tools"
	@mkdir -p $(TOOLS_DIR)
	go build -o $(TOOLS_DIR)/datagen cmd/datagen/main.go
	go build -o $(TOOLS_DIR)/accesskey-test cmd/accesskey-test/main.go

datagen:
	@go run cmd/datagen/main.go -table $(TABLE) -count $(or $(COUNT),1)

clean:
	rm -rf $(BIN_DIR)
	rm -f plugin/accesskey.pb.go plugin/accesskey_grpc.pb.go
	$(MAKE) -C $(WEB_DIR) clean
