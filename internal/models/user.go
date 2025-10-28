package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	// Core identity fields
	ID             uint   `gorm:"primaryKey" json:"id"`
	CognitoSubject string `gorm:"uniqueIndex;not null" json:"-"` // Cognito's sub claim
	Name           string `json:"name"`
	Email          string `gorm:"uniqueIndex" json:"email"`

	// Profile fields
	AvatarURL   string `json:"avatar_url,omitempty"`
	Bio         string `gorm:"type:text" json:"bio,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`

	// Address fields
	AddressLine1 string `json:"address_line1,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	Country      string `json:"country,omitempty"`
	PostalCode   string `json:"postal_code,omitempty"`

	// Account status
	Role          string `gorm:"default:'user';index" json:"role"`                // user, admin, author
	Status        string `gorm:"default:'active';index" json:"status"`            // active, suspended, deleted
	EmailVerified bool   `gorm:"default:false" json:"email_verified"`

	// User preferences (stored as JSON)
	Preferences datatypes.JSON `json:"preferences,omitempty"`

	// Account metadata
	LastLoginAt      *time.Time `json:"last_login_at,omitempty"`
	SuspendedAt      *time.Time `json:"suspended_at,omitempty"`
	SuspendedBy      *uint      `json:"suspended_by,omitempty"`
	SuspensionReason string     `json:"suspension_reason,omitempty"`

	// Relationships
	BlogPosts []BlogPost `gorm:"foreignKey:AuthorID" json:"blog_posts,omitempty"`
	Orders    []Order    `gorm:"foreignKey:UserID" json:"orders,omitempty"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete support
}
