# Makefile for TTS Reader
# A Go application for text-to-speech conversion

# Variables
BINARY_NAME=tts-reader
BINARY_PATH=build/$(BINARY_NAME)
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe
BINARY_DARWIN=$(BINARY_NAME)_darwin

# Go related variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(shell command -v gotest 2>/dev/null || echo "$(GOCMD) test")
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-X main.VersionInfo=$(shell git branch --show-current)-$(shell git rev-parse --short HEAD)-$(shell date +%Y%m%d-%H%M)"

# Default target
.DEFAULT_GOAL := build

# Build the application
build: vendor fmt lint test
	@mkdir -p build
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH) ./cmd/tts-reader

# Build for current platform with race detection
build-race:
	@mkdir -p build
	$(GOBUILD) $(LDFLAGS) -race -o $(BINARY_PATH) ./cmd/tts-reader

# Build for multiple platforms
build-all: build-linux build-windows build-darwin

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o build/$(BINARY_UNIX) ./cmd/tts-reader

# Build for Windows
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o build/$(BINARY_WINDOWS) ./cmd/tts-reader

# Build for macOS
build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o build/$(BINARY_DARWIN) ./cmd/tts-reader

pub: build
	cp -f build/$(BINARY_NAME) ~/icloud/bin/

# Run the application
run:
	$(GOCMD) run ./cmd/tts-reader

# Run with verbose output
run-verbose:
	$(GOCMD) run ./cmd/tts-reader -v

# Run with specific language
run-fr:
	$(GOCMD) run ./cmd/tts-reader -l fr "Bonjour le monde"

run-jp:
	$(GOCMD) run ./cmd/tts-reader -l jp "こんにちは世界"

run-pl:
	$(GOCMD) run ./cmd/tts-reader -l pl "Witaj świecie"

# Test the application
test: vendor
	$(GOTEST) -v ./...

# Test with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Test with race detection
test-race:
	$(GOTEST) -race -v ./...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f build/$(BINARY_NAME)
	rm -f build/$(BINARY_UNIX)
	rm -f build/$(BINARY_WINDOWS)
	rm -f build/$(BINARY_DARWIN)
	rm -f coverage.out
	rm -f coverage.html

# Install dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Update dependencies
deps-update:
	$(GOMOD) get -u ./...
	$(GOMOD) tidy

# Format code
fmt: vendor
	$(GOCMD) fmt ./...

# Run linter (requires golangci-lint)
lint: vendor
	golangci-lint run

# Install golangci-lint if not present
install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Check for security vulnerabilities
security:
	gosec ./...

# Install gosec if not present
install-security:
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# Development helpers
dev-setup: deps install-lint install-security

# Quick development cycle
dev: fmt lint test build

# Show help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  build-race     - Build with race detection"
	@echo "  build-all      - Build for Linux, Windows, and macOS"
	@echo "  run            - Run the application"
	@echo "  run-verbose    - Run with verbose output"
	@echo "  run-fr         - Run with French example"
	@echo "  run-jp         - Run with Japanese example"
	@echo "  run-pl         - Run with Polish example"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-race      - Run tests with race detection"
	@echo "  clean          - Clean build artifacts"
	@echo "  deps           - Install dependencies"
	@echo "  deps-update    - Update dependencies"
	@echo "  fmt            - Format code"
	@echo "  lint           - Run linter"
	@echo "  security       - Check for security vulnerabilities"
	@echo "  dev-setup      - Setup development environment"
	@echo "  dev            - Format, lint, test, and build"
	@echo "  help           - Show this help message"

.PHONY: build build-race build-all build-linux build-windows build-darwin \
        run run-verbose run-fr run-jp run-pl test test-coverage test-race \
        clean deps deps-update fmt lint install-lint security install-security \
        dev-setup dev help vendor

vendor:
	go mod vendor
