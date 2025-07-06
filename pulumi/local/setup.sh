#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}Setting up local Kubernetes environment with k3d and Pulumi...${NC}"

# Check if k3d is installed
if ! command -v k3d &> /dev/null; then
    echo -e "${RED}k3d is not installed. Please install it first:${NC}"
    echo "https://k3d.io/#installation"
    exit 1
fi

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo -e "${RED}kubectl is not installed. Please install it first:${NC}"
    echo "https://kubernetes.io/docs/tasks/tools/install-kubectl/"
    exit 1
fi

# Check if Pulumi is installed
if ! command -v pulumi &> /dev/null; then
    echo -e "${RED}Pulumi is not installed. Please install it first:${NC}"
    echo "https://www.pulumi.com/docs/get-started/install/"
    exit 1
fi

# Check if Docker is running
if ! docker info &> /dev/null; then
    echo -e "${RED}Docker is not running. Please start Docker first.${NC}"
    exit 1
fi

# Create k3d cluster if it doesn't exist
if ! k3d cluster list | grep -q "useful-cookery"; then
    echo -e "${YELLOW}Creating k3d cluster...${NC}"
    k3d cluster create useful-cookery \
        --api-port 6550 \
        --servers 1 \
        --agents 2 \
        --port 80:80@loadbalancer \
        --port 443:443@loadbalancer \
        --k3s-arg "--disable=traefik@server:0" \
        --wait
else
    echo -e "${YELLOW}Cluster already exists, skipping creation...${NC}"
fi

# Set kubectl context
kubectl config use-context k3d-useful-cookery

# Install NGINX Ingress Controller
echo -e "${YELLOW}Installing NGINX Ingress Controller...${NC}"
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.7.0/deploy/static/provider/cloud/deploy.yaml

# Wait for NGINX Ingress Controller to be ready
echo -e "${YELLOW}Waiting for NGINX Ingress Controller to be ready...${NC}"
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s

# Add hosts entry
echo -e "${YELLOW}Adding hosts entry for useful-cookery.local...${NC}"
if ! grep -q "useful-cookery.local" /etc/hosts; then
    echo "127.0.0.1 useful-cookery.local" | sudo tee -a /etc/hosts
else
    echo -e "${YELLOW}Hosts entry already exists, skipping...${NC}"
fi

# Initialize Pulumi stack if it doesn't exist
if ! pulumi stack ls | grep -q "dev"; then
    echo -e "${YELLOW}Initializing Pulumi stack...${NC}"
    pulumi stack init dev
else
    echo -e "${YELLOW}Pulumi stack already exists, skipping initialization...${NC}"
fi

# Deploy Pulumi stack
echo -e "${YELLOW}Deploying Pulumi stack...${NC}"
pulumi up --yes

echo -e "${GREEN}Local Kubernetes environment setup complete!${NC}"
echo -e "${GREEN}You can access the application at:${NC}"
echo -e "${YELLOW}UI: http://useful-cookery.local${NC}"
echo -e "${YELLOW}API: http://useful-cookery.local/api${NC}"

echo -e "${BLUE}Cleanup Commands:${NC}"
echo -e "${YELLOW}To destroy Pulumi resources:${NC}"
echo -e "${YELLOW}pulumi destroy${NC}"
echo -e "${YELLOW}To delete the k3d cluster when not needed:${NC}"
echo -e "${YELLOW}k3d cluster delete useful-cookery${NC}"
