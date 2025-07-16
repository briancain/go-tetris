.PHONY: build clean run test test-verbose test-coverage mod-tidy mod-tidy-check lint fmt fmt-check

# Binary name
BINARY_NAME=tetris
# Output directory
BIN_DIR=bin
# Coverage output directory
COVERAGE_DIR=coverage

# Build the application
build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY_NAME) ./cmd

# Clean build artifacts
clean:
	rm -f $(BIN_DIR)/$(BINARY_NAME)
	rm -rf $(COVERAGE_DIR)

# Run the application
run: build
	./$(BIN_DIR)/$(BINARY_NAME)

# Run tests
test:
	go test ./...

# Run tests with verbose output
test-verbose:
	go test -v ./...

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
