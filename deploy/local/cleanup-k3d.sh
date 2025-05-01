#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print header
echo -e "${GREEN}=========================================${NC}"
echo -e "${GREEN}Docker Overlay Filesystem Cleanup Script${NC}"
echo -e "${GREEN}=========================================${NC}"
echo -e "${YELLOW}This script cleans up Docker's overlay filesystem to free up disk space.${NC}"
echo

# Function to display disk usage
check_disk_usage() {
    echo -e "${BLUE}Current disk usage:${NC}"
    df -h | grep -E 'overlay|Filesystem'
}

# Check initial disk usage
echo -e "${BLUE}Initial disk usage:${NC}"
check_disk_usage

# Perform cleanup operations
echo -e "${YELLOW}Cleaning up Docker resources to free overlay filesystem space...${NC}"

# Remove unused containers
echo -e "${YELLOW}Removing unused containers...${NC}"
docker container prune -f

# Remove unused networks
echo -e "${YELLOW}Removing unused networks...${NC}"
docker network prune -f

# Remove unused volumes (be careful with this one)
echo -e "${YELLOW}Removing unused volumes...${NC}"
docker volume prune -f

# Remove unused images (this is a major source of overlay filesystem usage)
echo -e "${YELLOW}Removing unused images (major overlay filesystem consumer)...${NC}"
docker image prune -a -f

# Clean Docker build cache (frees overlay filesystem space)
echo -e "${YELLOW}Cleaning Docker build cache...${NC}"
docker builder prune -a -f

# Run a full system prune
echo -e "${YELLOW}Running full system prune...${NC}"
docker system prune -a -f --volumes

# Check final disk usage
echo -e "${BLUE}Final disk usage after cleanup:${NC}"
check_disk_usage

echo -e "${GREEN}Cleanup complete!${NC}"
echo -e "${BLUE}About the overlay filesystem:${NC}"
echo -e "${YELLOW}Docker's overlay filesystem stores container layers and images.${NC}"
echo -e "${YELLOW}It can fill up when many images are built or containers are run.${NC}"
echo -e "${YELLOW}Regular cleanup using this script will help prevent overlay filesystem issues.${NC}"
echo
echo -e "${BLUE}Additional commands:${NC}"
echo -e "${YELLOW}To delete the k3d cluster when not in use:${NC}"
echo -e "${YELLOW}k3d cluster delete useful-cookery${NC}"
echo -e "${YELLOW}To restart Docker (last resort if overlay is still full):${NC}"
echo -e "${YELLOW}sudo systemctl restart docker${NC}"
