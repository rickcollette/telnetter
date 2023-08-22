# Makefile for telnetter tests

# Variables
PKG = telnetter

# Phony targets
.PHONY: test clean

# Default target
all: test

# Run tests
test:
	@echo "Running tests for $(PKG)..."
	@go test -v .

# Clean up
clean:
	@echo "Cleaning up..."
	@go clean ./...

