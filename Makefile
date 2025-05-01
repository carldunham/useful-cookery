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

# Docker parameters
API_IMAGE=useful-cookery-api
UI_IMAGE=useful-cookery-ui
MIGRATIONS_IMAGE=useful-cookery-migrations
IMAGE_TAG=local
K3D_CLUSTER=useful-cookery
K8S_NAMESPACE=useful-cookery-local
API_DEPLOYMENT=useful-cookery-api
UI_DEPLOYMENT=useful-cookery-ui
MIGRATIONS_JOB=useful-cookery-migrations

# Tools
GQLGEN=github.com/99designs/gqlgen
GOLANGCI_LINT=github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: all build build-go build-ui clean test test-go test-ui lint lint-go lint-ui lint-md format format-ui format-md format-check format-check-ui format-check-md fix generate tidy help deploy-api-local deploy-ui-local deploy-migrations-local deploy-local migrate-create migrate-up migrate-down migrate-status

all: generate lint test build

# Build the application
build: build-go build-ui
	@echo "Building..."

# Build Go binaries
build-go:
	@echo "Building Go binaries..."
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

# Build UI
build-ui:
	@echo "Building UI..."
	cd ui && npm run build

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f $(MODELS_GEN)
	rm -rf $(GENERATED_DIR)

# Run tests
test: test-go test-ui
	@echo "Running all tests..."

# Run UI tests
test-ui:
	@echo "Running UI tests..."
	cd ui && npm run test:nowatch

# Run Go tests
test-go:
	@echo "Running Go tests..."
	$(GOTEST) -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# Run linting
lint: lint-go lint-ui lint-md
	@echo "Running project linter..."
	npm run lint

# Run Go linting
lint-go:
	@echo "Running Go linter..."
	$(GOTOOL) $(GOLANGCI_LINT) run ./...

# Run UI linting
lint-ui:
	@echo "Running UI linter..."
	cd ui && npm run lint

# Run Markdown linting
lint-md:
	@echo "Running Markdown linter..."
	npm run lint:md

# Fix linting issues
fix: format
	@echo "Fixing Go linting issues..."
	$(GOTOOL) $(GOLANGCI_LINT) run --fix ./...
	@echo "Fixing UI linting issues..."
	cd ui && npm run lint:fix
	@echo "Fixing Markdown linting issues..."
	npm run lint:md:fix
	@echo "Fixing project linting issues..."
	npm run lint:fix

# Format code
format: format-ui format-md
	@echo "Formatting code..."
	npm run format

# Format UI code
format-ui:
	@echo "Formatting UI code..."
	cd ui && npm run format

# Format Markdown
format-md:
	@echo "Formatting Markdown..."
	npm run format

# Check formatting
format-check: format-check-ui format-check-md
	@echo "Checking code formatting..."
	npx prettier --check "**/*.{js,jsx,ts,tsx,json,yml,yaml}"

# Check UI code formatting
format-check-ui:
	@echo "Checking UI code formatting..."
	cd ui && npx prettier --check "src/**/*.{js,jsx,ts,tsx,json,css,scss,md}"

# Check Markdown formatting
format-check-md:
	@echo "Checking Markdown formatting..."
	npx prettier --check "**/*.md"

# Generate code as needed
generate:
	@echo "Generating Go code..."
	$(GOGEN) ./...
	@echo "Generating Typescript models..."
	( cd ui && npm run generate )

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
	@echo ""
	@echo "Build targets:"
	@echo "  make build        Build the application (Go and UI)"
	@echo "  make build-go     Build Go binaries"
	@echo "  make build-ui     Build UI"
	@echo ""
	@echo "Test targets:"
	@echo "  make test         Run all tests"
	@echo "  make test-go      Run Go tests"
	@echo "  make test-ui      Run UI tests"
	@echo ""
	@echo "Lint targets:"
	@echo "  make lint         Run all linters"
	@echo "  make lint-go      Run Go linting"
	@echo "  make lint-ui      Run UI linting"
	@echo "  make lint-md      Run Markdown linting"
	@echo ""
	@echo "Format targets:"
	@echo "  make format       Format all code"
	@echo "  make format-ui    Format UI code"
	@echo "  make format-md    Format Markdown"
	@echo "  make format-check Check all code formatting"
	@echo "  make format-check-ui Check UI code formatting"
	@echo "  make format-check-md Check Markdown formatting"
	@echo ""
	@echo "Fix targets:"
	@echo "  make fix          Fix linting issues and format code"
	@echo ""
	@echo "Other targets:"
	@echo "  make clean        Clean build artifacts"
	@echo "  make generate     Generate code with gqlgen"
	@echo "  make tidy         Tidy up dependencies"
	@echo "  make run          Run the application"
	@echo ""
	@echo "Deployment targets:"
	@echo "  make deploy-api-local        Build and deploy API locally"
	@echo "  make deploy-ui-local         Build and deploy UI locally"
	@echo "  make deploy-migrations-local Build and deploy migrations locally"
	@echo "  make deploy-local            Build and deploy all components locally"
	@echo ""
	@echo "Migration targets:"
	@echo "  make migrate-create          Create a new database migration"
	@echo "  make migrate-up              Apply database migrations"
	@echo "  make migrate-down            Revert database migrations"
	@echo "  make migrate-status          Show current migration status"
	@echo ""
	@echo "Help:"
	@echo "  make help         Show this help message"

# Build and deploy API locally
deploy-api-local:
	@echo "Building API Docker image..."
	docker build -t $(API_IMAGE):$(IMAGE_TAG) -f Dockerfile.api .
	@echo "Importing API image to k3d..."
	k3d image import $(API_IMAGE):$(IMAGE_TAG) -c $(K3D_CLUSTER)
	@echo "Deploying API to local Kubernetes..."
	kubectl delete job $(MIGRATIONS_JOB) -n $(K8S_NAMESPACE) --ignore-not-found
	kubectl apply -k deploy/kubernetes/overlays/local
	@echo "Restarting API deployment..."
	kubectl rollout restart deployment $(API_DEPLOYMENT) -n $(K8S_NAMESPACE)

# Build and deploy UI locally
deploy-ui-local:
	@echo "Building UI Docker image..."
	docker build -t $(UI_IMAGE):$(IMAGE_TAG) -f Dockerfile.ui .
	@echo "Importing UI image to k3d..."
	k3d image import $(UI_IMAGE):$(IMAGE_TAG) -c $(K3D_CLUSTER)
	@echo "Deploying UI to local Kubernetes..."
	kubectl delete job $(MIGRATIONS_JOB) -n $(K8S_NAMESPACE) --ignore-not-found
	kubectl apply -k deploy/kubernetes/overlays/local
	@echo "Restarting UI deployment..."
	kubectl rollout restart deployment $(UI_DEPLOYMENT) -n $(K8S_NAMESPACE)

# Build and deploy migrations locally
deploy-migrations-local:
	@echo "Building migrations Docker image..."
	docker build -t $(MIGRATIONS_IMAGE):$(IMAGE_TAG) -f Dockerfile.migrations .
	@echo "Importing migrations image to k3d..."
	k3d image import $(MIGRATIONS_IMAGE):$(IMAGE_TAG) -c $(K3D_CLUSTER)
	@echo "Deploying migrations to local Kubernetes..."
	@echo "Running migrations job..."
	kubectl delete job $(MIGRATIONS_JOB) -n $(K8S_NAMESPACE) --ignore-not-found
	kubectl apply -k deploy/kubernetes/overlays/local
	@echo "Waiting for migrations to complete..."
	kubectl wait --for=condition=complete job/$(MIGRATIONS_JOB) -n $(K8S_NAMESPACE) --timeout=60s || true

# Build and deploy all components locally
deploy-local:
	@echo "Building all Docker images..."
	docker build -t $(API_IMAGE):$(IMAGE_TAG) -f Dockerfile.api .
	docker build -t $(UI_IMAGE):$(IMAGE_TAG) -f Dockerfile.ui .
	docker build -t $(MIGRATIONS_IMAGE):$(IMAGE_TAG) -f Dockerfile.migrations .
	@echo "Importing images to k3d..."
	k3d image import $(API_IMAGE):$(IMAGE_TAG) $(UI_IMAGE):$(IMAGE_TAG) $(MIGRATIONS_IMAGE):$(IMAGE_TAG) -c $(K3D_CLUSTER)
	@echo "Deploying to local Kubernetes..."
	@echo "Running migrations job..."
	kubectl delete job $(MIGRATIONS_JOB) -n $(K8S_NAMESPACE) --ignore-not-found
	kubectl apply -k deploy/kubernetes/overlays/local
	@echo "Waiting for migrations to complete..."
	kubectl wait --for=condition=complete job/$(MIGRATIONS_JOB) -n $(K8S_NAMESPACE) --timeout=60s || true
	@echo "Restarting deployments..."
	kubectl rollout restart deployment $(API_DEPLOYMENT) $(UI_DEPLOYMENT) -n $(K8S_NAMESPACE)

# Database migration commands
migrate-create:
	@echo "Creating new migration..."
	@read -p "Enter migration name: " name; \
	$(GORUN) ./cmd/migration/main.go db create $$name

migrate-up:
	@echo "Applying database migrations..."
	$(GORUN) ./cmd/migration/main.go db up

migrate-down:
	@echo "Reverting database migrations..."
	$(GORUN) ./cmd/migration/main.go db down

migrate-status:
	@echo "Checking migration status..."
	$(GORUN) ./cmd/migration/main.go db status
