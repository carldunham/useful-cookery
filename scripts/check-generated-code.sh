#!/bin/bash
# Script to check if generated code is up to date

set -e

# Store the current directory
CURRENT_DIR=$(pwd)

# Navigate to the project root (assuming this script is in the scripts directory)
cd "$(dirname "$0")/.."

echo "Checking if generated code is up to date..."

# Create a temporary directory for git operations
TEMP_DIR=$(mktemp -d)
trap 'rm -rf "$TEMP_DIR"' EXIT

# Initialize a temporary git repository
git init "$TEMP_DIR" > /dev/null 2>&1
git --git-dir="$TEMP_DIR/.git" config --local user.email "ci@example.com"
git --git-dir="$TEMP_DIR/.git" config --local user.name "CI"

# Copy current files to the temporary directory
rsync -a --exclude=".git" --exclude="node_modules" --exclude="ui/node_modules" . "$TEMP_DIR/"

# Navigate to the temporary directory
cd "$TEMP_DIR"

# Add all files and commit
git add .
git commit -m "Initial commit" > /dev/null 2>&1

# Run code generation
echo "Running 'make generate'..."
make generate

# Check if any files were changed
if [[ -n "$(git status --porcelain)" ]]; then
    echo "ERROR: Generated code is not up to date!"
    echo "The following files need to be updated:"
    git status --porcelain
    echo ""
    echo "Diff:"
    git diff
    echo ""
    echo "Please run 'make generate' locally and commit the changes."
    exit 1
else
    echo "Success: Generated code is up to date!"
fi

# Return to the original directory
cd "$CURRENT_DIR"

exit 0
