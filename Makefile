.PHONY: build build-server build-embed build-embed-windows build-embed-linux-arm64 build-embed-all build-web run clean build-tool datagen build-accesskey-auth-server build-apisix-runner proto

VERSION := $(shell git describe --tags --always 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -s -w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME)

SERVER_OUT := mcp_plat-console
EMBED_OUT := mcp_plat
ACCESSKEY_SERVER_OUT := accesskey-auth-server
APISIX_RUNNER_OUT := apisix-go-runner
ARCH ?= amd64
WEB_DIR := web

build: build-web build-embed

build-server:
	@echo "==> building server version $(VERSION)"
	go build -ldflags "$(LDFLAGS)" -o $(SERVER_OUT) .

build-accesskey-auth-server: proto
	@echo "==> building accesskey gRPC server (linux/$(ARCH))"
	GOOS=linux GOARCH=$(ARCH) go build -ldflags "$(LDFLAGS)" -o $(ACCESSKEY_SERVER_OUT)-linux-$(ARCH) ./cmd/accesskey-auth-server/

build-apisix-runner: proto
	@echo "==> building APISIX go plugin runner (linux/$(ARCH))"
	GOOS=linux GOARCH=$(ARCH) go build -ldflags "$(LDFLAGS)" -o $(APISIX_RUNNER_OUT)-linux-$(ARCH) ./cmd/apisix-runner/

proto:
	@echo "==> generating protobuf code"
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative plugin/accesskey.proto

build-embed: build-web
	@echo "==> building standalone (embed web) version $(VERSION)"
	go build -tags embed -ldflags "$(LDFLAGS)" -o $(EMBED_OUT) .

build-embed-windows: build-web
	@echo "==> building standalone (embed web) windows-amd64 version $(VERSION)"
	GOOS=windows GOARCH=amd64 go build -tags embed -ldflags "$(LDFLAGS)" -o $(EMBED_OUT)-windows-amd64.exe .

build-embed-linux-arm64: build-web
	@echo "==> building standalone (embed web) linux-arm64 version $(VERSION)"
	GOOS=linux GOARCH=arm64 go build -tags embed -ldflags "$(LDFLAGS)" -o $(EMBED_OUT)-linux-arm64 .

build-embed-all: build-embed build-embed-windows build-embed-linux-arm64

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
	@mkdir -p tools
	go build -o tools/datagen cmd/datagen/main.go

datagen:
	@go run cmd/datagen/main.go -table $(TABLE) -count $(or $(COUNT),1)

clean:
	rm -f $(SERVER_OUT) $(EMBED_OUT) $(EMBED_OUT)-windows-amd64.exe $(EMBED_OUT)-linux-arm64
	rm -f $(ACCESSKEY_SERVER_OUT)-linux-* $(APISIX_RUNNER_OUT)-linux-*
	rm -f plugin/accesskey.pb.go plugin/accesskey_grpc.pb.go
	rm -rf tools/
	$(MAKE) -C $(WEB_DIR) clean
