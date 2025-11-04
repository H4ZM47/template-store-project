package models

import "time"

// EmailVerificationToken stores tokens for email verification and email change requests
type EmailVerificationToken struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"index;not null" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	Token      string    `gorm:"uniqueIndex;not null" json:"-"` // Don't expose in JSON
	Email      string    `gorm:"not null" json:"email"`          // Email to verify (for email change requests)
	ExpiresAt  time.Time `gorm:"index;not null" json:"expires_at"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

// IsExpired returns true if the token has expired
func (evt *EmailVerificationToken) IsExpired() bool {
	return time.Now().After(evt.ExpiresAt)
}

// IsVerified returns true if the token has already been verified
func (evt *EmailVerificationToken) IsVerified() bool {
	return evt.VerifiedAt != nil
}

// IsValid returns true if the token is valid (not expired and not verified)
func (evt *EmailVerificationToken) IsValid() bool {
	return !evt.IsExpired() && !evt.IsVerified()
}
