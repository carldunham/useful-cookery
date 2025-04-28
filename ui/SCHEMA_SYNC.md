# Keeping UI Queries in Sync with GraphQL Schema

This document outlines a process for ensuring that UI GraphQL queries remain in sync with the backend schema as it evolves.

## Current Issue

The UI queries were not in sync with the current schema, which was recently changed to use Relay Connections for pagination. This required manual updates to the queries and the components that use them.

## Recommended Process

### 1. Use GraphQL Code Generator

[GraphQL Code Generator](https://www.graphql-code-generator.com/) can automatically generate TypeScript types and React hooks based on your GraphQL schema and operations.

#### Setup

You can set up the GraphQL Code Generator using the provided setup script:

```bash
# From the project root
./ui/scripts/setup-codegen.sh
```

This script will:

1. Install the required dependencies
1. Add the generate script to package.json if it doesn't exist
1. Create the generated directory
1. Run the code generator

Alternatively, you can set it up manually:

1. Install the required packages:

```bash
cd ui
npm install --save-dev @graphql-codegen/cli @graphql-codegen/typescript @graphql-codegen/typescript-operations @graphql-codegen/typescript-react-apollo
```

1. Create a `codegen.yml` configuration file in the `ui` directory:

```yaml
overwrite: true
schema: "../schema/graphql/schema.graphqls"
documents: "src/**/*.{js,jsx,ts,tsx}"
generates:
  src/generated/graphql.tsx:
    plugins:
      - "typescript"
      - "typescript-operations"
      - "typescript-react-apollo"
    config:
      withHooks: true
      withComponent: false
      withHOC: false
      skipTypename: false
      avoidOptionals: true
      preResolveTypes: true
      scalars:
        DateTime: string
```

1. Add a script to `package.json`:

```json
"scripts": {
  "generate": "graphql-codegen --config codegen.yml"
}
```

1. Run the generator:

```bash
npm run generate
```

### 2. Use Generated Hooks

Instead of manually writing GraphQL queries and using `useQuery` directly, use the generated hooks:

```javascript
// Before
const SEARCH_RECIPES = gql`...`;
const { data } = useQuery(SEARCH_RECIPES, { variables });

// After
import { useSearchRecipesQuery } from '../generated/graphql';
const { data } = useSearchRecipesQuery({ variables });
```

### 3. Implement CI Checks

Add a CI check that verifies the generated code is up to date:

```bash
# In your CI pipeline
cd ui
npm run generate
git diff --exit-code src/generated/
```

Alternatively, you can use the provided check script:

```bash
# In your CI pipeline
./ui/scripts/check-graphql-types.sh
```

This script ensures it's running from the correct directory and provides helpful output.

This will fail if the generated code is not committed, indicating that the queries might be out of sync with the schema.

### 4. Schema Change Process

When making schema changes:

1. Update the GraphQL schema files
1. Run the code generator to update the TypeScript types and hooks
1. The TypeScript compiler will show errors where the UI code is not compatible with the new schema
1. Fix the UI code to work with the new schema
1. Commit both the schema changes and the UI changes together

### 5. Consider Schema Versioning

For larger changes, consider versioning your schema to allow for a gradual transition:

1. Add new fields/types alongside existing ones
1. Update the UI to use the new fields/types
1. Once all clients are updated, remove the old fields/types

## Additional Recommendations

### Use Unique Operation Names

Ensure that each GraphQL operation (query, mutation, subscription) has a unique name across your entire codebase. This is required by the GraphQL Code Generator to generate unique hook names.

```javascript
// Bad: Multiple queries with the same name in different files
// File 1
const SEARCH_RECIPES = gql`
  query SearchRecipes($query: String!) {
    // ...
  }
`;

// File 2
const SEARCH_RECIPES = gql`
  query SearchRecipes($query: String!) {
    // ...
  }
`;

// Good: Each query has a unique name
// File 1
const SEARCH_RECIPES = gql`
  query HomeSearchRecipes($query: String!) {
    // ...
  }
`;

// File 2
const SEARCH_RECIPES = gql`
  query DetailSearchRecipes($query: String!) {
    // ...
  }
`;
```

### Use Fragment Colocation

Define GraphQL fragments alongside the components that use them. This makes it clear what data each component needs and helps prevent over-fetching.

```javascript
export const RECIPE_CARD_FRAGMENT = gql`
  fragment RecipeCardFragment on Recipe {
    id
    title
    description
    # ...other fields
  }
`;

// Then in queries
const SEARCH_RECIPES = gql`
  query SearchRecipes($query: String!) {
    searchRecipes(query: $query) {
      edges {
        node {
          ...RecipeCardFragment
        }
      }
      # ...other fields
    }
  }
  ${RECIPE_CARD_FRAGMENT}
`;
```

### Document Schema Changes

Maintain a changelog of schema changes to help track when and why changes were made.

### Regular Audits

Periodically audit your GraphQL queries to ensure they're not over-fetching data and are using the schema correctly.

## Conclusion

By implementing these processes, we can ensure that UI queries stay in sync with the GraphQL schema as it evolves, reducing the likelihood of runtime errors and making schema changes easier to manage.
