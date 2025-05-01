# Contributing to Useful Cookery

Thank you for your interest in contributing to Useful Cookery! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

Please be respectful and considerate of others when contributing to this project. We aim to foster an inclusive and welcoming community.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/useful-cookery.git`
3. Add the upstream repository: `git remote add upstream https://github.com/carldunham/useful-cookery.git`
4. Create a new branch for your changes: `git checkout -b USE-123-feature-name`
   - Always include the Linear issue ID (USE-\*) in your branch name

## Development Workflow

1. Make your changes in your feature branch
2. Write or update tests as necessary
3. Ensure all tests pass locally
4. Commit your changes with clear, descriptive commit messages
5. Push your changes to your fork
6. Create a pull request to the main repository

## Pull Request Process

1. Create a pull request from your fork to the main repository
2. Include the Linear issue ID (USE-\*) in the PR title or description
3. Provide a clear description of the changes and the problem they solve
4. Ensure all CI checks pass
5. Request a review from a maintainer
6. Address any feedback from reviewers
7. Once approved, a maintainer will merge your PR

## Continuous Integration

All pull requests go through our CI pipeline before they can be merged. The CI pipeline includes:

### Backend Checks

- Go linting with golangci-lint
- Go tests with race detection and coverage reporting
- Go build verification

### Frontend Checks

- JavaScript/TypeScript linting
- React component tests
- Build verification

### Additional Checks

- Markdown linting
- Code formatting verification
- Project-wide linting

### PR Validation

- Linear issue reference check (USE-\*)
- Merge conflict detection
- Required status checks verification

## Running CI Checks Locally

Before submitting a PR, run the same checks locally that will be run in CI:

```bash
# Backend checks
make lint-go     # Run Go linting
make test-go     # Run Go tests
make build-go    # Build Go binaries

# Frontend checks
make lint-ui     # Run UI linting
make test-ui     # Run UI tests
make build-ui    # Build UI

# Additional checks
make lint-md     # Run Markdown linting
make format-check     # Check code formatting
make format-check-ui  # Check UI code formatting
make format-check-md  # Check Markdown formatting

# Formatting code
make format      # Format all code
make format-ui   # Format UI code
make format-md   # Format Markdown

# Run all checks
make lint        # Run all linters
make test        # Run all tests
make build       # Build everything
```

For a complete list of available make targets, run:

```bash
make help
```

## Coding Standards

### Go Code

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Document all exported functions, types, and constants
- Write tests for all new functionality

### JavaScript/TypeScript Code

- Follow the ESLint configuration in the project
- Use Prettier for code formatting
- Write tests for all new components and functionality
- Follow React best practices

### Commit Messages

- Use clear, descriptive commit messages
- Include the Linear issue ID (USE-\*) in the commit message
- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")

## Documentation

- Update documentation when changing functionality
- Use clear, concise language
- Include examples where appropriate
- Follow Markdown best practices

## Questions?

If you have any questions or need help with the contribution process, please reach out to the maintainers.

Thank you for contributing to Useful Cookery!
