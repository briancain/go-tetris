.PHONY: build clean run test test-verbose test-coverage

# Binary name
BINARY_NAME=tetris
# Output directory
BIN_DIR=bin
# Coverage output directory
COVERAGE_DIR=coverage

# Build the application
build:
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

# Default target
all: build
