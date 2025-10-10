package services

import (
	"fmt"
	"template-store/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService handles authentication operations
type AuthService struct {
	db         *gorm.DB
	jwtService *JWTService
}

// NewAuthService creates a new authentication service
func NewAuthService(db *gorm.DB, jwtService *JWTService) *AuthService {
	return &AuthService{
		db:         db,
		jwtService: jwtService,
	}
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	Token string          `json:"token"`
	User  UserInfoResponse `json:"user"`
}

// UserInfoResponse represents user information in auth response
type UserInfoResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Register creates a new user account
func (s *AuthService) Register(req RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	} else if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "customer", // Default role
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Name, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResponse{
		Token: token,
		User: UserInfoResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(req LoginRequest) (*AuthResponse, error) {
	// Find user by email
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid email or password")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Update last login time
	s.db.Model(&user).Update("last_login_at", gorm.Expr("NOW()"))

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Name, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResponse{
		Token: token,
		User: UserInfoResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

// GetUserFromToken validates token and returns user info
func (s *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	// Validate token
	claims, err := s.jwtService.ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Get user from database
	var user models.User
	if err := s.db.First(&user, claims.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// RefreshToken generates a new token for the user
func (s *AuthService) RefreshToken(tokenString string) (string, error) {
	return s.jwtService.RefreshToken(tokenString)
}

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	// Get user
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return fmt.Errorf("invalid old password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	if err := s.db.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// ResetPasswordRequest initiates password reset (to be implemented with email)
func (s *AuthService) ResetPasswordRequest(email string) error {
	// Check if user exists
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Don't reveal if email exists or not for security
			return nil
		}
		return fmt.Errorf("failed to check user: %w", err)
	}

	// TODO: Generate reset token
	// TODO: Send reset email via SendGrid
	// For now, just return success
	return nil
}

// ValidateResetToken validates a password reset token
func (s *AuthService) ValidateResetToken(token string) (uint, error) {
	// TODO: Implement reset token validation
	// For now, just validate as JWT
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return 0, fmt.Errorf("invalid reset token")
	}
	return claims.UserID, nil
}

// ResetPassword resets a user's password using a reset token
func (s *AuthService) ResetPassword(token, newPassword string) error {
	// Validate token
	userID, err := s.ValidateResetToken(token)
	if err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	if err := s.db.Model(&models.User{}).Where("id = ?", userID).Update("password", string(hashedPassword)).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}
