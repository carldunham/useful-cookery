#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}Cleaning up local Kubernetes environment...${NC}"

# Check if Pulumi is installed
if ! command -v pulumi &> /dev/null; then
    echo -e "${RED}Pulumi is not installed. Please install it first:${NC}"
    echo "https://www.pulumi.com/docs/get-started/install/"
    exit 1
fi

# Check if k3d is installed
if ! command -v k3d &> /dev/null; then
    echo -e "${RED}k3d is not installed. Please install it first:${NC}"
    echo "https://k3d.io/#installation"
    exit 1
fi

# Destroy Pulumi stack
echo -e "${YELLOW}Destroying Pulumi stack...${NC}"
pulumi destroy --yes

# Delete k3d cluster
echo -e "${YELLOW}Deleting k3d cluster...${NC}"
if k3d cluster list | grep -q "useful-cookery"; then
    k3d cluster delete useful-cookery
else
    echo -e "${YELLOW}Cluster doesn't exist, skipping deletion...${NC}"
fi

# Clean up Docker resources
echo -e "${YELLOW}Cleaning up Docker resources...${NC}"
docker system prune -f
docker image prune -a -f

# Remove hosts entry
echo -e "${YELLOW}Removing hosts entry for useful-cookery.local...${NC}"
if grep -q "useful-cookery.local" /etc/hosts; then
    sudo sed -i '/useful-cookery.local/d' /etc/hosts
else
    echo -e "${YELLOW}Hosts entry doesn't exist, skipping removal...${NC}"
fi

echo -e "${GREEN}Local Kubernetes environment cleanup complete!${NC}"
