.PHONY: help build run dev test clean kill deps fmt lint check

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
build: ## Build the main application
	@echo "Building application..."
	@go build -o prj-start .

build-all: ## Build all projects in dev-docs directory
	@echo "Building all projects..."
	@for dir in dev-docs/*/; do \
		if [ -f "$$dir/go.mod" ]; then \
			echo "Building $$dir"; \
			cd "$$dir" && go build -o main . && cd ../..; \
		fi \
	done

# Run targets
run: build ## Build and run the application
	@echo "Running application..."
	@./prj-start

run-dev: ## Run application without building (for development)
	@echo "Running application in development mode..."
	@go run .

# Development targets
dev: ## Run with Air hot reload (requires Air to be installed)
	@echo "Starting development server with hot reload..."
	@which air > /dev/null || (echo "Air not installed. Install with: go install github.com/cosmtrek/air@latest" && exit 1)
	@air

dev-all: ## Run all projects with Air (for testing multiple projects)
	@echo "Starting all projects with Air..."
	@for dir in dev-docs/*/; do \
		if [ -f "$$dir/go.mod" ] && [ -f "$$dir/.air.toml" -o -f "$$dir/.air.linux.conf" -o -f "$$dir/.air.windows.conf" ]; then \
			echo "Starting $$dir with Air"; \
			cd "$$dir" && air & \
			cd ../..; \
		fi \
	done

# Test targets
test: ## Run all tests
	@echo "Running tests..."
	@go test ./...

test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	@go test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-all: ## Run tests for all projects
	@echo "Running tests for all projects..."
	@for dir in dev-docs/*/; do \
		if [ -f "$$dir/go.mod" ]; then \
			echo "Testing $$dir"; \
			cd "$$dir" && go test ./... && cd ../..; \
		fi \
	done

# Clean targets
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -f prj-start prj-start.exe coverage.out coverage.html
	@go clean -cache

clean-all: ## Clean all projects
	@echo "Cleaning all projects..."
	@for dir in dev-docs/*/; do \
		if [ -f "$$dir/go.mod" ]; then \
			echo "Cleaning $$dir"; \
			cd "$$dir" && rm -f main main.exe && go clean -cache && cd ../..; \
		fi \
	done
	@rm -f prj-start prj-start.exe coverage.out coverage.html

# Process management
kill: ## Kill running Go processes
	@echo "Killing running Go processes..."
	@pkill -f "go run" || true
	@pkill -f "prj-start" || true
	@pkill -f "air" || true
	@echo "Processes killed"

kill-port: ## Kill process using specific port (usage: make kill-port PORT=3000)
	@if [ -z "$(PORT)" ]; then \
		echo "Usage: make kill-port PORT=3000"; \
		exit 1; \
	fi
	@echo "Killing process on port $(PORT)..."
	@lsof -ti:$(PORT) | xargs kill -9 || true

# Dependency management
deps: ## Download and tidy dependencies
	@echo "Managing dependencies..."
	@go mod tidy
	@go mod download
	@go mod verify

deps-all: ## Manage dependencies for all projects
	@echo "Managing dependencies for all projects..."
	@for dir in dev-docs/*/; do \
		if [ -f "$$dir/go.mod" ]; then \
			echo "Managing dependencies for $$dir"; \
			cd "$$dir" && go mod tidy && go mod download && go mod verify && cd ../..; \
		fi \
	done

# Code quality
fmt: ## Format Go code
	@echo "Formatting code..."
	@go fmt ./...

fmt-all: ## Format code for all projects
	@echo "Formatting code for all projects..."
	@for dir in dev-docs/*/; do \
		if [ -f "$$dir/go.mod" ]; then \
			echo "Formatting $$dir"; \
			cd "$$dir" && go fmt ./... && cd ../..; \
		fi \
	done

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	@golangci-lint run

check: fmt lint test ## Run all code quality checks

check-all: ## Run checks for all projects
	@echo "Running checks for all projects..."
	@for dir in dev-docs/*/; do \
		if [ -f "$$dir/go.mod" ]; then \
			echo "Checking $$dir"; \
			cd "$$dir" && go fmt ./... && go test ./... && cd ../..; \
		fi \
	done

# Setup targets
setup: ## Setup development environment
	@echo "Setting up development environment..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Development tools installed"

setup-colors: ## Install colored logging packages
	@echo "Installing colored logging packages..."
	@go get github.com/sirupsen/logrus
	@go get github.com/fatih/color
	@go get github.com/sirupsen/logrus
	@echo "Colored logging packages installed"

# Utility targets
version: ## Show Go version
	@go version

list-projects: ## List all Go projects in dev-docs
	@echo "Go projects in dev-docs:"
	@for dir in dev-docs/*/; do \
		if [ -f "$$dir/go.mod" ]; then \
			echo "  $$dir"; \
		fi \
	done

# Default port for development
PORT ?= 3000

# Environment variables for development
export GO_ENV=development