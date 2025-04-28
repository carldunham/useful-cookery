#!/bin/bash
set -e

# This script checks if the GraphQL types are up to date
# It's meant to be run in CI to ensure that the generated types match the schema

# Ensure we're in the ui directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
UI_DIR="$( cd "$SCRIPT_DIR/.." && pwd )"
cd "$UI_DIR"

echo "Checking if GraphQL types are up to date from $(pwd)..."
echo "Node version: $(node -v)"
echo "NPM version: $(npm -v)"

# Check if package.json exists
if [ ! -f "package.json" ]; then
  echo "❌ Error: package.json not found in $(pwd)"
  exit 1
fi

# Check if codegen.yml exists
if [ ! -f "codegen.yml" ]; then
  echo "❌ Error: codegen.yml not found in $(pwd)"
  exit 1
fi

# Check if the generate script is defined in package.json
if ! grep -q '"generate"' package.json; then
  echo "❌ Error: 'generate' script not found in package.json"
  exit 1
fi

# Save the current state of the generated directory
if [ -d "src/generated" ]; then
  mkdir -p /tmp/graphql-check
  cp -r src/generated /tmp/graphql-check/
fi

# Run the code generator with more verbose output
echo "Running GraphQL code generator..."
echo "Command: npm run generate"
npm run generate || {
  echo "❌ Error: Failed to run 'npm run generate'"

  # Check if the error is due to duplicate operation names
  if grep -q "Not all operations have an unique name" src/generated/graphql.tsx 2>/dev/null; then
    echo "❌ Error: Duplicate operation names detected in your GraphQL queries."
    echo "Make sure each query has a unique name across your codebase."
    echo "For example, if you have multiple files with a query named 'SearchRecipes',"
    echo "rename them to be unique like 'HomeSearchRecipes', 'DetailSearchRecipes', etc."
    echo ""
    echo "The following operation names are duplicated:"
    grep -A 10 "Not all operations have an unique name" src/generated/graphql.tsx 2>/dev/null || true
    exit 1
  fi

  echo "Please check that @graphql-codegen/cli and other required packages are installed"
  echo "You may need to run: npm install --save-dev @graphql-codegen/cli @graphql-codegen/typescript @graphql-codegen/typescript-operations @graphql-codegen/typescript-react-apollo"
  exit 1
}

# Check if there are differences
if [ -d "/tmp/graphql-check/generated" ]; then
  if diff -r src/generated /tmp/graphql-check/generated > /dev/null; then
    echo "✅ GraphQL types are up to date!"
    exit 0
  else
    echo "❌ GraphQL types are not up to date!"
    echo "Please run 'npm run generate' and commit the changes."
    exit 1
  fi
else
  echo "⚠️ No previous generated files found. Assuming this is the first run."
  echo "✅ GraphQL types have been generated!"
  exit 0
fi
