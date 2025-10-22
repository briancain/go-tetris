.PHONY: build clean run test test-verbose test-coverage mod-tidy mod-tidy-check lint fmt fmt-check help build-windows build-macos build-macos-arm64 build-all

# Binary name
BINARY_NAME=tetris
# Output directory
BIN_DIR=bin
# Coverage output directory
COVERAGE_DIR=coverage

# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

# Help target
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "  ${YELLOW}%-20s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
	@echo ''

# Build all applications
build: build-desktop build-web build-server

# Build the desktop application
build-desktop:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY_NAME) ./cmd

# Build the server application
build-server:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/server ./cmd/server

# Build for web (WebAssembly)
build-web:
	mkdir -p $(BIN_DIR)/web
	GOOS=js GOARCH=wasm go build -o $(BIN_DIR)/web/$(BINARY_NAME).wasm ./cmd
	cp $$(go env GOROOT)/lib/wasm/wasm_exec.js $(BIN_DIR)/web/
	cp web/index.html $(BIN_DIR)/web/

# Clean build artifacts
clean:
	rm -f $(BIN_DIR)/$(BINARY_NAME)
	rm -rf $(BIN_DIR)/web
	rm -rf $(COVERAGE_DIR)

# Run the application
run: build-desktop
	./$(BIN_DIR)/$(BINARY_NAME)

# Run the web application
run-web: build-web
	@echo "Starting web server at http://localhost:8000"
	cd $(BIN_DIR)/web && python3 -m http.server 8000

# Run the server application
run-server: build-server
	./$(BIN_DIR)/server

# Run all tests
test:
	go test ./...

# Run server tests only
test-server:
	go test ./internal/server/...

# Run integration tests (requires server to be running)
test-integration:
	go test ./test/integration -v

# Run manual test client
test-manual:
	@echo "Starting manual test client..."
	@echo "Make sure server is running with 'make run-server'"
	go run test/manual/manual_client.go $(USER)

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Run tests in short mode (skips tests that require a display)
test-short:
	go test -short ./...

# Run tests in short mode with verbose output
test-short-verbose:
	go test -v -short ./...

# Run tests with coverage report
test-coverage:
	mkdir -p $(COVERAGE_DIR)
	go test -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated at $(COVERAGE_DIR)/coverage.html"

# Run go mod tidy
mod-tidy:
	go mod tidy

# Check if go mod tidy would make changes
mod-tidy-check:
	@echo "Checking if go mod tidy would make changes..."
	@go mod tidy
	@echo "No changes needed."

# Format code
fmt:
	gofmt -w .
	goimports -w .

# Check formatting
fmt-check:
	@echo "Checking code formatting..."
	@if [ "$$(gofmt -l . | wc -l)" -gt 0 ]; then \
		echo "The following files are not formatted correctly:"; \
		gofmt -l .; \
		exit 1; \
	fi
	@if [ "$$(goimports -l . | wc -l)" -gt 0 ]; then \
		echo "The following files have import formatting issues:"; \
		goimports -l .; \
		exit 1; \
	fi
	@echo "All files are properly formatted."

# Run linter
lint:
	golangci-lint run

# Verify all checks
verify: fmt-check mod-tidy-check lint test

# Default target
all: build

# Cross-platform build targets
.PHONY: build-windows build-macos build-macos-arm64 build-all
# Cross-compilation build tags:
# - headless: Our custom implementation for CI builds
# - ebitennogl: Disables OpenGL dependencies in Ebiten
# - ebitennonscreen: Disables screen-related functionality in Ebiten

# Build for Windows
build-windows:
	mkdir -p $(BIN_DIR)
	# Use pure Go mode for cross-compilation
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags=headless,ebitennogl,ebitennonscreen -o $(BIN_DIR)/$(BINARY_NAME).exe ./cmd

# Build for macOS (amd64)
build-macos:
	mkdir -p $(BIN_DIR)
	# Use pure Go mode for cross-compilation
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -tags=headless,ebitennogl,ebitennonscreen -o $(BIN_DIR)/$(BINARY_NAME)_darwin_amd64 ./cmd

# Build for macOS (arm64)
build-macos-arm64:
	mkdir -p $(BIN_DIR)
	# Use pure Go mode for cross-compilation
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -tags=headless,ebitennogl,ebitennonscreen -o $(BIN_DIR)/$(BINARY_NAME)_darwin_arm64 ./cmd

# Build for all platforms
build-all: build-desktop build-web build-server
	@echo "All builds completed successfully!"
