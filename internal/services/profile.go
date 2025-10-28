package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
	"template-store/internal/models"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidProfileData    = errors.New("invalid profile data")
	ErrAvatarUploadFailed    = errors.New("avatar upload failed")
	ErrInvalidAvatarFormat   = errors.New("invalid avatar format")
	ErrAvatarTooLarge        = errors.New("avatar file too large")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrInvalidPassword       = errors.New("invalid password")
)

const (
	MaxAvatarSizeMB        = 5
	MaxAvatarSizeBytes     = MaxAvatarSizeMB * 1024 * 1024
	AllowedAvatarFormats   = ".jpg,.jpeg,.png,.webp"
)

// ProfileService defines the interface for user profile operations
type ProfileService interface {
	// Profile operations
	GetProfile(userID uint) (*models.User, error)
	GetPublicProfile(userID uint, viewerID *uint) (*models.User, error)
	UpdateProfile(userID uint, updates map[string]interface{}) (*models.User, error)

	// Avatar operations
	UploadAvatar(ctx context.Context, userID uint, file *multipart.FileHeader) (string, error)
	DeleteAvatar(userID uint) error

	// Preferences
	GetPreferences(userID uint) (*models.UserPreferences, error)
	UpdatePreferences(userID uint, prefs *models.UserPreferences) error

	// Account operations
	RequestEmailChange(userID uint, newEmail, password string) error
	DeactivateAccount(userID uint, reason string) error
	DeleteAccount(userID uint) error
}

// ProfileServiceImpl implements the ProfileService interface
type ProfileServiceImpl struct {
	db             *gorm.DB
	storageService StorageService
	authService    AuthService
}

// NewProfileService creates a new ProfileService instance
func NewProfileService(db *gorm.DB, storageService StorageService, authService AuthService) ProfileService {
	return &ProfileServiceImpl{
		db:             db,
		storageService: storageService,
		authService:    authService,
	}
}

// GetProfile retrieves a user's full profile
func (s *ProfileServiceImpl) GetProfile(userID uint) (*models.User, error) {
	var user models.User

	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	return &user, nil
}

// GetPublicProfile retrieves a user's public profile with privacy settings applied
func (s *ProfileServiceImpl) GetPublicProfile(userID uint, viewerID *uint) (*models.User, error) {
	var user models.User

	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	// If viewer is the same as user, return full profile
	if viewerID != nil && *viewerID == userID {
		return &user, nil
	}

	// Parse preferences to check privacy settings
	prefs, err := s.parsePreferences(&user)
	if err != nil {
		// If preferences can't be parsed, default to restrictive privacy
		prefs = &models.UserPreferences{
			ProfileVisibility: "private",
			ShowEmail:         false,
		}
	}

	// Apply privacy settings
	if prefs.ProfileVisibility == "private" {
		// Return minimal public info
		return &models.User{
			ID:        user.ID,
			Name:      user.Name,
			AvatarURL: user.AvatarURL,
		}, nil
	}

	// For public profiles, hide sensitive information
	publicUser := user
	if !prefs.ShowEmail {
		publicUser.Email = ""
	}

	// Always hide sensitive fields
	publicUser.CognitoSubject = ""
	publicUser.PhoneNumber = ""
	publicUser.Preferences = nil
	publicUser.LastLoginAt = nil
	publicUser.SuspendedAt = nil
	publicUser.SuspendedBy = nil
	publicUser.SuspensionReason = ""

	return &publicUser, nil
}

// UpdateProfile updates a user's profile information
func (s *ProfileServiceImpl) UpdateProfile(userID uint, updates map[string]interface{}) (*models.User, error) {
	var user models.User

	// First check if user exists
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Validate that restricted fields are not being updated
	restrictedFields := []string{
		"id", "cognito_subject", "role", "status", "email_verified",
		"suspended_at", "suspended_by", "suspension_reason", "created_at",
	}

	for _, field := range restrictedFields {
		if _, exists := updates[field]; exists {
			return nil, fmt.Errorf("%w: cannot update %s", ErrInvalidProfileData, field)
		}
	}

	// Perform update
	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// Reload user with updated data
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("failed to reload profile: %w", err)
	}

	return &user, nil
}

// UploadAvatar uploads and sets a user's avatar
func (s *ProfileServiceImpl) UploadAvatar(ctx context.Context, userID uint, file *multipart.FileHeader) (string, error) {
	// Validate file size
	if file.Size > MaxAvatarSizeBytes {
		return "", ErrAvatarTooLarge
	}

	// Validate file format
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !strings.Contains(AllowedAvatarFormats, ext) {
		return "", ErrInvalidAvatarFormat
	}

	// Check if user exists
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrUserNotFound
		}
		return "", fmt.Errorf("failed to find user: %w", err)
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	newFilename := fmt.Sprintf("avatars/user_%d_%d%s", userID, timestamp, ext)

	// Create a new file header with the new filename
	avatarFile := *file
	avatarFile.Filename = newFilename

	// Upload to storage
	avatarURL, err := s.storageService.UploadFile(ctx, &avatarFile)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrAvatarUploadFailed, err)
	}

	// Delete old avatar if exists
	if user.AvatarURL != "" {
		// Non-blocking delete, log error but don't fail
		_ = s.deleteAvatarFromStorage(user.AvatarURL)
	}

	// Update user's avatar URL in database
	if err := s.db.Model(&user).Update("avatar_url", avatarURL).Error; err != nil {
		return "", fmt.Errorf("failed to update avatar URL: %w", err)
	}

	return avatarURL, nil
}

// DeleteAvatar removes a user's avatar
func (s *ProfileServiceImpl) DeleteAvatar(userID uint) error {
	var user models.User

	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	// If no avatar, nothing to delete
	if user.AvatarURL == "" {
		return nil
	}

	// Delete from storage
	if err := s.deleteAvatarFromStorage(user.AvatarURL); err != nil {
		// Log error but continue to clear DB
		fmt.Printf("Warning: failed to delete avatar from storage: %v\n", err)
	}

	// Clear avatar URL in database
	if err := s.db.Model(&user).Update("avatar_url", "").Error; err != nil {
		return fmt.Errorf("failed to clear avatar URL: %w", err)
	}

	return nil
}

// GetPreferences retrieves a user's preferences
func (s *ProfileServiceImpl) GetPreferences(userID uint) (*models.UserPreferences, error) {
	var user models.User

	if err := s.db.Select("preferences").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get preferences: %w", err)
	}

	return s.parsePreferences(&user)
}

// UpdatePreferences updates a user's preferences
func (s *ProfileServiceImpl) UpdatePreferences(userID uint, prefs *models.UserPreferences) error {
	// Serialize preferences to JSON
	prefsJSON, err := json.Marshal(prefs)
	if err != nil {
		return fmt.Errorf("failed to marshal preferences: %w", err)
	}

	// Update in database
	if err := s.db.Model(&models.User{}).Where("id = ?", userID).Update("preferences", prefsJSON).Error; err != nil {
		return fmt.Errorf("failed to update preferences: %w", err)
	}

	return nil
}

// RequestEmailChange initiates an email change request
func (s *ProfileServiceImpl) RequestEmailChange(userID uint, newEmail, password string) error {
	var user models.User

	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Check if new email is already in use
	var existingUser models.User
	if err := s.db.Where("email = ? AND id != ?", newEmail, userID).First(&existingUser).Error; err == nil {
		return ErrEmailAlreadyExists
	}

	// Verify password with Cognito
	// Note: This would require Cognito integration
	// For now, we'll create the email verification token

	// TODO: Create email verification token and send verification email
	// This will be implemented when we create the SecurityService and EmailService

	return nil
}

// DeactivateAccount deactivates a user's account
func (s *ProfileServiceImpl) DeactivateAccount(userID uint, reason string) error {
	updates := map[string]interface{}{
		"status":            "inactive",
		"suspension_reason": reason,
		"suspended_at":      time.Now(),
	}

	if err := s.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to deactivate account: %w", err)
	}

	return nil
}

// DeleteAccount soft deletes a user's account
func (s *ProfileServiceImpl) DeleteAccount(userID uint) error {
	var user models.User

	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Soft delete (sets DeletedAt timestamp)
	if err := s.db.Delete(&user).Error; err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	return nil
}

// Helper functions

func (s *ProfileServiceImpl) parsePreferences(user *models.User) (*models.UserPreferences, error) {
	if user.Preferences == nil || len(user.Preferences) == 0 {
		// Return default preferences
		defaultPrefs := models.DefaultPreferences()
		return &defaultPrefs, nil
	}

	var prefs models.UserPreferences
	if err := json.Unmarshal(user.Preferences, &prefs); err != nil {
		return nil, fmt.Errorf("failed to parse preferences: %w", err)
	}

	return &prefs, nil
}

func (s *ProfileServiceImpl) deleteAvatarFromStorage(avatarURL string) error {
	if avatarURL == "" {
		return nil
	}

	ctx := context.Background()
	if err := s.storageService.DeleteFile(ctx, avatarURL); err != nil {
		return fmt.Errorf("failed to delete avatar from storage: %w", err)
	}

	return nil
}
