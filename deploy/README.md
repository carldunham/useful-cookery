# AWS RDS PostgreSQL + EKS Deployment Solution

This directory contains a comprehensive deployment solution for the Useful Cookery application using AWS RDS PostgreSQL and EKS, with a local Kubernetes environment that mirrors production.

## Directory Structure

- `kubernetes/` - Kubernetes manifests for deploying the application
- `local/` - Scripts and configurations for local development
- `nginx/` - NGINX configuration for the UI service
- `terraform/` - Terraform configuration for AWS infrastructure

## Kubernetes Configuration

### Base Kubernetes Resources

- API and UI deployments with proper resource limits and health checks
- Service definitions for internal communication
- ConfigMaps for application configuration

### Local Development Environment

- k3d-based local Kubernetes cluster setup script
- Local PostgreSQL database deployment
- Ingress configuration for local development
- Resource patches optimized for local development

### Production Environment

- Production-ready deployment configurations with proper scaling
- High availability setup with pod anti-affinity
- AWS ALB Ingress with TLS support
- CloudFront CDN integration for the UI

## AWS Infrastructure (Terraform)

### EKS Cluster

- Properly configured node groups
- Security groups and IAM roles
- AWS Load Balancer Controller integration

### RDS PostgreSQL

- Multi-AZ deployment for high availability
- Automated backups and scaling
- Secure network configuration in private subnets

### CloudFront CDN

- Distribution configuration for UI assets
- TLS certificate integration
- Proper caching policies

## Container Images

### API Service

- Multi-stage build for minimal image size
- Proper configuration handling
- Health check endpoints

### UI Service

- Optimized React build process
- NGINX configuration with proper caching
- API proxy setup

## How to Use

### Local Development

1. Run the local setup script:

   ```sh
   ./deploy/local/setup-k3d.sh
   ```

2. Build and deploy the application:

   ```sh
   docker build -t useful-cookery-api:local -f Dockerfile.api .
   docker build -t useful-cookery-ui:local -f Dockerfile.ui .
   k3d image import useful-cookery-api:local useful-cookery-ui:local -c useful-cookery
   kubectl apply -k deploy/kubernetes/overlays/local
   ```

3. Access the application at <http://useful-cookery.local>

### Production Deployment

1. Initialize Terraform:

   ```sh
   cd deploy/terraform
   terraform init
   ```

2. Apply the Terraform configuration:

   ```sh
   terraform apply -var="environment=production" -var="rds_password=<secure-password>" -var="jwt_secret=<secure-secret>" -var="openai_api_key=<api-key>"
   ```

3. Build and push container images:

   ```sh
   aws ecr get-login-password | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com
   docker build -t ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/useful-cookery-api:latest -f Dockerfile.api .
   docker build -t ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/useful-cookery-ui:latest -f Dockerfile.ui .
   docker push ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/useful-cookery-api:latest
   docker push ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/useful-cookery-ui:latest
   ```

4. Deploy to Kubernetes:

   ```sh
   aws eks update-kubeconfig --region ${AWS_REGION} --name useful-cookery-production
   kubectl apply -k deploy/kubernetes/overlays/production
   ```

## Benefits

This deployment architecture provides:

- A production-ready environment on AWS with RDS PostgreSQL and EKS
- A local development environment that closely mirrors production
- CloudFront CDN integration for optimal performance
- Proper separation of configuration for different environments
- Infrastructure as code with Terraform
