.PHONY: build-local build-local-release build-linux-amd64-release build-linux-arm64-release \
        build-server build-embed \
        build-accesskey-auth-server build-accesskey-auth-server-release \
        build-audit-log-server build-audit-log-server-release \
        build-server-linux-amd64 build-embed-linux-amd64 \
        build-accesskey-auth-server-linux-amd64 build-audit-log-server-linux-amd64 \
        build-embed-linux-arm64 build-server-linux-arm64 \
        build-accesskey-auth-server-linux-arm64 build-audit-log-server-linux-arm64 \
        build-tools build-auditgen build-apisix-plugin build-web \
        docker-build docker-build-accesskey-auth-server docker-build-audit-log-server docker-build-mcp-plat-embed \
        docker-build-mcp-plat-server docker-build-arm64-mcp-plat-server \
        docker-build-mcp_plat_apisix \
        docker-build-arm64 docker-build-arm64-accesskey-auth-server docker-build-arm64-audit-log-server docker-build-arm64-mcp-plat-embed \
        makedockerbuild docker-build-env \
        docker-clean \
        proto run dev datagen auditgen clean

VERSION := $(shell git describe --tags --always 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

XFLAGS := -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME)
DEBUG_LDFLAGS := $(XFLAGS)
RELEASE_LDFLAGS := -s -w $(XFLAGS) -X 'github.com/gin-gonic/gin.mode=release'

BUILD_DIR := build
BIN_DIR := $(BUILD_DIR)/bin
SERVER_OUT := $(BIN_DIR)/mcp_plat-server
EMBED_OUT := $(BIN_DIR)/mcp_plat_embed
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

build-linux-arm64-release: build-embed-linux-arm64 build-server-linux-arm64 \
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

build-server-linux-arm64:
	@echo "==> building server (linux/arm64 release, gin=release) $(VERSION)"
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=arm64 go build -ldflags "$(RELEASE_LDFLAGS)" -o $(SERVER_OUT)-linux-arm64 .

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
# docker
# ------------------------------------------------------------------

DOCKER_TAG ?= $(GIT_COMMIT)

docker-build: build-linux-amd64-release docker-build-accesskey-auth-server docker-build-audit-log-server docker-build-mcp-plat-embed

docker-build-accesskey-auth-server:
	@echo "==> building docker image accesskey-auth-server:$(DOCKER_TAG)"
	DOCKER_BUILDKIT=0 docker build -f docker/Dockerfile.accesskey-auth-server -t accesskey-auth-server:$(DOCKER_TAG) .

docker-build-audit-log-server:
	@echo "==> building docker image audit-log-server:$(DOCKER_TAG)"
	DOCKER_BUILDKIT=0 docker build -f docker/Dockerfile.audit-log-server -t audit-log-server:$(DOCKER_TAG) .

docker-build-mcp-plat-embed:
	@echo "==> building docker image mcp_plat_embed:$(DOCKER_TAG)"
	DOCKER_BUILDKIT=0 docker build -f docker/Dockerfile.mcp_plat_embed -t mcp_plat_embed:$(DOCKER_TAG) .

docker-build-mcp-plat-server: build-server-linux-amd64
	@echo "==> building docker image mcp_plat_server:$(DOCKER_TAG) (纯后端，不内嵌前端)"
	DOCKER_BUILDKIT=0 docker build -f docker/Dockerfile.mcp_plat_server -t mcp_plat_server:$(DOCKER_TAG) .

docker-build-arm64-mcp-plat-server: build-server-linux-arm64
	@echo "==> building docker image mcp_plat_server:$(DOCKER_TAG) (arm64, 纯后端，不内嵌前端)"
	DOCKER_BUILDKIT=0 docker build -f docker/Dockerfile.mcp_plat_server-arm64 -t mcp_plat_server:$(DOCKER_TAG) .

docker-build-mcp_plat_apisix: build-apisix-plugin
	@echo "==> building docker image mcp_plat_apisix:$(DOCKER_TAG)"
	DOCKER_BUILDKIT=0 docker build -f docker/Dockerfile.mcp_plat_apisix -t mcp_plat_apisix:$(DOCKER_TAG) .

docker-build-arm64: build-linux-arm64-release docker-build-arm64-accesskey-auth-server docker-build-arm64-audit-log-server docker-build-arm64-mcp-plat-embed

docker-build-arm64-accesskey-auth-server:
	@echo "==> building docker image accesskey-auth-server:$(DOCKER_TAG) (arm64)"
	DOCKER_BUILDKIT=0 docker build -f docker/Dockerfile.accesskey-auth-server-arm64 -t accesskey-auth-server:$(DOCKER_TAG) .

docker-build-arm64-audit-log-server:
	@echo "==> building docker image audit-log-server:$(DOCKER_TAG) (arm64)"
	DOCKER_BUILDKIT=0 docker build -f docker/Dockerfile.audit-log-server-arm64 -t audit-log-server:$(DOCKER_TAG) .

docker-build-arm64-mcp-plat-embed:
	@echo "==> building docker image mcp_plat_embed:$(DOCKER_TAG) (arm64)"
	DOCKER_BUILDKIT=0 docker build -f docker/Dockerfile.mcp_plat_embed-arm64 -t mcp_plat_embed:$(DOCKER_TAG) .

# 构建镜像并生成 docker-compose 可用的 .env 文件（仅镜像标签）
makedockerbuild: docker-build docker-build-env

docker-build-env:
	@echo "==> generating docker/.env for docker-compose"
	@mkdir -p docker
	@cat > docker/.env <<-EOF
	# Generated by 'make makedockerbuild' — 供 docker-compose 使用
	MCP_PLAT_EMBED_IMAGE=mcp_plat_embed:$(DOCKER_TAG)
	ACCESSKEY_AUTH_SERVER_IMAGE=accesskey-auth-server:$(DOCKER_TAG)
	AUDIT_LOG_SERVER_IMAGE=audit-log-server:$(DOCKER_TAG)
	MCP_PLAT_APISIX_IMAGE=mcp_plat_apisix:$(DOCKER_TAG)
	EOF
	@echo "==> docker/.env written:"
	@cat docker/.env
	@echo ""
	@echo "==> 提示：已将镜像标签写入 docker/.env，docker-compose 可直接使用。"
	@echo "==> 参考用法（docker-compose.yml 中）："
	@echo "==>   image: \$${MCP_PLAT_EMBED_IMAGE}"
	@echo "==> 启动：docker compose up -d"

docker-clean:
	@echo "==> removing docker images (tag: $(DOCKER_TAG))"
	@docker rmi mcp_plat_embed:$(DOCKER_TAG) 2>/dev/null || true
	@docker rmi mcp_plat_server:$(DOCKER_TAG) 2>/dev/null || true
	@docker rmi accesskey-auth-server:$(DOCKER_TAG) 2>/dev/null || true
	@docker rmi audit-log-server:$(DOCKER_TAG) 2>/dev/null || true
	@docker rmi mcp_plat_apisix:$(DOCKER_TAG) 2>/dev/null || true

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

build-auditdatagenerate:
	@echo "==> building auditgen"
	@mkdir -p $(TOOLS_DIR)
	go build -o $(TOOLS_DIR)/auditgen cmd/auditgen/main.go

clean:
	rm -rf $(BIN_DIR)
	rm -f plugin/accesskey.pb.go plugin/accesskey_grpc.pb.go plugin/auditlog.pb.go plugin/auditlog_grpc.pb.go
	$(MAKE) -C $(WEB_DIR) clean
	$(MAKE) docker-clean
