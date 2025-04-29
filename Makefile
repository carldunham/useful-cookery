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

.PHONY: all build clean test lint fix generate tidy help deploy-api-local deploy-ui-local deploy-migrations-local deploy-local migrate-create migrate-up migrate-down migrate-status

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
	@echo "Running Go tests..."
	$(GOTEST) ./...
	@echo "Running UI tests..."
	cd ui && npm run test:nowatch

# Run linting
lint:
	@echo "Running Go linter..."
	$(GOTOOL) $(GOLANGCI_LINT) run ./...
	@echo "Running UI linter..."
	cd ui && npm run lint
	@echo "Running Markdown linter..."
	npm run lint:md
	@echo "Running project linter..."
	npm run lint

# Fix linting issues
fix:
	@echo "Fixing Go linting issues..."
	$(GOTOOL) $(GOLANGCI_LINT) run --fix ./...
	@echo "Fixing UI linting issues..."
	cd ui && npm run lint:fix
	@echo "Fixing Markdown linting issues..."
	npm run lint:md:fix
	@echo "Fixing project linting issues..."
	npm run lint:fix
	@echo "Formatting code..."
	npm run format:all

# Generate code as needed
generate:
	@echo "Generating Go code..."
	$(GOGEN) ./...

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
	@echo "  make test         Run Go and UI tests"
	@echo "  make lint         Run all linters (Go, UI, Markdown, project)"
	@echo "  make fix          Fix linting issues and format code"
	@echo "  make generate     Generate code with gqlgen"
	@echo "  make tidy         Tidy up dependencies"
	@echo "  make run          Run the application"
	@echo "  make deploy-api-local        Build and deploy API locally"
	@echo "  make deploy-ui-local         Build and deploy UI locally"
	@echo "  make deploy-migrations-local Build and deploy migrations locally"
	@echo "  make deploy-local            Build and deploy all components locally"
	@echo "  make migrate-create          Create a new database migration"
	@echo "  make migrate-up              Apply database migrations"
	@echo "  make migrate-down            Revert database migrations"
	@echo "  make migrate-status          Show current migration status"
	@echo "  make help         Show this help message"

# Build and deploy API locally
deploy-api-local:
	@echo "Building API Docker image..."
	docker build -t $(API_IMAGE):$(IMAGE_TAG) -f Dockerfile.api .
	@echo "Importing API image to k3d..."
	k3d image import $(API_IMAGE):$(IMAGE_TAG) -c $(K3D_CLUSTER)
	@echo "Deploying API to local Kubernetes..."
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
	kubectl apply -k deploy/kubernetes/overlays/local
	@echo "Running migrations job..."
	kubectl delete job $(MIGRATIONS_JOB) -n $(K8S_NAMESPACE) --ignore-not-found
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
	kubectl apply -k deploy/kubernetes/overlays/local
	@echo "Running migrations job..."
	kubectl delete job $(MIGRATIONS_JOB) -n $(K8S_NAMESPACE) --ignore-not-found
	kubectl apply -k deploy/kubernetes/overlays/local
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
