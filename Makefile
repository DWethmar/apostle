# Go project Makefile

# Project variables
APP_NAME := $(shell basename $(CURDIR))
GO_LINT := $(shell command -v golangci-lint 2> /dev/null)
BUILD_DIR := bin
BIN := $(BUILD_DIR)/$(APP_NAME)

# Build options
MAIN_PKG ?= .         # change to ./cmd/<app> if needed
TAGS ?= release
LDFLAGS ?=

# WASM build settings
WASM_DIR := web
WASM_BIN := $(WASM_DIR)/$(APP_NAME).wasm
WASM_EXEC := $(shell go env GOROOT)/lib/wasm/wasm_exec.js 
INDEX_SRC := index.html
SERVE_PORT ?= 8080

# Default target
.PHONY: all
all: build

# =========================
#      NATIVE BUILD
# =========================

.PHONY: build
build:
	@echo "👉 Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -tags=$(TAGS) -ldflags="$(LDFLAGS)" -o $(BIN) $(MAIN_PKG)

.PHONY: run
run: build
	@echo "🚀 Running $(APP_NAME)..."
	@$(BIN)

# =========================
#      TEST & LINT
# =========================

.PHONY: test
test:
	@echo "🧪 Running tests..."
	@go test ./... -v -coverprofile=coverage.out
	@go tool cover -func=coverage.out | tail -n 1

.PHONY: fmt
fmt:
	@echo "🧹 Formatting code..."
	@go fmt ./...

.PHONY: lint
lint:
ifndef GO_LINT
	$(error "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest")
endif
	@echo "🔍 Linting code..."
	@golangci-lint run

.PHONY: tidy
tidy:
	@echo "🧾 Tidying modules..."
	@go mod tidy

# =========================
#        WASM BUILD
# =========================

.PHONY: wasm
wasm: wasm-prepare
	@echo "🧩 Building WASM → $(WASM_BIN)"
	@GOOS=js GOARCH=wasm go build -tags=$(TAGS) -ldflags="$(LDFLAGS)" -o $(WASM_BIN) $(MAIN_PKG)

.PHONY: wasm-prepare
wasm-prepare:
	@mkdir -p $(WASM_DIR)
	@cp -f $(WASM_EXEC) $(WASM_DIR)/wasm_exec.js
	@if [ -f "$(INDEX_SRC)" ]; then \
		echo "📄 Copying index.html → $(WASM_DIR)/index.html"; \
		cp $(INDEX_SRC) $(WASM_DIR)/index.html; \
	else \
		echo "⚠️  No index.html found at root!"; \
	fi

.PHONY: serve
serve: wasm
	@echo "🌐 Serving $(WASM_DIR) at http://localhost:$(SERVE_PORT)"
	@cd $(WASM_DIR) && python3 -m http.server $(SERVE_PORT)

# =========================
#        HOUSEKEEPING
# =========================

.PHONY: clean
clean:
	@echo "🗑️ Cleaning..."
	@rm -rf $(BUILD_DIR) coverage.out $(WASM_DIR)/$(APP_NAME).wasm

.PHONY: generate
generate:
	@echo "⚙️ Running code generation..."
	@go generate ./...
