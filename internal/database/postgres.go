package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver.
	"github.com/pkg/errors"

	"github.com/carldunham/useful-cookery/internal/database/dbtypes"
	"github.com/carldunham/useful-cookery/internal/model"
)

// Database connection pool constants.
const (
	MaxOpenConns    = 25
	MaxIdleConns    = 5
	ConnMaxLifetime = 5 * time.Minute
	CacheDefaultTTL = 5 * time.Minute
)

// Error constants.
var (
	ErrRecipeIDRequired = errors.New("recipe ID is required")
	ErrAuthorIDRequired = errors.New("author ID is required")
)

// PostgresDatabase implements database operations using PostgreSQL.
// It is exported for testing purposes.
type PostgresDatabase struct {
	// DB is exported for testing purposes.
	DB    *sql.DB
	cache dbtypes.Cache
}

// NewPostgresDatabase creates a new PostgreSQL database instance.
func NewPostgresDatabase(options dbtypes.DatabaseOptions) (*PostgresDatabase, error) {
	db, err := sql.Open("postgres", options.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("opening database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	// Set connection pool parameters
	db.SetMaxOpenConns(MaxOpenConns)
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetConnMaxLifetime(ConnMaxLifetime)

	pgDB := &PostgresDatabase{
		DB: db,
	}

	// Set up cache if enabled
	if options.CacheEnabled && options.Cache != nil {
		pgDB.cache = options.Cache
	}

	// Initialize database schema if needed
	if err := pgDB.initSchema(); err != nil {
		return nil, fmt.Errorf("initializing schema: %w", err)
	}

	return pgDB, nil
}

// Close closes the database connection.
func (db *PostgresDatabase) Close() error {
	if err := db.DB.Close(); err != nil {
		return fmt.Errorf("closing database connection: %w", err)
	}
	return nil
}

// initSchema initializes the database schema if it doesn't exist.
//
//nolint:cyclop,funlen // This function is necessarily complex and long due to the database schema creation
func (db *PostgresDatabase) initSchema() error {
	// Create users table
	_, err := db.DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			role TEXT NOT NULL,
			preferences JSONB,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("creating users table: %w", err)
	}

	// Create categories table
	_, err = db.DB.Exec(`
		CREATE TABLE IF NOT EXISTS categories (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("creating categories table: %w", err)
	}

	// Create recipes table
	_, err = db.DB.Exec(`
		CREATE TABLE IF NOT EXISTS recipes (
			id TEXT PRIMARY KEY,
			original_id TEXT,
			title TEXT NOT NULL,
			description TEXT,
			notes TEXT,
			author_id TEXT REFERENCES users(id),
			cuisine TEXT,
			prep_time INTEGER,
			cook_time INTEGER,
			servings INTEGER,
			difficulty TEXT,
			nutrition_info JSONB,
			tags TEXT[],
			likes INTEGER DEFAULT 0,
			average_rating FLOAT DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("creating recipes table: %w", err)
	}

	// Create recipe_categories junction table
	_, err = db.DB.Exec(`
		CREATE TABLE IF NOT EXISTS recipe_categories (
			recipe_id TEXT REFERENCES recipes(id) ON DELETE CASCADE,
			category_id TEXT REFERENCES categories(id) ON DELETE CASCADE,
			PRIMARY KEY (recipe_id, category_id)
		)
	`)
	if err != nil {
		return fmt.Errorf("creating recipe_categories table: %w", err)
	}

	// Create ingredients table
	_, err = db.DB.Exec(`
		CREATE TABLE IF NOT EXISTS ingredients (
			id TEXT PRIMARY KEY,
			recipe_id TEXT REFERENCES recipes(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			quantity FLOAT,
			unit TEXT,
			preparation TEXT,
			substitutes TEXT[],
			is_optional BOOLEAN DEFAULT FALSE
		)
	`)
	if err != nil {
		return fmt.Errorf("creating ingredients table: %w", err)
	}

	// Create steps table
	_, err = db.DB.Exec(`
		CREATE TABLE IF NOT EXISTS steps (
			id TEXT PRIMARY KEY,
			recipe_id TEXT REFERENCES recipes(id) ON DELETE CASCADE,
			order_index INTEGER NOT NULL,
			description TEXT NOT NULL,
			time_estimate INTEGER
		)
	`)
	if err != nil {
		return fmt.Errorf("creating steps table: %w", err)
	}

	// Create images table
	_, err = db.DB.Exec(`
		CREATE TABLE IF NOT EXISTS images (
			id TEXT PRIMARY KEY,
			url TEXT NOT NULL,
			alt TEXT,
			width INTEGER,
			height INTEGER
		)
	`)
	if err != nil {
		return fmt.Errorf("creating images table: %w", err)
	}

	// Create recipe_images junction table
	_, err = db.DB.Exec(`
		CREATE TABLE IF NOT EXISTS recipe_images (
			recipe_id TEXT REFERENCES recipes(id) ON DELETE CASCADE,
			image_id TEXT REFERENCES images(id) ON DELETE CASCADE,
			PRIMARY KEY (recipe_id, image_id)
		)
	`)
	if err != nil {
		return fmt.Errorf("creating recipe_images table: %w", err)
	}

	// Create step_images junction table
	_, err = db.DB.Exec(`
		CREATE TABLE IF NOT EXISTS step_images (
			step_id TEXT REFERENCES steps(id) ON DELETE CASCADE,
			image_id TEXT REFERENCES images(id) ON DELETE CASCADE,
			PRIMARY KEY (step_id, image_id)
		)
	`)
	if err != nil {
		return fmt.Errorf("creating step_images table: %w", err)
	}

	// Create reviews table
	_, err = db.DB.Exec(`
		CREATE TABLE IF NOT EXISTS reviews (
			id TEXT PRIMARY KEY,
			recipe_id TEXT REFERENCES recipes(id) ON DELETE CASCADE,
			author_id TEXT REFERENCES users(id) ON DELETE CASCADE,
			rating INTEGER NOT NULL,
			comment TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("creating reviews table: %w", err)
	}

	// Create saved_recipes junction table
	_, err = db.DB.Exec(`
		CREATE TABLE IF NOT EXISTS saved_recipes (
			user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
			recipe_id TEXT REFERENCES recipes(id) ON DELETE CASCADE,
			PRIMARY KEY (user_id, recipe_id)
		)
	`)
	if err != nil {
		return fmt.Errorf("creating saved_recipes table: %w", err)
	}

	return nil
}

// Query executes a raw SQL query against PostgreSQL.
// This is a simplified implementation that doesn't support all DGraph query features.
func (db *PostgresDatabase) Query(_ context.Context, _ string, _ map[string]string, _ any) error {
	return errors.New("raw queries not supported in PostgreSQL implementation")
}

// GetUser fetches a user by ID.
//
//nolint:cyclop // This function is necessarily complex due to caching and preference handling
func (db *PostgresDatabase) GetUser(ctx context.Context, userID string) (*model.User, error) {
	if userID == "" {
		return nil, dbtypes.ErrInvalidID
	}

	// Check cache first if enabled
	if db.cache != nil {
		cacheKey := "user:" + userID
		cachedData, err := db.cache.Get(ctx, cacheKey)
		if err == nil && cachedData != nil {
			var user model.User
			if err := json.Unmarshal(cachedData, &user); err == nil {
				return &user, nil
			}
		}
	}

	query := `
		SELECT id, name, email, password, role, preferences, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var (
		user        model.User
		preferences []byte
	)

	err := db.DB.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&preferences,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dbtypes.ErrNotFound
		}
		return nil, fmt.Errorf("querying user: %w", err)
	}

	// Parse preferences if present
	if len(preferences) > 0 {
		user.Preferences = &model.UserPreferences{}
		if err := json.Unmarshal(preferences, user.Preferences); err != nil {
			return nil, fmt.Errorf("unmarshaling preferences: %w", err)
		}
	}

	// Cache the result if caching is enabled
	if db.cache != nil {
		userData, err := json.Marshal(user)
		if err == nil {
			cacheKey := "user:" + userID
			_ = db.cache.Set(ctx, cacheKey, userData, CacheDefaultTTL)
		}
	}

	return &user, nil
}

// GetUserByEmail fetches a user by email.
func (db *PostgresDatabase) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	if email == "" {
		return nil, dbtypes.ErrEmailRequired
	}

	query := `
		SELECT id, name, email, password, role, preferences, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var (
		user        model.User
		preferences []byte
	)

	err := db.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&preferences,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dbtypes.ErrNotFound
		}
		return nil, fmt.Errorf("querying user by email: %w", err)
	}

	// Parse preferences if present
	if len(preferences) > 0 {
		user.Preferences = &model.UserPreferences{}
		if err := json.Unmarshal(preferences, user.Preferences); err != nil {
			return nil, fmt.Errorf("unmarshaling preferences: %w", err)
		}
	}

	return &user, nil
}

// CreateUser creates a new user.
func (db *PostgresDatabase) CreateUser(ctx context.Context, user *model.User) error {
	if user.Email == "" {
		return dbtypes.ErrEmailRequired
	}

	// Check if user already exists
	existingUser, err := db.GetUserByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return dbtypes.ErrAlreadyExists
	}

	// Generate ID if not provided
	if user.ID == "" {
		user.ID = model.NewID()
	}

	// Set timestamps
	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	user.UpdatedAt = now

	// Marshal preferences to JSON
	var preferencesJSON []byte
	var err2 error
	if user.Preferences != nil {
		preferencesJSON, err2 = json.Marshal(user.Preferences)
		if err2 != nil {
			return fmt.Errorf("marshaling preferences: %w", err2)
		}
	}

	query := `
		INSERT INTO users (id, name, email, password, role, preferences, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err = db.DB.ExecContext(
		ctx,
		query,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.Role,
		preferencesJSON,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return dbtypes.ErrAlreadyExists
		}
		return fmt.Errorf("inserting user: %w", err)
	}

	return nil
}

// GetUsers fetches a list of users with pagination.
func (db *PostgresDatabase) GetUsers(ctx context.Context, limit, offset int) ([]*model.User, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}

	if offset < 0 {
		offset = 0
	}

	query := `
		SELECT id, name, email, role, preferences, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := db.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("querying users: %w", err)
	}
	defer rows.Close()

	var users []*model.User

	for rows.Next() {
		var (
			user        model.User
			preferences []byte
		)

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Role,
			&preferences,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scanning user row: %w", err)
		}

		// Parse preferences if present
		if len(preferences) > 0 {
			user.Preferences = &model.UserPreferences{}
			if err := json.Unmarshal(preferences, user.Preferences); err != nil {
				return nil, fmt.Errorf("unmarshaling preferences: %w", err)
			}
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating user rows: %w", err)
	}

	return users, nil
}

// UpdateUser updates a user.
func (db *PostgresDatabase) UpdateUser(ctx context.Context, user *model.User) error {
	if user.ID == "" {
		return dbtypes.ErrInvalidID
	}

	// Check if user exists
	_, err := db.GetUser(ctx, user.ID)
	if err != nil {
		return err
	}

	// Update timestamp
	user.UpdatedAt = time.Now()

	// Marshal preferences to JSON
	var preferencesJSON []byte
	var err2 error
	if user.Preferences != nil {
		preferencesJSON, err2 = json.Marshal(user.Preferences)
		if err2 != nil {
			return fmt.Errorf("marshaling preferences: %w", err2)
		}
	}

	query := `
		UPDATE users
		SET name = $1, email = $2, password = $3, role = $4, preferences = $5, updated_at = $6
		WHERE id = $7
	`

	result, err := db.DB.ExecContext(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Password,
		user.Role,
		preferencesJSON,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return dbtypes.ErrAlreadyExists
		}
		return fmt.Errorf("updating user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return dbtypes.ErrNotFound
	}

	// Invalidate cache if caching is enabled
	if db.cache != nil {
		cacheKey := "user:" + user.ID
		_ = db.cache.Delete(ctx, cacheKey)
	}

	return nil
}

// GetCategory fetches a category by ID.
func (db *PostgresDatabase) GetCategory(ctx context.Context, categoryID string) (*model.Category, error) {
	if categoryID == "" {
		return nil, dbtypes.ErrInvalidID
	}

	// Check cache first if enabled
	if db.cache != nil {
		cacheKey := "category:" + categoryID
		cachedData, err := db.cache.Get(ctx, cacheKey)
		if err == nil && cachedData != nil {
			var category model.Category
			if err := json.Unmarshal(cachedData, &category); err == nil {
				return &category, nil
			}
		}
	}

	query := `
		SELECT id, name, description, created_at, updated_at
		FROM categories
		WHERE id = $1
	`

	var category model.Category

	err := db.DB.QueryRowContext(ctx, query, categoryID).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dbtypes.ErrNotFound
		}
		return nil, fmt.Errorf("querying category: %w", err)
	}

	// Cache the result if caching is enabled
	if db.cache != nil {
		categoryData, err := json.Marshal(category)
		if err == nil {
			cacheKey := "category:" + categoryID
			_ = db.cache.Set(ctx, cacheKey, categoryData, CacheDefaultTTL)
		}
	}

	return &category, nil
}

// CreateCategory creates a new category.
func (db *PostgresDatabase) CreateCategory(ctx context.Context, category *model.Category) error {
	// Generate ID if not provided
	if category.ID == "" {
		category.ID = model.NewID()
	}

	// Set timestamps
	now := time.Now()
	if category.CreatedAt.IsZero() {
		category.CreatedAt = now
	}
	category.UpdatedAt = now

	query := `
		INSERT INTO categories (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := db.DB.ExecContext(
		ctx,
		query,
		category.ID,
		category.Name,
		category.Description,
		category.CreatedAt,
		category.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("inserting category: %w", err)
	}

	return nil
}

// UpdateCategory updates an existing category.
func (db *PostgresDatabase) UpdateCategory(ctx context.Context, category *model.Category) error {
	if category.ID == "" {
		return dbtypes.ErrInvalidID
	}

	// Check if category exists
	_, err := db.GetCategory(ctx, category.ID)
	if err != nil {
		return err
	}

	// Update timestamp
	category.UpdatedAt = time.Now()

	query := `
		UPDATE categories
		SET name = $1, description = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := db.DB.ExecContext(
		ctx,
		query,
		category.Name,
		category.Description,
		category.UpdatedAt,
		category.ID,
	)

	if err != nil {
		return fmt.Errorf("updating category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return dbtypes.ErrNotFound
	}

	// Invalidate cache if caching is enabled
	if db.cache != nil {
		cacheKey := "category:" + category.ID
		_ = db.cache.Delete(ctx, cacheKey)
	}

	return nil
}

// DeleteCategory deletes a category.
func (db *PostgresDatabase) DeleteCategory(ctx context.Context, categoryID string) error {
	if categoryID == "" {
		return dbtypes.ErrInvalidID
	}

	// Check if category exists
	_, err := db.GetCategory(ctx, categoryID)
	if err != nil {
		return err
	}

	query := `DELETE FROM categories WHERE id = $1`

	result, err := db.DB.ExecContext(ctx, query, categoryID)
	if err != nil {
		return fmt.Errorf("deleting category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return dbtypes.ErrNotFound
	}

	// Invalidate cache if caching is enabled
	if db.cache != nil {
		cacheKey := "category:" + categoryID
		_ = db.cache.Delete(ctx, cacheKey)
	}

	return nil
}

// CreateReview creates a new review.
//
//nolint:cyclop // This function is necessarily complex due to transaction handling
func (db *PostgresDatabase) CreateReview(ctx context.Context, review *model.Review) error {
	if review.Recipe == nil || review.Recipe.ID == "" {
		return ErrRecipeIDRequired
	}

	if review.Author == nil || review.Author.ID == "" {
		return ErrAuthorIDRequired
	}

	// Generate ID if not provided
	if review.ID == "" {
		review.ID = model.NewID()
	}

	// Set timestamps
	now := time.Now()
	if review.CreatedAt.IsZero() {
		review.CreatedAt = now
	}
	review.UpdatedAt = now

	// Start a transaction
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Insert review
	query := `
		INSERT INTO reviews (id, recipe_id, author_id, rating, comment, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = tx.ExecContext(
		ctx,
		query,
		review.ID,
		review.Recipe.ID,
		review.Author.ID,
		review.Rating,
		review.Comment,
		review.CreatedAt,
		review.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("inserting review: %w", err)
	}

	// Update recipe's average rating
	err = db.updateRecipeRating(ctx, tx, review.Recipe.ID)
	if err != nil {
		return fmt.Errorf("updating recipe rating: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

// GetReview fetches a review by ID.
func (db *PostgresDatabase) GetReview(ctx context.Context, reviewID string) (*model.Review, error) {
	if reviewID == "" {
		return nil, dbtypes.ErrInvalidID
	}

	query := `
		SELECT r.id, r.rating, r.comment, r.created_at, r.updated_at,
			   a.id, a.name,
			   rc.id, rc.title
		FROM reviews r
		JOIN users a ON r.author_id = a.id
		JOIN recipes rc ON r.recipe_id = rc.id
		WHERE r.id = $1
	`

	var review model.Review
	var author model.User
	var recipe model.Recipe

	err := db.DB.QueryRowContext(ctx, query, reviewID).Scan(
		&review.ID,
		&review.Rating,
		&review.Comment,
		&review.CreatedAt,
		&review.UpdatedAt,
		&author.ID,
		&author.Name,
		&recipe.ID,
		&recipe.Title,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dbtypes.ErrNotFound
		}
		return nil, fmt.Errorf("querying review: %w", err)
	}

	review.Author = &author
	review.Recipe = &recipe

	return &review, nil
}

// UpdateReview updates an existing review.
func (db *PostgresDatabase) UpdateReview(ctx context.Context, review *model.Review) error {
	if review.ID == "" {
		return dbtypes.ErrInvalidID
	}

	// Get existing review to check if it exists and get recipe ID
	existingReview, err := db.GetReview(ctx, review.ID)
	if err != nil {
		return err
	}

	// Update timestamp
	review.UpdatedAt = time.Now()

	// Start a transaction
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Update review
	query := `
		UPDATE reviews
		SET rating = $1, comment = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := tx.ExecContext(
		ctx,
		query,
		review.Rating,
		review.Comment,
		review.UpdatedAt,
		review.ID,
	)

	if err != nil {
		return fmt.Errorf("updating review: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return dbtypes.ErrNotFound
	}

	// Update recipe's average rating
	err = db.updateRecipeRating(ctx, tx, existingReview.Recipe.ID)
	if err != nil {
		return fmt.Errorf("updating recipe rating: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

// DeleteReview deletes a review.
func (db *PostgresDatabase) DeleteReview(ctx context.Context, reviewID string) error {
	if reviewID == "" {
		return dbtypes.ErrInvalidID
	}

	// Get existing review to check if it exists and get recipe ID
	existingReview, err := db.GetReview(ctx, reviewID)
	if err != nil {
		return err
	}

	// Start a transaction
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Delete review
	query := `DELETE FROM reviews WHERE id = $1`

	result, err := tx.ExecContext(ctx, query, reviewID)
	if err != nil {
		return fmt.Errorf("deleting review: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return dbtypes.ErrNotFound
	}

	// Update recipe's average rating
	err = db.updateRecipeRating(ctx, tx, existingReview.Recipe.ID)
	if err != nil {
		return fmt.Errorf("updating recipe rating: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

// updateRecipeRating updates a recipe's average rating based on its reviews.
func (db *PostgresDatabase) updateRecipeRating(ctx context.Context, tx *sql.Tx, recipeID string) error {
	query := `
		UPDATE recipes
		SET average_rating = (
			SELECT COALESCE(AVG(rating), 0)
			FROM reviews
			WHERE recipe_id = $1
		)
		WHERE id = $1
	`

	_, err := tx.ExecContext(ctx, query, recipeID)
	if err != nil {
		return fmt.Errorf("updating recipe rating: %w", err)
	}

	return nil
}

// GetRecipe fetches a recipe by ID.
func (db *PostgresDatabase) GetRecipe(ctx context.Context, recipeID string) (*model.Recipe, error) {
	if recipeID == "" {
		return nil, dbtypes.ErrInvalidID
	}

	// This is a simplified implementation
	query := `
		SELECT id, title, description
		FROM recipes
		WHERE id = $1
	`

	var recipe model.Recipe
	err := db.DB.QueryRowContext(ctx, query, recipeID).Scan(
		&recipe.ID,
		&recipe.Title,
		&recipe.Description,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dbtypes.ErrNotFound
		}
		return nil, fmt.Errorf("querying recipe: %w", err)
	}

	return &recipe, nil
}

// GetRecipes fetches recipes based on filters.
//
//nolint:cyclop,funlen // This function is necessarily complex and long due to dynamic filter handling
func (db *PostgresDatabase) GetRecipes(
	ctx context.Context,
	filter map[string]string,
	first, offset int,
) ([]*model.Recipe, error) {
	if first <= 0 {
		first = 10 // Default limit
	}

	if offset < 0 {
		offset = 0
	}

	// Build query with filters
	query := `
		SELECT r.id, r.title, r.description, r.cuisine, r.created_at,
		       u.id, u.name
		FROM recipes r
		LEFT JOIN users u ON r.author_id = u.id
		WHERE 1=1
	`
	args := make([]interface{}, 0)
	argIndex := 1

	// Apply filters
	for key, value := range filter {
		if value == "" {
			continue
		}

		switch key {
		case "title":
			query += fmt.Sprintf(" AND r.title ILIKE $%d", argIndex)
			args = append(args, "%"+value+"%")
			argIndex++
		case "cuisine":
			query += fmt.Sprintf(" AND r.cuisine ILIKE $%d", argIndex)
			args = append(args, "%"+value+"%")
			argIndex++
		case "category":
			query += fmt.Sprintf(` AND r.id IN (
				SELECT recipe_id FROM recipe_categories rc
				JOIN categories c ON rc.category_id = c.id
				WHERE c.name ILIKE $%d
			)`, argIndex)
			args = append(args, "%"+value+"%")
			argIndex++
		case "authorID":
			query += fmt.Sprintf(" AND r.author_id = $%d", argIndex)
			args = append(args, value)
			argIndex++
		}
	}

	// Add order by and pagination
	query += " ORDER BY created_at DESC LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
	args = append(args, first, offset)

	// Execute query
	rows, err := db.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying recipes: %w", err)
	}
	defer rows.Close()

	var recipes []*model.Recipe
	for rows.Next() {
		var (
			recipe     model.Recipe
			cuisine    sql.NullString
			authorID   sql.NullString
			authorName sql.NullString
		)

		err := rows.Scan(
			&recipe.ID,
			&recipe.Title,
			&recipe.Description,
			&cuisine,
			&recipe.CreatedAt,
			&authorID,
			&authorName,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning recipe row: %w", err)
		}

		// Set nullable fields
		if cuisine.Valid {
			recipe.Cuisine = cuisine.String
		}

		// Set author if present
		if authorID.Valid && authorName.Valid {
			recipe.Author = &model.User{
				ID:   authorID.String,
				Name: authorName.String,
			}
		}

		recipes = append(recipes, &recipe)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating recipe rows: %w", err)
	}

	return recipes, nil
}

// CountRecipes returns the total count of recipes matching the given filters.
func (db *PostgresDatabase) CountRecipes(ctx context.Context, filter map[string]string) (int, error) {
	// Build query with filters
	query := `
		SELECT COUNT(*)
		FROM recipes r
		WHERE 1=1
	`
	args := make([]interface{}, 0)
	argIndex := 1

	// Apply filters
	for key, value := range filter {
		if value == "" {
			continue
		}

		switch key {
		case "title":
			query += fmt.Sprintf(" AND r.title ILIKE $%d", argIndex)
			args = append(args, "%"+value+"%")
			argIndex++
		case "cuisine":
			query += fmt.Sprintf(" AND r.cuisine ILIKE $%d", argIndex)
			args = append(args, "%"+value+"%")
			argIndex++
		case "category":
			query += fmt.Sprintf(` AND r.id IN (
				SELECT recipe_id FROM recipe_categories rc
				JOIN categories c ON rc.category_id = c.id
				WHERE c.name ILIKE $%d
			)`, argIndex)
			args = append(args, "%"+value+"%")
			argIndex++
		case "authorID":
			query += fmt.Sprintf(" AND r.author_id = $%d", argIndex)
			args = append(args, value)
			argIndex++
		}
	}

	// Execute query
	var count int
	err := db.DB.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting recipes: %w", err)
	}

	return count, nil
}

// CreateRecipe creates a new recipe.
func (db *PostgresDatabase) CreateRecipe(ctx context.Context, recipe *model.Recipe) error {
	// Generate ID if not provided
	if recipe.ID == "" {
		recipe.ID = model.NewID()
	}

	// Set timestamps
	now := time.Now()
	if recipe.CreatedAt.IsZero() {
		recipe.CreatedAt = now
	}
	recipe.UpdatedAt = now

	// Basic implementation - just insert the recipe
	query := `
		INSERT INTO recipes (id, title, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := db.DB.ExecContext(
		ctx,
		query,
		recipe.ID,
		recipe.Title,
		recipe.Description,
		recipe.CreatedAt,
		recipe.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("inserting recipe: %w", err)
	}

	return nil
}

// UpdateRecipe updates an existing recipe.
func (db *PostgresDatabase) UpdateRecipe(ctx context.Context, recipe *model.Recipe) error {
	if recipe.ID == "" {
		return dbtypes.ErrInvalidID
	}

	// Update timestamp
	recipe.UpdatedAt = time.Now()

	// Basic implementation - just update the recipe
	query := `
		UPDATE recipes
		SET title = $1, description = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := db.DB.ExecContext(
		ctx,
		query,
		recipe.Title,
		recipe.Description,
		recipe.UpdatedAt,
		recipe.ID,
	)

	if err != nil {
		return fmt.Errorf("updating recipe: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return dbtypes.ErrNotFound
	}

	return nil
}

// DeleteRecipe deletes a recipe.
func (db *PostgresDatabase) DeleteRecipe(ctx context.Context, recipeID string) error {
	if recipeID == "" {
		return dbtypes.ErrInvalidID
	}

	query := `DELETE FROM recipes WHERE id = $1`

	result, err := db.DB.ExecContext(ctx, query, recipeID)
	if err != nil {
		return fmt.Errorf("deleting recipe: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return dbtypes.ErrNotFound
	}

	return nil
}

// GetPopularRecipes returns popular recipes for PostgreSQL.
//
//nolint:funlen // This function is necessarily long due to the query and result processing
func (db *PostgresDatabase) GetPopularRecipes(ctx context.Context, limit, offset int) ([]*model.Recipe, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}

	if offset < 0 {
		offset = 0
	}

	query := `
		SELECT r.id, r.title, r.description, r.cuisine, r.created_at,
		       u.id, u.name
		FROM recipes r
		LEFT JOIN users u ON r.author_id = u.id
		ORDER BY r.likes DESC, r.average_rating DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := db.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("querying popular recipes: %w", err)
	}
	defer rows.Close()

	var recipes []*model.Recipe
	for rows.Next() {
		var (
			recipe     model.Recipe
			cuisine    sql.NullString
			authorID   sql.NullString
			authorName sql.NullString
		)

		err := rows.Scan(
			&recipe.ID,
			&recipe.Title,
			&recipe.Description,
			&cuisine,
			&recipe.CreatedAt,
			&authorID,
			&authorName,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning recipe row: %w", err)
		}

		// Set nullable fields
		if cuisine.Valid {
			recipe.Cuisine = cuisine.String
		}

		// Set author if present
		if authorID.Valid && authorName.Valid {
			recipe.Author = &model.User{
				ID:   authorID.String,
				Name: authorName.String,
			}
		}

		recipes = append(recipes, &recipe)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating recipe rows: %w", err)
	}

	return recipes, nil
}
