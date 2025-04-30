# Local Kubernetes Configuration

This directory contains Kubernetes configuration files for local development.

## Applying Configuration

To apply the local Kubernetes configuration:

```bash
kubectl apply -k .
```

### Troubleshooting: Job Immutability

If you encounter an error like `The Job "useful-cookery-migrations" is invalid: spec.template: Invalid value: ... field is immutable`, this is because Kubernetes Jobs are immutable resources. Once created, they cannot be modified.

To resolve this issue, delete the existing Job before reapplying:

```bash
kubectl delete job useful-cookery-migrations -n useful-cookery-local
kubectl apply -k .
```

This is a common issue when reapplying configurations that include Jobs.

## Secrets Management

To set up your local environment, you need to create a `secrets.yaml` file based on the provided example:

1. Copy the example file:

  ```bash
  cp secrets-example.yaml secrets.yaml
  ```

1. Edit `secrets.yaml` to add your actual secret values.

⚠️ **IMPORTANT**: Never commit the `secrets.yaml` file to Git! ⚠️

The `.gitignore` file has been configured to exclude `secrets.yaml` files, but always verify your commits to ensure sensitive information is not accidentally included.

## Why Secrets Should Never Be Committed

- Security credentials in Git history are permanently accessible
- Even if deleted later, they remain in the repository history
- Exposed secrets require immediate rotation
- Secret leakage can lead to unauthorized access and data breaches

Always use appropriate secret management practices for your Kubernetes deployments.
