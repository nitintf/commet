# Makefile for Commet - AI-powered commit message generator

.PHONY: build clean install test lint fmt help run dev

# Default target
help:
	@echo "Available commands:"
	@echo "  build    - Build the commet binary"
	@echo "  clean    - Remove build artifacts"
	@echo "  install  - Install commet to GOPATH/bin"
	@echo "  test     - Run tests"
	@echo "  lint     - Run linter"
	@echo "  fmt      - Format code"
	@echo "  run      - Run commet with arguments (use: make run ARGS='--help')"
	@echo "  dev      - Build and run in development mode"

# Build the binary
build:
	@echo "Building commet..."
	go build -o bin/commet .

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Install to GOPATH/bin
install: build
	@echo "Installing commet..."
	go install .

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run with arguments (use: make run ARGS='--help')
run: build
	@echo "Running commet..."
	./bin/commet $(ARGS)

# Development mode - build and run config
dev: build
	@echo "Running in development mode..."
	./bin/commet config set

# Create bin directory if it doesn't exist
bin:
	mkdir -p bin

# Ensure bin directory exists before building
build: bin