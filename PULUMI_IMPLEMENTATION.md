# Pulumi Implementation for Useful Cookery

This document outlines the implementation of Pulumi for infrastructure management in the Useful Cookery project.

## Overview

We've implemented Pulumi as an infrastructure-as-code solution for both cloud deployment and local development. This approach provides several benefits:

1. **Consistent Workflow**: The same programming model is used for both local and cloud environments
2. **Declarative Infrastructure**: Infrastructure is defined as code, making it easier to understand and maintain
3. **Drift Detection**: Pulumi can detect when the actual state of the infrastructure differs from the desired state
4. **Incremental Updates**: Only changed resources are updated, reducing deployment time
5. **Type Safety**: Using Go provides compile-time type checking for infrastructure code

## Project Structure

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

## Cloud Deployment

The cloud deployment uses Pulumi to provision AWS resources and Kubernetes resources on EKS. The main components are:

1. **VPC**: A Virtual Private Cloud with public and private subnets
2. **EKS**: A Kubernetes cluster for container orchestration
3. **RDS**: A PostgreSQL database for application data
4. **CloudFront**: A CDN distribution for content delivery
5. **Kubernetes Resources**: Namespace, Deployments, and Services for the application components

## Local Development

For local development, we've created a Pulumi-based alternative to the existing shell scripts in `deploy/local`. This approach uses the same Pulumi programming model as the cloud deployment, providing a consistent workflow.

The local development setup includes:

1. **k3d Kubernetes Cluster**: A lightweight Kubernetes cluster running in Docker
2. **PostgreSQL Database**: A containerized database for local development
3. **NGINX Ingress Controller**: For routing HTTP traffic
4. **Application Deployments**: For API and UI components

## Makefile Integration

We've updated the Makefile to include targets for the Pulumi-based local development approach:

- `make setup-local-pulumi`: Set up the local k3d environment with Pulumi
- `make deploy-local-pulumi`: Deploy to the local k3d environment with Pulumi
- `make cleanup-local-pulumi`: Clean up the local k3d environment with Pulumi

These targets provide a convenient way to use the Pulumi-based approach alongside the existing shell script-based approach.

## Comparison with Shell Scripts

The Pulumi-based approach offers several advantages over the shell script approach:

1. **Declarative vs. Imperative**: Pulumi defines the desired state of the infrastructure, while shell scripts define the steps to create the infrastructure
2. **Drift Detection**: Pulumi can detect when the actual state of the infrastructure differs from the desired state
3. **Incremental Updates**: Pulumi can make incremental updates to the infrastructure, only changing what needs to be changed
4. **Reusable Components**: The Pulumi code can be reused across different environments with minimal changes
5. **Consistent Workflow**: Using Pulumi for both local and cloud deployments provides a consistent workflow

## Getting Started

### Cloud Deployment

1. Install prerequisites (Pulumi CLI, Go, AWS CLI, kubectl)
2. Configure AWS credentials
3. Create a new Pulumi stack: `pulumi stack init dev`
4. Configure required variables
5. Deploy with `pulumi up`

### Local Development

1. Install prerequisites (Pulumi CLI, Go, k3d, Docker)
2. Run `make setup-local-pulumi` to set up the local environment
3. Access the application at <http://useful-cookery.local>
4. Clean up with `make cleanup-local-pulumi` when done

## Future Improvements

1. **CI/CD Integration**: Integrate Pulumi with CI/CD pipelines for automated deployments
2. **Secret Management**: Use Pulumi's secret management for sensitive information
3. **Multi-Environment Support**: Create Pulumi stacks for different environments (dev, staging, production)
4. **Infrastructure Testing**: Add tests for the infrastructure code
5. **Custom Components**: Create custom Pulumi components for common infrastructure patterns
