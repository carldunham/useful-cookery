package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/carldunham/useful-cookery/internal/database"
	"github.com/carldunham/useful-cookery/internal/model"
)

// Errors.
var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrInvalidToken         = errors.New("invalid token")
	ErrExpiredToken         = errors.New("token expired")
	ErrUserExists           = errors.New("user already exists")
	ErrUnexpectedSignMethod = errors.New("unexpected signing method")
)

// Claims represents JWT claims.
type Claims struct {
	UserID string     `json:"user_id"`
	Role   model.Role `json:"role"`
	jwt.RegisteredClaims
}

// Service provides authentication functionality.
type Service struct {
	jwtSecret   string
	tokenExpiry time.Duration
	dbClient    *database.DGraphClient
}

// NewService creates a new auth service.
func NewService(jwtSecret string, tokenExpiry time.Duration, dbClient *database.DGraphClient) *Service {
	return &Service{
		jwtSecret:   jwtSecret,
		tokenExpiry: tokenExpiry,
		dbClient:    dbClient,
	}
}

// RegisterUser registers a new user.
func (s *Service) RegisterUser(ctx context.Context, name, email, password string) (*model.User, error) {
	// Check if user already exists
	existingUser, err := s.dbClient.GetUserByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &model.User{
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		Role:      "USER", // Default role
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save user to database
	err = s.dbClient.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Don't return the password
	user.Password = ""

	return user, nil
}

// Login authenticates a user.
func (s *Service) Login(ctx context.Context, email, password string) (string, *model.User, error) {
	user, err := s.dbClient.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		log.Printf("Login: getting user: %v", err)
		return "", nil, ErrInvalidCredentials
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Login: comparing hash and password: %v", err)
		return "", nil, ErrInvalidCredentials
	}

	// Generate token
	token, err := s.generateToken(user.ID, user.Role)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Don't return the password
	user.Password = ""

	return token, user, nil
}

// generateToken generates a JWT token.
func (s *Service) generateToken(userID string, role model.Role) (string, error) {
	// Set token expiry
	expirationTime := time.Now().Add(s.tokenExpiry)

	// Create claims
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "useful-cookery",
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// VerifyToken verifies and parses a JWT token.
func (s *Service) VerifyToken(tokenString string) (*Claims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", ErrUnexpectedSignMethod, token.Header["alg"])
		}

		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		var validationErr *jwt.ValidationError
		if errors.As(err, &validationErr) {
			if validationErr.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrExpiredToken
			}
		}
		return nil, ErrInvalidToken
	}

	// Extract claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// AuthMiddleware is an HTTP middleware for authentication.
func (s *Service) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// Get Authorization header
		authHeader := request.Header.Get("Authorization")
		if authHeader == "" {
			// No token, continue without user info
			next(writer, request)
			return
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			next(writer, request)
			return
		}

		// Verify token
		claims, err := s.VerifyToken(parts[1])
		if err != nil {
			// Invalid token, continue without user info
			next(writer, request)
			return
		}

		// Get user from database
		user, err := s.dbClient.GetUser(request.Context(), claims.UserID)
		if err != nil || user == nil {
			// User not found, continue without user info
			next(writer, request)
			return
		}

		// Add user to context
		ctx := context.WithValue(request.Context(), UserContextKey, user)
		ctx = context.WithValue(ctx, RoleContextKey, claims.Role)

		// Continue with the new context
		next(writer, request.WithContext(ctx))
	}
}

// GetUserFromContext gets the user from context.
func GetUserFromContext(ctx context.Context) *model.User {
	user, ok := ctx.Value(UserContextKey).(*model.User)
	if !ok {
		return nil
	}
	return user
}

// GetRoleFromContext gets the role from context.
func GetRoleFromContext(ctx context.Context) string {
	role, ok := ctx.Value(RoleContextKey).(string)
	if !ok {
		return ""
	}
	return role
}

// Context keys.
type contextKey string

const (
	UserContextKey contextKey = "user"
	RoleContextKey contextKey = "role"
)
