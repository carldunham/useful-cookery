#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

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
          emptyDir: {}
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

echo -e "${GREEN}Local Kubernetes environment setup complete!${NC}"
echo -e "${GREEN}You can now build and deploy the application:${NC}"
echo -e "${YELLOW}docker build -t useful-cookery-api:local -f Dockerfile.api .${NC}"
echo -e "${YELLOW}docker build -t useful-cookery-ui:local -f Dockerfile.ui .${NC}"
echo -e "${YELLOW}k3d image import useful-cookery-api:local useful-cookery-ui:local -c useful-cookery${NC}"
echo -e "${YELLOW}kubectl apply -k deploy/kubernetes/overlays/local${NC}"
echo -e "${GREEN}Then access the application at http://useful-cookery.local${NC}"
