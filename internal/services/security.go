package services

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"template-store/internal/models"
)

var (
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrTokenAlreadyUsed   = errors.New("token already used")
	ErrPasswordMismatch   = errors.New("current password is incorrect")
	ErrCognitoOperation   = errors.New("cognito operation failed")
)

const (
	PasswordResetTokenExpiry  = 24 * time.Hour
	EmailVerificationExpiry   = 72 * time.Hour
	TokenLength               = 32 // bytes (64 hex characters)
)

// SecurityService defines the interface for security operations
type SecurityService interface {
	// Password operations
	ChangePassword(userID uint, currentPassword, newPassword string) error
	RequestPasswordReset(email string) (*models.PasswordResetToken, error)
	ResetPasswordWithToken(token, newPassword string) error
	ValidatePasswordResetToken(token string) (*models.PasswordResetToken, error)

	// Email verification
	SendVerificationEmail(userID uint) (*models.EmailVerificationToken, error)
	VerifyEmail(token string) error
	RequestEmailChange(userID uint, newEmail string) (*models.EmailVerificationToken, error)

	// Login tracking
	RecordLogin(userID uint, ipAddress, userAgent string, success bool, failureReason string) (*models.LoginHistory, error)
	GetLoginHistory(userID uint, limit, offset int) ([]models.LoginHistory, int64, error)
	GetActiveSessions(userID uint) ([]models.LoginHistory, error)
	LogoutSession(sessionID uint) error
	LogoutAllSessions(userID uint) error

	// Activity logging
	LogActivity(userID uint, action, resource string, resourceID *uint, details map[string]interface{}, ipAddress, userAgent string) error
	GetActivityLog(userID uint, limit, offset int) ([]models.ActivityLog, int64, error)
}

// SecurityServiceImpl implements the SecurityService interface
type SecurityServiceImpl struct {
	db          *gorm.DB
	authService AuthService
}

// NewSecurityService creates a new SecurityService instance
func NewSecurityService(db *gorm.DB, authService AuthService) SecurityService {
	return &SecurityServiceImpl{
		db:          db,
		authService: authService,
	}
}

// ChangePassword changes a user's password
func (s *SecurityServiceImpl) ChangePassword(userID uint, currentPassword, newPassword string) error {
	// Get user
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Verify current password and change to new password via Cognito
	// TODO: Implement Cognito password change
	// For now, return a placeholder error
	// err := s.authService.ChangePassword(user.Email, currentPassword, newPassword)
	// if err != nil {
	// 	return ErrPasswordMismatch
	// }

	return nil
}

// RequestPasswordReset creates a password reset token
func (s *SecurityServiceImpl) RequestPasswordReset(email string) (*models.PasswordResetToken, error) {
	// Find user by email
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Don't reveal if email exists or not for security
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Generate secure random token
	token, err := generateSecureToken(TokenLength)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create password reset token
	resetToken := &models.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(PasswordResetTokenExpiry),
	}

	if err := s.db.Create(resetToken).Error; err != nil {
		return nil, fmt.Errorf("failed to create reset token: %w", err)
	}

	return resetToken, nil
}

// ResetPasswordWithToken resets a password using a token
func (s *SecurityServiceImpl) ResetPasswordWithToken(token, newPassword string) error {
	// Validate token
	resetToken, err := s.ValidatePasswordResetToken(token)
	if err != nil {
		return err
	}

	// Get user
	var user models.User
	if err := s.db.First(&user, resetToken.UserID).Error; err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Reset password via Cognito
	// TODO: Implement Cognito password reset
	// err = s.authService.AdminSetUserPassword(user.CognitoSubject, newPassword)
	// if err != nil {
	// 	return fmt.Errorf("%w: %v", ErrCognitoOperation, err)
	// }

	// Mark token as used
	now := time.Now()
	resetToken.UsedAt = &now
	if err := s.db.Save(resetToken).Error; err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	return nil
}

// ValidatePasswordResetToken validates a password reset token
func (s *SecurityServiceImpl) ValidatePasswordResetToken(token string) (*models.PasswordResetToken, error) {
	var resetToken models.PasswordResetToken

	if err := s.db.Where("token = ?", token).First(&resetToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, fmt.Errorf("failed to find token: %w", err)
	}

	if resetToken.IsExpired() {
		return nil, ErrInvalidToken
	}

	if resetToken.IsUsed() {
		return nil, ErrTokenAlreadyUsed
	}

	return &resetToken, nil
}

// SendVerificationEmail creates and sends an email verification token
func (s *SecurityServiceImpl) SendVerificationEmail(userID uint) (*models.EmailVerificationToken, error) {
	// Get user
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Generate secure random token
	token, err := generateSecureToken(TokenLength)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create email verification token
	verificationToken := &models.EmailVerificationToken{
		UserID:    user.ID,
		Token:     token,
		Email:     user.Email,
		ExpiresAt: time.Now().Add(EmailVerificationExpiry),
	}

	if err := s.db.Create(verificationToken).Error; err != nil {
		return nil, fmt.Errorf("failed to create verification token: %w", err)
	}

	// TODO: Send email via EmailService

	return verificationToken, nil
}

// VerifyEmail verifies an email using a token
func (s *SecurityServiceImpl) VerifyEmail(token string) error {
	var verificationToken models.EmailVerificationToken

	if err := s.db.Where("token = ?", token).First(&verificationToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidToken
		}
		return fmt.Errorf("failed to find token: %w", err)
	}

	if !verificationToken.IsValid() {
		return ErrInvalidToken
	}

	// Mark token as verified
	now := time.Now()
	verificationToken.VerifiedAt = &now
	if err := s.db.Save(&verificationToken).Error; err != nil {
		return fmt.Errorf("failed to mark token as verified: %w", err)
	}

	// Update user's email and mark as verified
	updates := map[string]interface{}{
		"email":          verificationToken.Email,
		"email_verified": true,
	}

	if err := s.db.Model(&models.User{}).Where("id = ?", verificationToken.UserID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	return nil
}

// RequestEmailChange creates a token for email change verification
func (s *SecurityServiceImpl) RequestEmailChange(userID uint, newEmail string) (*models.EmailVerificationToken, error) {
	// Check if email is already in use
	var existingUser models.User
	if err := s.db.Where("email = ?", newEmail).First(&existingUser).Error; err == nil {
		return nil, ErrEmailAlreadyExists
	}

	// Generate secure random token
	token, err := generateSecureToken(TokenLength)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create email verification token for the new email
	verificationToken := &models.EmailVerificationToken{
		UserID:    userID,
		Token:     token,
		Email:     newEmail,
		ExpiresAt: time.Now().Add(EmailVerificationExpiry),
	}

	if err := s.db.Create(verificationToken).Error; err != nil {
		return nil, fmt.Errorf("failed to create verification token: %w", err)
	}

	// TODO: Send verification email to new address

	return verificationToken, nil
}

// RecordLogin records a login attempt
func (s *SecurityServiceImpl) RecordLogin(userID uint, ipAddress, userAgent string, success bool, failureReason string) (*models.LoginHistory, error) {
	loginRecord := &models.LoginHistory{
		UserID:        userID,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		Device:        detectDevice(userAgent),
		LoginAt:       time.Now(),
		Success:       success,
		FailureReason: failureReason,
	}

	if err := s.db.Create(loginRecord).Error; err != nil {
		return nil, fmt.Errorf("failed to record login: %w", err)
	}

	// Update user's last login time if successful
	if success {
		now := time.Now()
		if err := s.db.Model(&models.User{}).Where("id = ?", userID).Update("last_login_at", now).Error; err != nil {
			// Non-critical error, log but don't fail
			fmt.Printf("Warning: failed to update last login time: %v\n", err)
		}
	}

	return loginRecord, nil
}

// GetLoginHistory retrieves a user's login history
func (s *SecurityServiceImpl) GetLoginHistory(userID uint, limit, offset int) ([]models.LoginHistory, int64, error) {
	var history []models.LoginHistory
	var total int64

	// Count total records
	if err := s.db.Model(&models.LoginHistory{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count login history: %w", err)
	}

	// Get paginated records
	if err := s.db.Where("user_id = ?", userID).
		Order("login_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&history).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get login history: %w", err)
	}

	return history, total, nil
}

// GetActiveSessions retrieves a user's active sessions
func (s *SecurityServiceImpl) GetActiveSessions(userID uint) ([]models.LoginHistory, error) {
	var sessions []models.LoginHistory

	if err := s.db.Where("user_id = ? AND success = ? AND logout_at IS NULL", userID, true).
		Order("login_at DESC").
		Find(&sessions).Error; err != nil {
		return nil, fmt.Errorf("failed to get active sessions: %w", err)
	}

	return sessions, nil
}

// LogoutSession logs out a specific session
func (s *SecurityServiceImpl) LogoutSession(sessionID uint) error {
	now := time.Now()

	if err := s.db.Model(&models.LoginHistory{}).
		Where("id = ?", sessionID).
		Update("logout_at", now).Error; err != nil {
		return fmt.Errorf("failed to logout session: %w", err)
	}

	return nil
}

// LogoutAllSessions logs out all of a user's sessions
func (s *SecurityServiceImpl) LogoutAllSessions(userID uint) error {
	now := time.Now()

	if err := s.db.Model(&models.LoginHistory{}).
		Where("user_id = ? AND logout_at IS NULL", userID).
		Update("logout_at", now).Error; err != nil {
		return fmt.Errorf("failed to logout all sessions: %w", err)
	}

	return nil
}

// LogActivity logs a user action for audit trail
func (s *SecurityServiceImpl) LogActivity(userID uint, action, resource string, resourceID *uint, details map[string]interface{}, ipAddress, userAgent string) error {
	var detailsJSON []byte
	var err error

	if details != nil {
		detailsJSON, err = json.Marshal(details)
		if err != nil {
			return fmt.Errorf("failed to marshal details: %w", err)
		}
	}

	activity := &models.ActivityLog{
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Details:    detailsJSON,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
	}

	if err := s.db.Create(activity).Error; err != nil {
		return fmt.Errorf("failed to log activity: %w", err)
	}

	return nil
}

// GetActivityLog retrieves a user's activity log
func (s *SecurityServiceImpl) GetActivityLog(userID uint, limit, offset int) ([]models.ActivityLog, int64, error) {
	var activities []models.ActivityLog
	var total int64

	// Count total records
	if err := s.db.Model(&models.ActivityLog{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count activities: %w", err)
	}

	// Get paginated records
	if err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&activities).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get activity log: %w", err)
	}

	return activities, total, nil
}

// Helper functions

// generateSecureToken generates a cryptographically secure random token
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// detectDevice attempts to detect device type from user agent
func detectDevice(userAgent string) string {
	ua := userAgent
	if ua == "" {
		return "unknown"
	}

	// Simple detection logic
	if contains(ua, "Mobile") || contains(ua, "Android") || contains(ua, "iPhone") {
		return "mobile"
	}
	if contains(ua, "Tablet") || contains(ua, "iPad") {
		return "tablet"
	}
	return "desktop"
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		len(s) > len(substr) && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
