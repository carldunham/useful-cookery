#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to display disk usage
check_disk_usage() {
    echo -e "${BLUE}Current disk usage:${NC}"
    df -h | grep -E 'overlay|Filesystem'
}

# Get the directory where the script is located (not where it's running from)
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo -e "${GREEN}Setting up local Kubernetes environment with k3d...${NC}"

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

# Check if Docker is running
if ! docker info &> /dev/null; then
    echo -e "${RED}Docker is not running. Please start Docker first.${NC}"
    exit 1
fi

# Create volume directory, if it doesn't exist
VOLUME_DIR="$SCRIPT_DIR/k3d-volumes/useful-cookery"
OVERLAY_VOLUME="$VOLUME_DIR/overlay"
STORAGE_VOLUME="$VOLUME_DIR/storage"
if [ ! -d "$VOLUME_DIR" ]; then
    echo -e "${YELLOW}Creating volume directory...${NC}"
    mkdir -p "$VOLUME_DIR"
else
    echo -e "${YELLOW}Volume directory already exists, skipping creation...${NC}"
fi
if [ ! -d "$OVERLAY_VOLUME" ]; then
    echo -e "${YELLOW}Creating overlay volume directory...${NC}"
    mkdir -p "$OVERLAY_VOLUME"
else
    echo -e "${YELLOW}Overlay volume directory already exists, skipping creation...${NC}"
fi
if [ ! -d "$STORAGE_VOLUME" ]; then
    echo -e "${YELLOW}Creating storage volume directory...${NC}"
    mkdir -p "$STORAGE_VOLUME"
else
    echo -e "${YELLOW}Storage volume directory already exists, skipping creation...${NC}"
fi

# Clean up Docker resources to prevent overlay filesystem from filling up
echo -e "${YELLOW}Cleaning up Docker resources to prevent overlay filesystem from filling up...${NC}"
echo -e "${YELLOW}This helps prevent the 'overlay 100%' issue...${NC}"
docker system prune -f
docker image prune -f
docker builder prune -f

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
        --volume "$STORAGE_VOLUME:/var/lib/rancher/k3s/storage@all" \
        --volume "$OVERLAY_VOLUME:/var/lib/containerd@all" \
        --k3s-arg "--kubelet-arg=eviction-hard=imagefs.available<1%,nodefs.available<1%@agent:0" \
        --k3s-arg "--kubelet-arg=eviction-minimum-reclaim=imagefs.available=1%,nodefs.available=1%@agent:0" \
        --k3s-arg "--kubelet-arg=image-gc-high-threshold=85@agent:0" \
        --k3s-arg "--kubelet-arg=image-gc-low-threshold=80@agent:0" \
        --wait
else
    echo -e "${YELLOW}Cluster already exists, skipping creation...${NC}"
fi

# Set kubectl context
kubectl config use-context k3d-useful-cookery

# Create namespace if it doesn't exist
if ! kubectl get namespace useful-cookery-local &> /dev/null; then
    echo -e "${YELLOW}Creating namespace...${NC}"
    kubectl create namespace useful-cookery-local
else
    echo -e "${YELLOW}Namespace already exists, skipping creation...${NC}"
fi

# Install NGINX Ingress Controller
echo -e "${YELLOW}Installing NGINX Ingress Controller...${NC}"
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.7.0/deploy/static/provider/cloud/deploy.yaml

# Wait for NGINX Ingress Controller to be ready
echo -e "${YELLOW}Waiting for NGINX Ingress Controller to be ready...${NC}"
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s

# Create PVC for the PostgreSQL database volume
echo -e "${YELLOW}Creating PostgreSQL database volume PVC...${NC}"
kubectl apply -f - <<EOF
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-pvc
  namespace: useful-cookery-local
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: local-path
  resources:
    requests:
      storage: 1Gi
EOF

# Note: We don't wait for the PVC to be bound here because the local-path
# storage class uses WaitForFirstConsumer binding mode, which means the PVC
# will only be bound after a pod that uses it is created

# Create local PostgreSQL database
echo -e "${YELLOW}Creating local PostgreSQL database...${NC}"
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  namespace: useful-cookery-local
  labels:
    app: postgres
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
        - name: postgres
          image: postgres:14
          ports:
            - containerPort: 5432
          env:
            - name: POSTGRES_USER
              value: postgres
            - name: POSTGRES_PASSWORD
              value: postgres
            - name: POSTGRES_DB
              value: useful-cookery
          volumeMounts:
            - name: postgres-data
              mountPath: /var/lib/postgresql/data
      volumes:
        - name: postgres-data
          persistentVolumeClaim:
            claimName: postgres-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: postgres
  namespace: useful-cookery-local
spec:
  selector:
    app: postgres
  ports:
    - port: 5432
      targetPort: 5432
EOF

# Wait for PostgreSQL to be ready
echo -e "${YELLOW}Waiting for PostgreSQL to be ready...${NC}"
kubectl wait --namespace useful-cookery-local \
  --for=condition=ready pod \
  --selector=app=postgres \
  --timeout=90s

# Add hosts entry
echo -e "${YELLOW}Adding hosts entry for useful-cookery.local...${NC}"
if ! grep -q "useful-cookery.local" /etc/hosts; then
    echo "127.0.0.1 useful-cookery.local" | sudo tee -a /etc/hosts
else
    echo -e "${YELLOW}Hosts entry already exists, skipping...${NC}"
fi

# Check disk usage after setup
check_disk_usage

echo -e "${GREEN}Local Kubernetes environment setup complete!${NC}"
echo -e "${GREEN}You can now build and deploy the application:${NC}"
echo -e "${YELLOW}make deploy-local${NC}"
echo -e "${GREEN}Then access the application at http://useful-cookery.local${NC}"
echo -e "${GREEN}To run database migrations:${NC}"
echo -e "${YELLOW}make migrate-up${NC}"

echo -e "${BLUE}Overlay Filesystem Management:${NC}"
echo -e "${YELLOW}The Docker overlay filesystem can fill up when:${NC}"
echo -e "${YELLOW}- Many container images are built or pulled${NC}"
echo -e "${YELLOW}- Containers with large volumes are created${NC}"
echo -e "${YELLOW}- Build cache grows too large${NC}"
echo -e "${YELLOW}- Old/unused resources aren't cleaned up${NC}"
echo
echo -e "${BLUE}Maintenance Commands:${NC}"
echo -e "${YELLOW}To clean up the overlay filesystem:${NC}"
echo -e "${YELLOW}./$(basename $SCRIPT_DIR)/cleanup-k3d.sh${NC}"
echo -e "${YELLOW}To remove the cluster when not needed:${NC}"
echo -e "${YELLOW}k3d cluster delete useful-cookery${NC}"
echo -e "${YELLOW}To check overlay filesystem usage:${NC}"
echo -e "${YELLOW}df -h | grep overlay${NC}"
