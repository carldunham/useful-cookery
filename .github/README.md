# GitHub Actions Workflows

This directory contains GitHub Actions workflows for CI/CD processes in the Useful Cookery project.

## Workflows

### Backend CI (`backend-ci.yml`)

Runs checks for the Go backend code:

- **Trigger**: Push events affecting Go files, all pull requests
- **Jobs**:
  - **Lint**: Runs Go linting using golangci-lint
  - **Test**: Runs Go tests with race detection and coverage reporting
  - **Build**: Builds Go binaries to verify compilation

### Frontend CI (`frontend-ci.yml`)

Runs checks for the React frontend code:

- **Trigger**: Push events affecting UI files, all pull requests
- **Jobs**:
  - **Lint**: Runs JavaScript/TypeScript linting
  - **Test**: Runs frontend tests
  - **Build**: Builds the frontend and archives artifacts

### Additional Checks (`additional-checks.yml`)

Runs additional project-wide checks:

- **Trigger**: Push events affecting Markdown files and config files, all pull requests
- **Jobs**:
  - **Markdown Lint**: Validates Markdown files
  - **Format Check**: Ensures code formatting standards are met using Prettier
  - **Check Generated Code**: Verifies that generated code is up to date with the GraphQL schema

### PR Checks (`pr-checks.yml`)

Validates pull requests before they can be merged:

- **Trigger**: All pull requests
- **Jobs**:
  - **PR Validation**: Checks for Linear issue references (USE-\*) and merge conflicts
  - **Required Checks**: Ensures all individual CI jobs have passed

## Branch Protection Rules

To fully implement the CI/CD process, set up the following branch protection rules for the `main` branch:

1. Require status checks to pass before merging

   - Required status checks:
     - Backend CI
     - Frontend CI
     - Additional Checks
     - PR Checks

2. Require pull request reviews before merging

   - Require at least 1 approval
   - Dismiss stale pull request approvals when new commits are pushed

3. Require linear history

   - Prevent merge commits
   - Prefer squash merging

4. Do not allow bypassing the above settings

## Makefile Integration

The CI workflows leverage targets defined in the project's Makefile:

### Backend Targets

- `lint-go`: Runs Go linting
- `test-go`: Runs Go tests with coverage
- `build-go`: Builds Go binaries

### Frontend Targets

- `lint-ui`: Runs UI linting
- `test-ui`: Runs UI tests
- `build-ui`: Builds UI

### Additional Targets

- `lint-md`: Runs Markdown linting
- `format`: Formats all code
- `format-ui`: Formats UI code
- `format-md`: Formats Markdown
- `format-check`: Checks all code formatting
- `format-check-ui`: Checks UI code formatting
- `format-check-md`: Checks Markdown formatting
- `generate`: Generates code from GraphQL schema

### Combined Targets

- `lint`: Runs all linters
- `test`: Runs all tests
- `build`: Builds everything

These targets are also available for local development to ensure consistency between CI and local environments. For a complete list of available make targets, run:

```bash
make help
```
