.PHONY: build clean run

# Binary name
BINARY_NAME=tetris
# Output directory
BIN_DIR=bin

# Build the application
build:
	go build -o $(BIN_DIR)/$(BINARY_NAME) ./cmd

# Clean build artifacts
clean:
	rm -f $(BIN_DIR)/$(BINARY_NAME)

# Run the application
run: build
	./$(BIN_DIR)/$(BINARY_NAME)

# Default target
all: build
