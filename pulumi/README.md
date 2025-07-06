# Useful Cookery Infrastructure

This directory contains the Pulumi infrastructure code for the Useful Cookery application.

## Overview

The infrastructure is defined using Pulumi with Go. It provisions the following resources:

### AWS Resources

- VPC with public and private subnets
- EKS cluster for Kubernetes workloads
- RDS PostgreSQL database for application data
- CloudFront distribution for content delivery

### Kubernetes Resources

- Namespace for application resources
- Deployments for API and UI components
- Services for API and UI components

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

## Prerequisites

- [Pulumi CLI](https://www.pulumi.com/docs/get-started/install/)
- [Go 1.24 or later](https://golang.org/doc/install)
- [AWS CLI](https://aws.amazon.com/cli/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

## Configuration

Pulumi stacks are used to manage different environments. We have pre-configured stacks for development, staging, and production environments.

### Development Environment

```bash
cd pulumi
pulumi stack init dev
pulumi stack select dev
```

### Staging Environment

```bash
cd pulumi
pulumi stack init staging
pulumi stack select staging
```

### Production Environment

```bash
cd pulumi
pulumi stack init production
pulumi stack select production
```

### Environment-Specific Configuration

Each environment has its own configuration in the corresponding `Pulumi.<stack>.yaml` file. You can also set or update configuration values using the CLI:

```bash
# Set AWS region
pulumi config set aws:region us-west-2

# Set database password (will be encrypted)
pulumi config set dbPassword --secret

# Set ACM certificate ARN
pulumi config set acmCertificateARN arn:aws:acm:us-east-1:123456789012:certificate/abcdef-1234-5678-abcd-123456789012

# Production-specific settings
pulumi config set highAvailability true
pulumi config set rdsInstanceClass db.t3.medium
pulumi config set eksNodeCount 3
pulumi config set eksNodeType t3.medium
```

## Deployment

Deploy the infrastructure:

```bash
pulumi up
```

This will provision all the resources defined in the Pulumi program.

## Accessing the Application

After deployment, you can access the application using the CloudFront domain name:

```bash
pulumi stack output cdnDomainName
```

## Connecting to the Kubernetes Cluster

Configure kubectl to connect to the EKS cluster:

```bash
aws eks update-kubeconfig --name $(pulumi stack output eksClusterName) --region us-west-2
```

## Local Development

For local development, we provide a Pulumi-based alternative to the shell scripts in `deploy/local`. This approach uses the same Pulumi programming model as the cloud deployment, providing a consistent workflow.

### Setting Up Local Environment

To set up the local development environment:

```bash
cd pulumi/local
./setup.sh
```

This script will:

1. Create a k3d Kubernetes cluster
2. Install NGINX Ingress Controller
3. Add a hosts entry for useful-cookery.local
4. Deploy the application using Pulumi

### Accessing the Local Application

After deployment, you can access the application at:

- UI: <http://useful-cookery.local>
- API: <http://useful-cookery.local/api>

### Cleaning Up Local Environment

To clean up the local environment:

```bash
cd pulumi/local
./cleanup.sh
```

This script will:

1. Destroy the Pulumi stack
2. Delete the k3d cluster
3. Clean up Docker resources
4. Remove the hosts entry

For more details, see the [local development README](./local/README.md).

## Cleanup

To destroy all cloud resources:

```bash
pulumi destroy
```
