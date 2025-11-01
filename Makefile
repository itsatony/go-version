.PHONY: all build test lint clean help
.PHONY: test-unit test-integration test-race test-coverage test-bench
.PHONY: build-cli build-all docker-build

# Version information
VERSION := $(shell cat VERSION)
GIT_COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "dev")
GIT_TAG := $(shell git describe --tags --always 2>/dev/null || echo "v0.0.0")
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
BUILD_USER := $(shell whoami)

# Build flags
LDFLAGS := -s -w
LDFLAGS += -X github.com/itsatony/go-version.GitCommit=$(GIT_COMMIT)
LDFLAGS += -X github.com/itsatony/go-version.GitTag=$(GIT_TAG)
LDFLAGS += -X github.com/itsatony/go-version.BuildTime=$(BUILD_TIME)
LDFLAGS += -X github.com/itsatony/go-version.BuildUser=$(BUILD_USER)

# Default target
all: lint test build

help:
	@echo "go-version Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  all                Run lint, test, and build (default)"
	@echo "  build              Build example applications"
	@echo "  build-cli          Build CLI tool"
	@echo "  build-all          Build all binaries"
	@echo "  test               Run all tests (unit + integration + race)"
	@echo "  test-unit          Run unit tests only"
	@echo "  test-integration   Run integration tests only"
	@echo "  test-race          Run race detector"
	@echo "  test-coverage      Generate coverage report"
	@echo "  test-bench         Run benchmarks"
	@echo "  lint               Run linters"
	@echo "  clean              Remove build artifacts"
	@echo "  docker-build       Build Docker image"
	@echo ""
	@echo "Version: $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Git Tag: $(GIT_TAG)"

# Build targets
build:
	@echo "Building version $(VERSION)..."
	@mkdir -p bin
	@# Will build examples once they exist

build-cli:
	@echo "Building CLI tool..."
	@mkdir -p bin
	@go build -ldflags "$(LDFLAGS)" -o bin/go-version ./cmd/go-version

build-all: build build-cli

# Test targets
test: test-unit test-integration test-race
	@echo "All tests passed!"

test-unit:
	@echo "Running unit tests..."
	@go test -v -count=1 ./... -short

test-integration:
	@echo "Running integration tests..."
	@go test -v -count=1 ./... -run Integration

test-race:
	@echo "Running race detection..."
	@go test -race -count=1 ./...

test-coverage:
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage: $$(go tool cover -func=coverage.out | grep total | awk '{print $$3}')"
	@echo "Report saved to coverage.html"

test-bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Lint target
lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found. Install: https://golangci-lint.run/usage/install/"; \
		echo "Running basic go vet instead..."; \
		go vet ./...; \
	fi

# Docker target
docker-build:
	@echo "Building Docker image..."
	@docker build -t go-version:$(VERSION) .
	@docker tag go-version:$(VERSION) go-version:latest

# Clean target
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/ coverage.out coverage.html
	@go clean -testcache
