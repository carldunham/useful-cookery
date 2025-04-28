# Database Abstraction Layer

This package provides a database abstraction layer for the Useful Cookery application. It allows the application to use different database backends without changing the application code.

## Architecture

The database abstraction layer consists of the following components:

- **Database Interfaces**: Defined in caller packages, specifying only the operations they need.
- **Database Implementations**: Concrete implementations of the database interfaces for different database backends.
- **Factory Functions**: Functions to create database instances based on configuration.

## Directory Structure

- `dbtypes/`: Contains the full database interface and common types and errors used by all database implementations.
- `strategies/`: Contains concrete implementations of the database interfaces.
- `factory.go`: Contains factory functions to create database instances.
- `interfaces.go`: Re-exports common types and errors from dbtypes.
- `cache.go`: Provides caching functionality for database operations.
- `dgraph.go`: Legacy DGraph client implementation (will be deprecated).

## Usage

### Configuration

The database type and connection string are specified in the application configuration:

```yaml
database:
  type: dgraph # or "memory" for in-memory database
  connection_string: localhost:9080
```

For backward compatibility, the application also supports the old configuration format:

```yaml
dgraph:
  connection_string: localhost:9080
```

### Creating a Database Instance

```go
import (
    "github.com/carldunham/useful-cookery/internal/database"
    "github.com/carldunham/useful-cookery/internal/database/dbtypes"
    "github.com/carldunham/useful-cookery/internal/mypackage" // Your package with interfaces
)

// Create database options
dbOptions := database.Options{
    ConnectionString: cfg.Database.ConnectionString,
    CacheEnabled:     true,
    Cache:            cache,  // Optional cache implementation
}

// Create database instance based on type
var dbImpl dbtypes.Database
var err error

switch cfg.Database.Type {
case "dgraph", "":
    dbImpl, err = database.NewDGraphDatabase(dbOptions)
case "memory":
    dbImpl, err = database.NewInMemoryDatabase(dbOptions)
default:
    log.Fatalf("Unsupported database type: %s", cfg.Database.Type)
}

// Use the database implementation with your specific interface
var db mypackage.Database = dbImpl
```

### Using the Database

```go
// Define your own interface with only the methods you need
type UserDatabase interface {
    GetUser(ctx context.Context, userID string) (*model.User, error)
    CreateUser(ctx context.Context, user *model.User) error
}

// Get a user by ID
user, err := db.GetUser(ctx, userID)
if err != nil {
    // Handle error
}

// Create a new user
err = db.CreateUser(ctx, user)
if err != nil {
    // Handle error
}

// Execute a raw query (if your interface includes it)
var result struct {
    Users []*model.User `json:"users"`
}
err = db.Query(ctx, query, vars, &result)
if err != nil {
    // Handle error
}
```

## Adding a New Database Implementation

To add a new database implementation:

1. Create a new file in the `strategies/` directory.
2. Implement the full `dbtypes.Database` interface.
3. Add a new database type constant in `dbtypes/types.go`.
4. Add a factory function in `factory.go`.
5. Update the switch statement in the application code to use the new database type.

## Error Handling

Common database errors are defined in `dbtypes/errors.go`:

- `ErrNotFound`: Returned when an entity is not found.
- `ErrInvalidID`: Returned when an invalid ID is provided.
- `ErrAlreadyExists`: Returned when an entity already exists.
- `ErrEmailRequired`: Returned when an email is required but not provided.

These errors should be used by all database implementations to ensure consistent error handling.
