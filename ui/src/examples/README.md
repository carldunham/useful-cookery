# GraphQL Query Examples

This directory contains example implementations demonstrating best practices for working with GraphQL in the UI.

## Files

### RecipeCardWithFragment.js

Demonstrates the fragment colocation pattern, where GraphQL fragments are defined alongside the components that use them. This makes it clear what data each component needs and helps prevent over-fetching.

### RecipeSearchWithFragments.js

Shows how to use the fragments defined in components within GraphQL queries. This approach:

- Keeps queries DRY by reusing fragments
- Makes it clear what data is being requested for each component
- Ensures that if a component's data requirements change, you only need to update the fragment in one place

### RecipeSearchWithGeneratedHooks.js

Demonstrates how to use the generated hooks from GraphQL Code Generator. This approach:

- Provides type safety for queries and their results
- Eliminates the need to manually write GraphQL queries as strings
- Ensures queries stay in sync with the schema

## Best Practices

1. **Use Fragment Colocation**: Define fragments alongside the components that use them
2. **Use Generated Hooks**: Let the code generator create typed hooks for your queries
3. **Implement Pagination Properly**: Follow the Relay Connection pattern for cursor-based pagination
4. **Keep Types in Sync**: Run the code generator whenever the schema changes
5. **Document Schema Changes**: Update the schema changelog when making changes

## Getting Started

To use these patterns in your own components:

1. Define fragments for your components
2. Use those fragments in your queries
3. Run the code generator to create typed hooks
4. Use the generated hooks in your components

For more details, see the [Schema Sync Documentation](../../SCHEMA_SYNC.md).
