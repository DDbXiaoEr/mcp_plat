.PHONY: build clean run-card run-academic run-network run-library run-all

BIN_DIR := bin
SERVICES := card academic network library

GOOS ?= linux
GOARCH ?= amd64

build:
	@mkdir -p $(BIN_DIR)
	@for svc in $(SERVICES); do \
		echo "Building $$svc-service ($(GOOS)/$(GOARCH))..."; \
		GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BIN_DIR)/$$svc-service ./cmd/$$svc/; \
	done
	@echo "All services built successfully."

build-local:
	@mkdir -p $(BIN_DIR)
	@for svc in $(SERVICES); do \
		echo "Building $$svc-service (local)..."; \
		go build -o $(BIN_DIR)/$$svc-service ./cmd/$$svc/; \
	done
	@echo "All services built successfully."

clean:
	rm -rf $(BIN_DIR)

run-card:
	go run ./cmd/card/ -config configs/card.yaml

run-academic:
	go run ./cmd/academic/ -config configs/academic.yaml

run-network:
	go run ./cmd/network/ -config configs/network.yaml

run-library:
	go run ./cmd/library/ -config configs/library.yaml

run-all:
	@echo "Starting all campus MCP services..."
	@go run ./cmd/card/ -config configs/card.yaml &
	@go run ./cmd/academic/ -config configs/academic.yaml &
	@go run ./cmd/network/ -config configs/network.yaml &
	@go run ./cmd/library/ -config configs/library.yaml &
	@echo "All services started. Press Ctrl+C to stop."
	@wait

lint:
	go vet ./...

test:
	go test ./...

tidy:
	go mod tidy
