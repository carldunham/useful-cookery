# Database Migrations

This directory contains database migration files for the Useful Cookery application. We use [golang-migrate](https://github.com/golang-migrate/migrate) to manage database schema migrations.

## Migration Files

Migration files are stored in the `migrations` directory and follow the naming convention:

```
{version}_{name}.{up|down}.sql
```

- `version`: A numeric version number (e.g., 000001)
- `name`: A descriptive name for the migration
- `up|down`: Indicates whether the migration is for applying changes (up) or reverting them (down)

## Initial Schema

The initial schema includes the following tables:

1. `users`: Stores user information
2. `categories`: Stores recipe categories
3. `recipes`: Stores recipe information
4. `recipe_categories`: Junction table for recipe-category relationships
5. `ingredients`: Stores recipe ingredients
6. `steps`: Stores recipe preparation steps
7. `images`: Stores image information
8. `recipe_images`: Junction table for recipe-image relationships
9. `step_images`: Junction table for step-image relationships
10. `reviews`: Stores recipe reviews
11. `saved_recipes`: Junction table for user-saved recipes

## Usage

### Using Make Commands

The project includes several make commands for working with migrations:

```bash
# Create a new migration
make migrate-create

# Apply migrations
make migrate-up

# Revert migrations
make migrate-down

# Check migration status
make migrate-status
```

### Using the Migration Tool Directly

You can also use the migration tool directly:

```bash
# Create a new migration
go run cmd/migration/main.go db create <name>

# Apply migrations
go run cmd/migration/main.go db up

# Apply a specific number of migrations
go run cmd/migration/main.go db up -n <number>

# Revert migrations
go run cmd/migration/main.go db down

# Revert a specific number of migrations
go run cmd/migration/main.go db down -n <number>

# Check migration status
go run cmd/migration/main.go db status
```

### In Kubernetes

Migrations are automatically applied when deploying to Kubernetes using:

```bash
# Deploy migrations only
make deploy-migrations-local

# Deploy all components (including migrations)
make deploy-local
```

## Creating New Migrations

When creating a new migration, follow these guidelines:

1. Use a descriptive name that clearly indicates what the migration does
2. Always create both up and down migrations
3. Test migrations locally before deploying
4. Keep migrations small and focused on a single change
5. Ensure migrations are idempotent when possible

Example:

```sql
-- 000007_add_user_preferences.up.sql
ALTER TABLE users ADD COLUMN theme VARCHAR(50) DEFAULT 'light';

-- 000007_add_user_preferences.down.sql
ALTER TABLE users DROP COLUMN theme;
```
