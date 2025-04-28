# GraphQL Schema Changelog

This file documents changes to the GraphQL schema to help track when and why changes were made.

## [Unreleased]

## [1.0.0] - 2025-04-28

### Added

- Implemented Relay Connection pattern for pagination in the following queries:
  - `recipes` query
  - `searchRecipes` query
  - `recommendRecipes` query
- Added `PageInfo` type with `hasNextPage`, `hasPreviousPage`, `startCursor`, and `endCursor` fields
- Added `RecipeEdge` type with `node` and `cursor` fields
- Added `RecipeConnection` type with `edges`, `pageInfo`, and `totalCount` fields

### Changed

- Updated `recipes` query to return a `RecipeConnection` instead of a list of `Recipe`
- Updated `searchRecipes` query to return a `RecipeConnection` instead of a list of `Recipe`
- Updated `recommendRecipes` query to return a `RecipeConnection` instead of a list of `Recipe`
- Added pagination parameters (`first` and `after`) to all queries that return a `RecipeConnection`

### Removed

- Removed direct list return types from `recipes`, `searchRecipes`, and `recommendRecipes` queries

## How to Update Client Code

When the schema changes to use Relay Connections, client code needs to be updated to:

1. Update GraphQL queries to request the new connection structure:

   ```graphql
   query {
     recipes {
       edges {
         node {
           id
           title
           # other fields
         }
         cursor
       }
       pageInfo {
         hasNextPage
         hasPreviousPage
         startCursor
         endCursor
       }
       totalCount
     }
   }
   ```

2. Update component logic to handle the new data structure:

   ```javascript
   // Before
   const recipes = data?.recipes || [];

   // After
   const recipes = data?.recipes?.edges?.map(edge => edge.node) || [];
   const pageInfo = data?.recipes?.pageInfo;
   const totalCount = data?.recipes?.totalCount;
   ```

3. Implement pagination using the cursor-based approach:

   ```javascript
   const loadMore = () => {
     fetchMore({
       variables: {
         after: pageInfo.endCursor,
         first: 10
       },
     });
   };
