# Local Development with Pulumi

This directory contains Pulumi code for deploying the Useful Cookery application to a local Kubernetes cluster using k3d.

## Overview

The local deployment uses Pulumi to provision the following resources in a local k3d Kubernetes cluster:

- Namespace for application resources
- PostgreSQL database with persistent storage
- API deployment and service
- UI deployment and service
- Ingress for routing traffic

## Prerequisites

- [Pulumi CLI](https://www.pulumi.com/docs/get-started/install/)
- [Go 1.24 or later](https://golang.org/doc/install)
- [k3d](https://k3d.io/#installation)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [Docker](https://docs.docker.com/get-docker/)

## Setup

### 1. Create a k3d cluster

Before deploying with Pulumi, you need to create a k3d cluster:

```bash
k3d cluster create useful-cookery \
    --api-port 6550 \
    --servers 1 \
    --agents 2 \
    --port 80:80@loadbalancer \
    --port 443:443@loadbalancer \
    --k3s-arg "--disable=traefik@server:0"
```

### 2. Install NGINX Ingress Controller

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.7.0/deploy/static/provider/cloud/deploy.yaml
```

### 3. Add hosts entry

```bash
echo "127.0.0.1 useful-cookery.local" | sudo tee -a /etc/hosts
```

## Deployment

Initialize a new Pulumi stack:

```bash
cd pulumi/local
pulumi stack init dev
```

Deploy the infrastructure:

```bash
pulumi up
```

This will provision all the resources defined in the Pulumi program.

## Accessing the Application

After deployment, you can access the application at:

- UI: <http://useful-cookery.local>
- API: <http://useful-cookery.local/api>

## Cleanup

To destroy all resources:

```bash
pulumi destroy
```

To delete the k3d cluster when not needed:

```bash
k3d cluster delete useful-cookery
```

## Comparison with Shell Scripts

This Pulumi-based approach offers several advantages over the shell script approach:

1. **Declarative Infrastructure**: Pulumi defines the desired state of the infrastructure, making it easier to understand and maintain.

2. **Drift Detection**: Pulumi can detect when the actual state of the infrastructure differs from the desired state.

3. **Incremental Updates**: Pulumi can make incremental updates to the infrastructure, only changing what needs to be changed.

4. **Reusable Components**: The Pulumi code can be reused across different environments with minimal changes.

5. **Consistent Workflow**: Using Pulumi for both local and cloud deployments provides a consistent workflow.

## Integration with Cloud Deployment

The local Pulumi deployment can be used as a stepping stone to the cloud deployment. The same Pulumi code structure can be used for both local and cloud deployments, with environment-specific configurations.
