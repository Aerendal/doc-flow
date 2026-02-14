GO ?= go
GOFLAGS ?= -mod=vendor
GOCACHE ?= /tmp/go-cache

export GOFLAGS
export GOCACHE

BUILD_DIR ?= build

.PHONY: build test tools lint smoke vendor-sync

build: $(BUILD_DIR)/docflow

$(BUILD_DIR)/docflow:
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/docflow ./cmd/docflow

test:
	$(GO) test ./...

tools:
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/docflow-section-order-lint ./cmd/section-order-lint
	$(GO) build -o $(BUILD_DIR)/docflow-source-template-check ./cmd/source-template-check
	$(GO) build -o $(BUILD_DIR)/worklog-lint-days ./cmd/worklog-lint-days

lint:
	$(GO) vet ./...

smoke: build
	./$(BUILD_DIR)/docflow --help >/dev/null

vendor-sync:
	$(GO) mod tidy
	$(GO) mod vendor
	$(GO) mod verify
