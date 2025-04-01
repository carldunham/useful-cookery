package auth

import (
	"context"
	"fmt"

	"github.com/carldunham/useful-cookery/internal/database"
	"github.com/carldunham/useful-cookery/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// UpdateUser updates a user in the database
func (s *Service) UpdateUser(ctx context.Context, user *model.User) error {
	return s.dbClient.UpdateUser(ctx, user)
}

// HashPassword hashes a password using bcrypt
func (s *Service) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// ValidatePassword checks if a password matches a hash
func (s *Service) ValidatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// IsAdmin checks if a user has admin privileges
func (s *Service) IsAdmin(ctx context.Context, userID string) (bool, error) {
	user, err := s.dbClient.GetUser(ctx, userID)
	if err != nil {
		if err == database.ErrNotFound {
			return false, nil
		}
		return false, err
	}

	return user.Role == "ADMIN", nil
}

// CreateInitialAdminUser creates an admin user if no users exist
func (s *Service) CreateInitialAdminUser(ctx context.Context, email, password string) error {
	// Check if any users exist
	users, err := s.dbClient.GetUsers(ctx, 1, 0)
	if err != nil {
		return err
	}

	// If users exist, skip
	if len(users) > 0 {
		return nil
	}

	// Create admin user
	_, err = s.RegisterUser(ctx, "Admin", email, password)
	if err != nil {
		return err
	}

	// Update role to ADMIN
	admin, err := s.dbClient.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	admin.Role = "ADMIN"
	return s.dbClient.UpdateUser(ctx, admin)
}

// CheckPermissionForRecipe checks if a user has permission to modify a recipe
func (s *Service) CheckPermissionForRecipe(ctx context.Context, userID, recipeID string) (bool, error) {
	// Get recipe
	recipe, err := s.dbClient.GetRecipe(ctx, recipeID)
	if err != nil {
		return false, err
	}

	// Check if user is admin
	isAdmin, err := s.IsAdmin(ctx, userID)
	if err != nil {
		return false, err
	}

	// Admins can modify any recipe
	if isAdmin {
		return true, nil
	}

	// Check if user is the author
	return recipe.Author != nil && recipe.Author.ID == userID, nil
}
