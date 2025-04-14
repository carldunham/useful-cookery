# Makefile for useful-cookery

# Go parameters
GOCMD=go
GOGEN=$(GOCMD) generate
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOTOOL=$(GOCMD) tool
GOMOD=$(GOCMD) mod
GORUN=$(GOCMD) run

# Project parameters
MODULE=github.com/carldunham/useful-cookery
MAIN_PATH=./cmd/api
BINARY_NAME=useful-cookery-api
GENERATED_DIR=./internal/graphql/generated
MODELS_GEN=./internal/graphql/models/models_gen.go

# Tools
GQLGEN=github.com/99designs/gqlgen
GOLANGCI_LINT=github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: all build clean test lint generate tidy help

all: generate lint test build

# Build the application
build:
	@echo "Building..."
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f $(MODELS_GEN)
	rm -rf $(GENERATED_DIR)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) ./...

# Run linting
lint:
	@echo "Running linter..."
	$(GOTOOL) $(GOLANGCI_LINT) run ./...

# Generate code as needed
generate:
	@echo "Generating Go code..."
	$(GOGEN)

# Tidy up dependencies
tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

# Run the application
run:
	@echo "Running application..."
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)
	./$(BINARY_NAME)

# Show help
help:
	@echo "Makefile for useful-cookery"
	@echo ""
	@echo "Usage:"
	@echo "  make              Run generate, lint, test, and build"
	@echo "  make build        Build the application"
	@echo "  make clean        Clean build artifacts"
	@echo "  make test         Run tests"
	@echo "  make lint         Run linter"
	@echo "  make generate     Generate code with gqlgen"
	@echo "  make tidy         Tidy up dependencies"
	@echo "  make run          Run the application"
	@echo "  make help         Show this help message"
