.PHONY: test test-verbose test-coverage test-unit test-controllers test-watch clean-test build

# Default: Run all tests
test:
	@echo "Running all tests..."
	go test ./... -timeout 30s

# Verbose test output
test-verbose:
	@echo "Running tests with verbose output..."
	go test ./... -v -timeout 30s

# Test with coverage report
test-coverage:
	@echo "Generating coverage report..."
	go test ./... -cover -coverprofile=coverage.out -timeout 30s
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Unit tests only (models, middleware)
test-unit:
	@echo "Running unit tests..."
	go test ./models ./middleware ./validation -v -timeout 30s

# Controller tests (integration with mocks)
test-controllers:
	@echo "Running controller tests..."
	go test ./controllers -v -timeout 60s

# Watch mode - run tests on file changes
test-watch:
	@echo "Watching for changes and running tests..."
	go test ./... -v --watch

# Clean coverage files
clean-test:
	@echo "Cleaning test artifacts..."
	rm -f coverage.out coverage.html

# Build the backend
build:
	@echo "Building backend..."
	go build -o legionstore-backend .
	@echo "Build complete: legionstore-backend"

# Run tests with short timeout (quick feedback)
test-short:
	go test ./... -short -timeout 10s

# Full test suite with benchmarks
test-bench:
	go test ./... -v -bench=. -timeout 60s

# List all test files
list-tests:
	@find . -name "*_test.go" -type f | sort

# Help command
help:
	@echo "Available commands:"
	@echo "  make test              - Run all tests"
	@echo "  make test-verbose      - Run tests with verbose output"
	@echo "  make test-coverage     - Run tests with coverage report"
	@echo "  make test-unit         - Run unit tests only"
	@echo "  make test-controllers  - Run controller tests"
	@echo "  make test-watch        - Run tests in watch mode"
	@echo "  make test-bench        - Run tests with benchmarks"
	@echo "  make build             - Build the backend binary"
	@echo "  make clean-test        - Clean coverage files"
	@echo "  make list-tests        - List all test files"
	@echo "  make help              - Show this help message"
