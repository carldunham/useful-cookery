# Migration from Shell Scripts to Pulumi

This document outlines the migration from the shell script-based deployment approach to the Pulumi-based approach.

## Overview

We've migrated from using shell scripts and Kubernetes manifests in the `deploy` directory to using Pulumi for infrastructure management. This migration provides several benefits:

1. **Consistent Workflow**: The same programming model is used for both local and cloud environments
2. **Declarative Infrastructure**: Infrastructure is defined as code, making it easier to understand and maintain
3. **Drift Detection**: Pulumi can detect when the actual state of the infrastructure differs from the desired state
4. **Incremental Updates**: Only changed resources are updated, reducing deployment time
5. **Type Safety**: Using Go provides compile-time type checking for infrastructure code

## Changes Made

### 1. New Directory Structure

We've created a new `pulumi` directory with the following structure:

```
pulumi/
├── aws/                    # AWS infrastructure modules
│   ├── vpc.go              # VPC configuration
│   ├── eks.go              # EKS cluster configuration
│   ├── rds.go              # RDS database configuration
│   └── cloudfront.go       # CloudFront distribution configuration
├── kubernetes/             # Kubernetes resource modules
│   ├── namespace.go        # Namespace configuration
│   ├── deployment.go       # Deployment configuration
│   └── service.go          # Service configuration
├── local/                  # Local development with k3d
│   ├── main.go             # Local Kubernetes deployment
│   ├── setup.sh            # Setup script for local environment
│   ├── cleanup.sh          # Cleanup script for local environment
│   └── README.md           # Documentation for local development
├── main.go                 # Main Pulumi program for cloud deployment
├── go.mod                  # Go module definition
├── Pulumi.yaml             # Pulumi project configuration
└── README.md               # Project documentation
```

### 2. Makefile Updates

We've updated the Makefile to remove the old deployment targets and add new ones for the Pulumi-based approach:

- Removed:
  - `deploy-api-local`
  - `deploy-ui-local`
  - `deploy-migrations-local`
  - `deploy-local`

- Added:
  - `setup-local`: Set up the local k3d environment with Pulumi
  - `deploy-local`: Deploy to the local k3d environment with Pulumi
  - `cleanup-local`: Clean up the local k3d environment with Pulumi
  - `deploy-cloud`: Deploy to the cloud with Pulumi
  - `cleanup-cloud`: Clean up cloud resources with Pulumi

### 3. Removed `deploy` Directory

The `deploy` directory is no longer needed, as all deployment functionality has been migrated to the Pulumi-based approach.

### 4. Multi-Environment Support

We've implemented multi-environment support using Pulumi stacks:

- `dev`: Development environment
- `staging`: Staging environment
- `production`: Production environment

Each environment has its own configuration in the corresponding `Pulumi.<stack>.yaml` file, allowing for environment-specific settings such as:

- High availability configuration
- Instance types and sizes
- Node counts
- Other environment-specific parameters

### 5. CI/CD Integration

We've added GitHub Actions workflows for CI/CD integration:

- Pull Request: Runs tests and `pulumi preview` to show changes without applying them
- Push to main: Runs tests and `pulumi up` to apply changes
- Manual trigger: Allows for manual deployment

The workflow handles:

- Running unit tests for all Pulumi code
- Building Docker images
- Pushing images to ECR
- Deploying infrastructure with Pulumi
- Environment-specific configuration

### 6. Unit Tests

We've added unit tests for the Pulumi code following Go's standard practices:

- Tests are in separate `_test` packages to ensure they only test the public API
- Helper functions in the main package are exported with capitalized names to make them accessible to tests
- Tests for helper functions that handle environment-specific configuration
- Tests for the main Pulumi program
- Tests for the local development Pulumi program

The tests use mocks to simulate the Pulumi runtime environment and verify that the code behaves as expected. The tests are run as part of the CI/CD pipeline to ensure that changes to the infrastructure code don't break existing functionality.

## 7. Go Workspaces

We've set up Go workspaces to manage the multiple Go modules in the repository:

- The main application module in the root directory
- The main Pulumi module in `pulumi/`
- The local development module in `pulumi/local/`

Go workspaces (introduced in Go 1.18) allow for easier management of multi-module repositories. We've created a `go.work` file at the root of the repository that includes all modules, making it easier to work with the code and run tests across all modules.

The Makefile has been updated to use Go workspaces for running tests and tidying dependencies:

```makefile
# Run Pulumi tests
test-pulumi:
 @echo "Running Pulumi tests..."
 cd $(PULUMI_DIR) && $(GOTEST) -v ./...
 cd $(PULUMI_LOCAL_DIR) && $(GOTEST) -v ./...

# Tidy up dependencies
tidy:
 @echo "Tidying dependencies..."
 $(GOMOD) tidy
 cd $(PULUMI_DIR) && $(GOMOD) tidy
 cd $(PULUMI_LOCAL_DIR) && $(GOMOD) tidy
```

This approach simplifies development and testing, as you can now run tests for all Pulumi modules with a single command.

## Migration Steps

If you're still using the old deployment approach, follow these steps to migrate to the Pulumi-based approach:

### 1. Install Prerequisites

- [Pulumi CLI](https://www.pulumi.com/docs/get-started/install/)
- [Go 1.24 or later](https://golang.org/doc/install)
- [k3d](https://k3d.io/#installation) (for local development)
- [Docker](https://docs.docker.com/get-docker/) (for local development)
- [AWS CLI](https://aws.amazon.com/cli/) (for cloud deployment)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

### 2. Set Up Local Development

```bash
make setup-local
```

This will create a k3d cluster, install NGINX Ingress Controller, and set up the local environment.

### 3. Deploy Locally

```bash
make deploy-local
```

This will build the Docker images, import them to k3d, and deploy the application using Pulumi.

### 4. Clean Up Local Environment

```bash
make cleanup-local
```

This will destroy the Pulumi stack, delete the k3d cluster, and clean up Docker resources.

### 5. Deploy to Cloud

```bash
# Select the appropriate environment
pulumi stack select dev  # or staging, or production

# Deploy
make deploy-cloud
```

This will deploy the application to AWS using Pulumi with the selected environment's configuration.

### 6. Clean Up Cloud Resources

```bash
make cleanup-cloud
```

This will destroy the cloud resources created by Pulumi.

## Benefits of the Migration

1. **Simplified Workflow**: The Pulumi-based approach provides a simpler workflow for both local and cloud deployment.
2. **Improved Maintainability**: The infrastructure code is more maintainable and easier to understand.
3. **Better Error Handling**: Pulumi provides better error handling and reporting than shell scripts.
4. **Incremental Updates**: Pulumi can make incremental updates to the infrastructure, only changing what needs to be changed.
5. **Drift Detection**: Pulumi can detect when the actual state of the infrastructure differs from the desired state.
6. **Type Safety**: Using Go provides compile-time type checking for infrastructure code.
7. **Multi-Environment Support**: Easily manage different environments with environment-specific configurations.
8. **CI/CD Integration**: Automated deployments through GitHub Actions.

## Conclusion

The migration from shell scripts to Pulumi provides a more robust and maintainable approach to infrastructure management. The Pulumi-based approach is easier to understand, maintain, and extend, and provides a consistent workflow for both local and cloud deployment. With the addition of multi-environment support and CI/CD integration, we now have a complete infrastructure-as-code solution that meets our needs for both development and production environments.
