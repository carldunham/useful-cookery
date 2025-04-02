package auth

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/carldunham/useful-cookery/internal/database"
	"github.com/carldunham/useful-cookery/internal/model"
)

var (
	ErrInitialUsersExist = errors.New("initial user(s) already exist")
)

// UpdateUser updates a user in the database.
func (s *Service) UpdateUser(ctx context.Context, user *model.User) error {
	err := s.dbClient.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// HashPassword hashes a password using bcrypt.
func (s *Service) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// ValidatePassword checks if a password matches a hash.
func (s *Service) ValidatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// IsAdmin checks if a user has admin privileges.
func (s *Service) IsAdmin(ctx context.Context, userID string) (bool, error) {
	user, err := s.dbClient.GetUser(ctx, userID)
	if err != nil {
		if err == database.ErrNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	return user.Role == model.AdminRole, nil
}

// CreateInitialAdminUser creates an admin user if no users exist.
func (s *Service) CreateInitialAdminUser(ctx context.Context, email, password string) error {
	// Check if any users exist
	users, err := s.dbClient.GetUsers(ctx, 1, 0)
	if err != nil {
		return fmt.Errorf("failed to check existing users: %w", err)
	}

	// If users exist, skip
	if len(users) > 0 {
		return ErrInitialUsersExist
	}

	// Create admin user
	_, err = s.RegisterUser(ctx, "Admin", email, password)
	if err != nil {
		return fmt.Errorf("failed to register admin user: %w", err)
	}

	// Update role to ADMIN
	admin, err := s.dbClient.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to get admin user by email: %w", err)
	}

	admin.Role = model.AdminRole
	err = s.dbClient.UpdateUser(ctx, admin)
	if err != nil {
		return fmt.Errorf("failed to update admin role: %w", err)
	}
	return nil
}

// CheckPermissionForRecipe checks if a user has permission to modify a recipe.
func (s *Service) CheckPermissionForRecipe(ctx context.Context, userID, recipeID string) (bool, error) {
	// Get recipe
	recipe, err := s.dbClient.GetRecipe(ctx, recipeID)
	if err != nil {
		return false, fmt.Errorf("failed to get recipe: %w", err)
	}

	// Check if user is admin
	isAdmin, err := s.IsAdmin(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to check admin status: %w", err)
	}

	// Admins can modify any recipe
	if isAdmin {
		return true, nil
	}

	// Check if user is the author
	return recipe.Author != nil && recipe.Author.ID == userID, nil
}
