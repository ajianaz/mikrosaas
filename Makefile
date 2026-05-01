# Makefile for MikroSaaS
# Usage: make <target>

BINARY    := bin/server
GO        := go
GOFLAGS   := -trimpath -ldflags "-s -w -X github.com/ajianaz/mikrosaas.version=$(shell git describe --tags --always --dirty 2>/dev/null || echo dev)"
BUILD_DIR := bin

.PHONY: all build run test lint clean migrate-up migrate-down help

## all: Build the binary
all: build

## build: Compile the server binary
build:
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BINARY) ./cmd/server

## run: Build and run the server (uses env vars for config)
run: build
	./$(BINARY)

## dev: Run with go run (faster iteration, no binary)
dev:
	$(GO) run ./cmd/server

## test: Run all tests
test:
	$(GO) test -race -count=1 ./...

## test-cover: Run tests with coverage report
test-cover:
	$(GO) test -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## lint: Run golangci-lint
lint:
	@if command -v golangci-lint > /dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## fmt: Format code
fmt:
	$(GO) fmt ./...
	gofmt -s -w .

## vet: Run go vet
vet:
	$(GO) vet ./...

## tidy: Tidy and verify module dependencies
tidy:
	$(GO) mod tidy
	$(GO) mod verify

## clean: Remove build artifacts
clean:
	rm -rf $(BUILD_DIR) coverage.out coverage.html

## migrate-up: Run core migrations (requires DATABASE_URL)
migrate-up:
	@echo "Use: psql $$DATABASE_URL -f migrations/core/001_init.up.sql"

## migrate-down: Rollback core migrations (requires DATABASE_URL)
migrate-down:
	@echo "Use: psql $$DATABASE_URL -f migrations/core/001_init.down.sql"

## help: Show this help message
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ": "}; {printf "\033[36m%-20s\033[0m %s\n", $$2, $$1}' | sed 's/## //'
