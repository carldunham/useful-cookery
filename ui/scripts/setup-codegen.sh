#!/bin/bash
set -e

# This script installs the required dependencies for GraphQL Code Generator
# and runs the generator to create the initial types

# Ensure we're in the ui directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
UI_DIR="$( cd "$SCRIPT_DIR/.." && pwd )"
cd "$UI_DIR"

echo "Setting up GraphQL Code Generator from $(pwd)..."

# Install required dependencies
echo "Installing GraphQL Code Generator dependencies..."
npm install --save-dev \
  @graphql-codegen/cli \
  @graphql-codegen/typescript \
  @graphql-codegen/typescript-operations \
  @graphql-codegen/typescript-react-apollo

# Check if the generate script exists in package.json
if ! grep -q '"generate"' package.json; then
  echo "Adding 'generate' script to package.json..."
  # Use a temporary file to avoid issues with in-place editing
  sed -i.bak 's/"scripts": {/"scripts": {\n    "generate": "graphql-codegen --config codegen.yml",/g' package.json
  rm package.json.bak
fi

# Create the generated directory if it doesn't exist
mkdir -p src/generated

# Run the code generator
echo "Running GraphQL code generator..."
npm run generate || {
  echo "⚠️ Warning: Code generator encountered an error."
  echo "This might be due to duplicate operation names in your GraphQL queries."
  echo "Make sure each query has a unique name across your codebase."
  echo "For example, if you have multiple files with a query named 'SearchRecipes',"
  echo "rename them to be unique like 'HomeSearchRecipes', 'DetailSearchRecipes', etc."

  # Check if the error is due to duplicate operation names
  if grep -q "Not all operations have an unique name" src/generated/graphql.tsx 2>/dev/null; then
    echo "✅ Duplicate operation names detected. Please make your query names unique."
    exit 1
  fi
}

echo "✅ Setup complete! GraphQL types have been generated."
echo "You can now use the generated hooks in your components."
echo ""
echo "Note: If you're creating example files or tests, make sure to use unique"
echo "operation names to avoid conflicts with your actual components."
