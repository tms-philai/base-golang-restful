# Makefile for base-gin project

# Variables
BINARY_NAME=base-gin
BINARY_UNIX=$(BINARY_NAME)_unix
MAIN_PATH=./cmd/server
TEST_PATH=./test
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(MAIN_PATH)

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(COVERAGE_FILE)
	rm -f $(COVERAGE_HTML)

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)
	./$(BINARY_NAME)

# Run with Air (hot reload for development)
dev:
	air

# Generate Swagger documentation
docs:
	swag init -g $(MAIN_PATH)/main.go --parseDependency --parseInternal

# Generate Swagger docs and run with Air
dev-docs:
	swag init -g $(MAIN_PATH)/main.go --parseDependency --parseInternal
	air

# Run tests
test:
	$(GOTEST) -v ./...

# Run unit tests only
test-unit:
	$(GOTEST) -v ./test/unit/...

# Run integration tests only
test-integration:
	$(GOTEST) -v ./test/integration/...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=$(COVERAGE_FILE) ./...
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)

# Run unit tests with coverage
test-unit-coverage:
	$(GOTEST) -v -coverprofile=$(COVERAGE_FILE) ./test/unit/...
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)

# Run integration tests with coverage
test-integration-coverage:
	$(GOTEST) -v -coverprofile=$(COVERAGE_FILE) ./test/integration/...
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)

# Run tests with race detection
test-race:
	$(GOTEST) -v -race ./...

# Run benchmarks
benchmark:
	$(GOTEST) -v -bench=. ./...

# Format code
fmt:
	$(GOFMT) ./...

# Vet code
vet:
	$(GOVET) ./...

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Install dependencies
install:
	$(GOGET) -d ./...

# Generate mocks (if using mockgen)
generate:
	$(GOCMD) generate ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Lint and fix
lint-fix:
	golangci-lint run --fix

# Security scan (requires gosec)
security:
	gosec ./...

# Run all quality checks
quality: fmt vet lint security

# Run all tests (unit + integration)
test-all: test-unit test-integration

# Run all tests with coverage
test-all-coverage: test-unit-coverage test-integration-coverage

# Setup test environment
test-setup:
	@echo "Setting up test environment..."
	@mkdir -p ./configs
	@mkdir -p ./tmp
	@echo "Test environment setup complete"

# Clean test environment
test-clean:
	@echo "Cleaning test environment..."
	@rm -rf ./configs/.env.test
	@rm -rf ./tmp/test_files
	@echo "Test environment cleaned"

# Database setup for tests
test-db-setup:
	@echo "Setting up test database..."
	@echo "Please ensure PostgreSQL is running and create test database:"
	@echo "CREATE DATABASE test_base_gin;"
	@echo "CREATE USER test_user WITH PASSWORD 'test_password';"
	@echo "GRANT ALL PRIVILEGES ON DATABASE test_base_gin TO test_user;"

# Run specific test package
test-package:
	@read -p "Enter package path (e.g., ./test/unit): " package; \
	$(GOTEST) -v $$package

# Run specific test function
test-function:
	@read -p "Enter test function (e.g., TestAuthHandler_Register): " function; \
	$(GOTEST) -v -run $$function ./...

# Show test coverage in terminal
test-coverage-term:
	$(GOTEST) -v -coverprofile=$(COVERAGE_FILE) ./...
	$(GOCMD) tool cover -func=$(COVERAGE_FILE)

# Run tests with verbose output and show test names
test-verbose:
	$(GOTEST) -v -count=1 ./...

# Run tests with timeout
test-timeout:
	$(GOTEST) -v -timeout=30s ./...

# Run tests in parallel
test-parallel:
	$(GOTEST) -v -parallel=4 ./...

# Generate test report
test-report:
	$(GOTEST) -v -json ./... > test-results.json

# Docker commands
docker-build:
	docker build -t $(BINARY_NAME) .

docker-run:
	docker run -p 8001:8001 $(BINARY_NAME)

# Development commands
dev-setup: deps test-setup test-db-setup
	@echo "Development environment setup complete"

dev-test: test-setup test-all
	@echo "All tests completed"

dev-clean: clean test-clean
	@echo "Development environment cleaned"

# CI/CD commands
ci-test: deps test-setup test-all-coverage
	@echo "CI tests completed"

ci-build: deps build
	@echo "CI build completed"

# Help
help:
	@echo "Available commands:"
	@echo "  build              - Build the application"
	@echo "  build-linux        - Build for Linux"
	@echo "  clean              - Clean build artifacts"
	@echo "  run                - Run the application"
	@echo "  dev                - Run with Air (hot reload)"
	@echo "  docs               - Generate Swagger documentation"
	@echo "  dev-docs           - Generate docs and run with Air"
	@echo "  test               - Run all tests"
	@echo "  test-unit          - Run unit tests only"
	@echo "  test-integration   - Run integration tests only"
	@echo "  test-coverage      - Run tests with coverage"
	@echo "  test-unit-coverage - Run unit tests with coverage"
	@echo "  test-integration-coverage - Run integration tests with coverage"
	@echo "  test-race          - Run tests with race detection"
	@echo "  benchmark          - Run benchmarks"
	@echo "  fmt                - Format code"
	@echo "  vet                - Vet code"
	@echo "  deps               - Download dependencies"
	@echo "  install            - Install dependencies"
	@echo "  generate           - Generate mocks"
	@echo "  lint               - Lint code"
	@echo "  lint-fix           - Lint and fix"
	@echo "  security           - Security scan"
	@echo "  quality            - Run all quality checks"
	@echo "  test-all           - Run all tests (unit + integration)"
	@echo "  test-all-coverage  - Run all tests with coverage"
	@echo "  test-setup         - Setup test environment"
	@echo "  test-clean         - Clean test environment"
	@echo "  test-db-setup      - Setup test database"
	@echo "  test-package       - Run specific test package"
	@echo "  test-function      - Run specific test function"
	@echo "  test-coverage-term - Show test coverage in terminal"
	@echo "  test-verbose       - Run tests with verbose output"
	@echo "  test-timeout       - Run tests with timeout"
	@echo "  test-parallel      - Run tests in parallel"
	@echo "  test-report        - Generate test report"
	@echo "  docker-build       - Build Docker image"
	@echo "  docker-run         - Run Docker container"
	@echo "  dev-setup          - Setup development environment"
	@echo "  dev-test           - Run development tests"
	@echo "  dev-clean          - Clean development environment"
	@echo "  ci-test            - Run CI tests"
	@echo "  ci-build           - Run CI build"
	@echo "  help               - Show this help message"

# Default target
.DEFAULT_GOAL := help

# Phony targets
.PHONY: build build-linux clean run test test-unit test-integration test-coverage test-unit-coverage test-integration-coverage test-race benchmark fmt vet deps install generate lint lint-fix security quality test-all test-all-coverage test-setup test-clean test-db-setup test-package test-function test-coverage-term test-verbose test-timeout test-parallel test-report docker-build docker-run dev-setup dev-test dev-clean ci-test ci-build help
