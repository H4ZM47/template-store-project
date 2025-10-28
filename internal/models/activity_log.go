package models

import (
	"time"

	"gorm.io/datatypes"
)

// ActivityLog tracks user and admin actions for audit trail
type ActivityLog struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"index;not null" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	// Action details
	Action     string         `gorm:"index" json:"action"` // profile_updated, password_changed, order_placed, etc.
	Resource   string         `json:"resource"`            // user, order, template, blog_post, etc.
	ResourceID *uint          `json:"resource_id,omitempty"`
	Details    datatypes.JSON `json:"details,omitempty"` // Additional context as JSON

	// Request metadata
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`

	CreatedAt time.Time `json:"created_at"`
}

// Common activity action constants
const (
	// Profile actions
	ActivityProfileUpdated = "profile_updated"
	ActivityAvatarUploaded = "avatar_uploaded"
	ActivityAvatarDeleted  = "avatar_deleted"

	// Security actions
	ActivityPasswordChanged      = "password_changed"
	ActivityPasswordResetRequest = "password_reset_requested"
	ActivityPasswordReset        = "password_reset"
	ActivityEmailVerified        = "email_verified"
	ActivityEmailChangeRequested = "email_change_requested"
	ActivityEmailChanged         = "email_changed"

	// Account actions
	ActivityAccountCreated     = "account_created"
	ActivityAccountDeactivated = "account_deactivated"
	ActivityAccountDeleted     = "account_deleted"
	ActivityPreferencesUpdated = "preferences_updated"

	// Login actions
	ActivityLoginSuccess = "login_success"
	ActivityLoginFailed  = "login_failed"
	ActivityLogout       = "logout"

	// Order actions
	ActivityOrderPlaced   = "order_placed"
	ActivityOrderCanceled = "order_canceled"

	// Admin actions
	ActivityUserSuspended   = "user_suspended"
	ActivityUserUnsuspended = "user_unsuspended"
	ActivityRoleChanged     = "role_changed"
	ActivityUserDeleted     = "user_deleted"
)
