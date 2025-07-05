# Makefile for Golang API REST

# Variables
APP_NAME=golang-api-rest
BINARY_NAME=cmd/api/main.go
DOCKER_IMAGE=golang-api-rest
DOCKER_TAG=latest

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GORUN=$(GOCMD) run
GOLINT=golangci-lint
GOSEC=gosec

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(shell git describe --tags --always --dirty) -X main.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S')"

.PHONY: all build clean test coverage lint security-check fmt vet deps tidy run dev docker-build docker-run migrate-up migrate-down seed help

# Default target
all: clean build

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o bin/$(APP_NAME) $(BINARY_NAME)
	@echo "Build completed: bin/$(APP_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -f bin/$(APP_NAME)
	@echo "Clean completed"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linting
lint:
	@echo "Running linter..."
	$(GOLINT) run ./...

# Run security check
security-check:
	@echo "Running security check..."
	$(GOSEC) ./...

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOGET) -v -t -d ./...

# Tidy modules
tidy:
	@echo "Tidying modules..."
	$(GOMOD) tidy

# Run the application
run:
	@echo "Running $(APP_NAME)..."
	$(GORUN) $(BINARY_NAME)

# Run in development mode
dev:
	@echo "Running in development mode..."
	LOG_LEVEL=debug $(GORUN) $(BINARY_NAME)

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Docker run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

# Run database migrations up
migrate-up:
	@echo "Running database migrations up..."
	$(GORUN) cmd/seeds/main.go migrate up

# Run database migrations down
migrate-down:
	@echo "Running database migrations down..."
	$(GORUN) cmd/seeds/main.go migrate down

# Seed database
seed:
	@echo "Seeding database..."
	$(GORUN) cmd/seeds/main.go seed

# Generate swagger documentation
swagger:
	@echo "Generating swagger documentation..."
	swag init -g $(BINARY_NAME) -o docs

# Install development tools
install-tools:
	@echo "Installing development tools..."
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint
	$(GOGET) -u github.com/securecodewarrior/gosec/v2/cmd/gosec
	$(GOGET) -u github.com/swaggo/swag/cmd/swag

# Pre-commit checks
pre-commit: fmt vet lint security-check test

# CI/CD pipeline
ci: deps tidy fmt vet lint security-check test coverage

# Show help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"
	@echo "  coverage       - Run tests with coverage"
	@echo "  lint           - Run linter"
	@echo "  security-check - Run security check"
	@echo "  fmt            - Format code"
	@echo "  vet            - Run go vet"
	@echo "  deps           - Install dependencies"
	@echo "  tidy           - Tidy modules"
	@echo "  run            - Run the application"
	@echo "  dev            - Run in development mode"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  migrate-up     - Run database migrations up"
	@echo "  migrate-down   - Run database migrations down"
	@echo "  seed           - Seed database"
	@echo "  swagger        - Generate swagger documentation"
	@echo "  install-tools  - Install development tools"
	@echo "  pre-commit     - Run pre-commit checks"
	@echo "  ci             - Run CI/CD pipeline"
	@echo "  help           - Show this help" 