.PHONY: build-local build-local-release build-linux-amd64-release build-linux-arm64-release \
        build-server build-embed \
        build-accesskey-auth-server build-accesskey-auth-server-release \
        build-audit-log-server build-audit-log-server-release \
        build-server-linux-amd64 build-embed-linux-amd64 \
        build-accesskey-auth-server-linux-amd64 build-audit-log-server-linux-amd64 \
        build-embed-linux-arm64 \
        build-accesskey-auth-server-linux-arm64 build-audit-log-server-linux-arm64 \
        build-tools build-auditgen build-apisix-plugin build-web \
        proto run dev datagen auditgen clean

VERSION := $(shell git describe --tags --always 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

XFLAGS := -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME)
DEBUG_LDFLAGS := $(XFLAGS)
RELEASE_LDFLAGS := -s -w $(XFLAGS) -X 'github.com/gin-gonic/gin.mode=release'

BUILD_DIR := build
BIN_DIR := $(BUILD_DIR)/bin
SERVER_OUT := $(BIN_DIR)/mcp_plat-console
EMBED_OUT := $(BIN_DIR)/mcp_plat
ACCESSKEY_SERVER_OUT := $(BIN_DIR)/accesskey-auth-server
AUDIT_LOG_SERVER_OUT := $(BIN_DIR)/audit-log-server
APISIX_RUNNER_OUT := $(BIN_DIR)/apisix-go-runner
TOOLS_DIR := $(BIN_DIR)/tools
ARCH ?= amd64
WEB_DIR := web

# ------------------------------------------------------------------
# build groups
# ------------------------------------------------------------------

build-local: build-server build-accesskey-auth-server build-audit-log-server

build-local-release: build-embed build-accesskey-auth-server-release build-audit-log-server-release

build-linux-amd64-release: build-server-linux-amd64 build-embed-linux-amd64 \
	build-accesskey-auth-server-linux-amd64 build-audit-log-server-linux-amd64

build-linux-arm64-release: build-embed-linux-arm64 \
	build-accesskey-auth-server-linux-arm64 build-audit-log-server-linux-arm64

# ------------------------------------------------------------------
# build-local: debug (gin=debug)
# ------------------------------------------------------------------

build-server:
	@echo "==> building server (local debug, gin=debug) $(VERSION)"
	@mkdir -p $(BIN_DIR)
	go build -ldflags "$(DEBUG_LDFLAGS)" -o $(SERVER_OUT) .

build-accesskey-auth-server: proto
	@echo "==> building accesskey gRPC server (local debug) $(VERSION)"
	@mkdir -p $(BIN_DIR)
	go build -ldflags "$(DEBUG_LDFLAGS)" -o $(ACCESSKEY_SERVER_OUT) ./cmd/accesskey-auth-server/

build-audit-log-server: proto
	@echo "==> building audit log gRPC server (local debug) $(VERSION)"
	@mkdir -p $(BIN_DIR)
	go build -ldflags "$(DEBUG_LDFLAGS)" -o $(AUDIT_LOG_SERVER_OUT) ./cmd/audit-log-server/

# ------------------------------------------------------------------
# build-local-release: static link -s -w (gin=release)
# ------------------------------------------------------------------

build-embed: build-web
	@echo "==> building standalone embed web (local release, gin=release) $(VERSION)"
	@mkdir -p $(BIN_DIR)
	go build -tags embed -ldflags "$(RELEASE_LDFLAGS)" -o $(EMBED_OUT) .

build-accesskey-auth-server-release: proto
	@echo "==> building accesskey gRPC server (local release) $(VERSION)"
	@mkdir -p $(BIN_DIR)
	go build -ldflags "$(RELEASE_LDFLAGS)" -o $(ACCESSKEY_SERVER_OUT) ./cmd/accesskey-auth-server/

build-audit-log-server-release: proto
	@echo "==> building audit log gRPC server (local release) $(VERSION)"
	@mkdir -p $(BIN_DIR)
	go build -ldflags "$(RELEASE_LDFLAGS)" -o $(AUDIT_LOG_SERVER_OUT) ./cmd/audit-log-server/

# ------------------------------------------------------------------
# build-linux-amd64-release: static link -s -w (gin=release)
# ------------------------------------------------------------------

build-server-linux-amd64:
	@echo "==> building server (linux/amd64 release, gin=release) $(VERSION)"
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags "$(RELEASE_LDFLAGS)" -o $(SERVER_OUT)-linux-amd64 .

build-embed-linux-amd64: build-web
	@echo "==> building standalone embed web (linux/amd64 release) $(VERSION)"
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 go build -tags embed -ldflags "$(RELEASE_LDFLAGS)" -o $(EMBED_OUT)-linux-amd64 .

build-accesskey-auth-server-linux-amd64: proto
	@echo "==> building accesskey gRPC server (linux/amd64 release) $(VERSION)"
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags "$(RELEASE_LDFLAGS)" -o $(ACCESSKEY_SERVER_OUT)-linux-amd64 ./cmd/accesskey-auth-server/

build-audit-log-server-linux-amd64: proto
	@echo "==> building audit log gRPC server (linux/amd64 release) $(VERSION)"
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags "$(RELEASE_LDFLAGS)" -o $(AUDIT_LOG_SERVER_OUT)-linux-amd64 ./cmd/audit-log-server/

# ------------------------------------------------------------------
# build-linux-arm64-release: static link -s -w (gin=release)
# ------------------------------------------------------------------

build-embed-linux-arm64: build-web
	@echo "==> building standalone embed web (linux/arm64 release) $(VERSION)"
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=arm64 go build -tags embed -ldflags "$(RELEASE_LDFLAGS)" -o $(EMBED_OUT)-linux-arm64 .

build-accesskey-auth-server-linux-arm64: proto
	@echo "==> building accesskey gRPC server (linux/arm64 release) $(VERSION)"
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=arm64 go build -ldflags "$(RELEASE_LDFLAGS)" -o $(ACCESSKEY_SERVER_OUT)-linux-arm64 ./cmd/accesskey-auth-server/

build-audit-log-server-linux-arm64: proto
	@echo "==> building audit log gRPC server (linux/arm64 release) $(VERSION)"
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=arm64 go build -ldflags "$(RELEASE_LDFLAGS)" -o $(AUDIT_LOG_SERVER_OUT)-linux-arm64 ./cmd/audit-log-server/

# ------------------------------------------------------------------
# other build targets
# ------------------------------------------------------------------

build-tools:
	@echo "==> building tools"
	@mkdir -p $(TOOLS_DIR)
	go build -o $(TOOLS_DIR)/datagen cmd/datagen/main.go
	go build -o $(TOOLS_DIR)/accesskey-test cmd/accesskey-test/main.go
	go build -o $(TOOLS_DIR)/audit-log-server cmd/audit-log-server/main.go
	go build -o $(TOOLS_DIR)/auditgen cmd/auditgen/main.go

build-apisix-plugin: proto
	@echo "==> building APISIX go plugin runner (linux/$(ARCH))"
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=$(ARCH) go build -ldflags "$(RELEASE_LDFLAGS)" -o $(APISIX_RUNNER_OUT)-linux-$(ARCH) ./cmd/apisix-runner/

build-web:
	@echo "==> building web"
	$(MAKE) -C $(WEB_DIR) build

proto:
	@echo "==> generating protobuf code"
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative plugin/accesskey.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative plugin/auditlog.proto

# ------------------------------------------------------------------
# utils
# ------------------------------------------------------------------

run:
	go run .

dev:
	@echo "==> starting dev server on :8080"
	go run . &
	@echo "==> starting dev web on :5174"
	cd $(WEB_DIR) && npm run dev

datagen:
	@go run cmd/datagen/main.go -table $(TABLE) -count $(or $(COUNT),1)

auditgen:
	@go run ./cmd/auditgen -days $(or $(DAYS),30) -per-day $(or $(PER_DAY),200)

build-auditgen:
	@echo "==> building auditgen"
	@mkdir -p $(TOOLS_DIR)
	go build -o $(TOOLS_DIR)/auditgen cmd/auditgen/main.go

clean:
	rm -rf $(BIN_DIR)
	rm -f plugin/accesskey.pb.go plugin/accesskey_grpc.pb.go plugin/auditlog.pb.go plugin/auditlog_grpc.pb.go
	$(MAKE) -C $(WEB_DIR) clean
