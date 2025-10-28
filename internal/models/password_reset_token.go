package models

import "time"

// PasswordResetToken stores tokens for password reset requests
type PasswordResetToken struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"index;not null" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	Token     string    `gorm:"uniqueIndex;not null" json:"-"` // Don't expose in JSON
	ExpiresAt time.Time `gorm:"index;not null" json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

// IsExpired returns true if the token has expired
func (prt *PasswordResetToken) IsExpired() bool {
	return time.Now().After(prt.ExpiresAt)
}

// IsUsed returns true if the token has already been used
func (prt *PasswordResetToken) IsUsed() bool {
	return prt.UsedAt != nil
}

// IsValid returns true if the token is valid (not expired and not used)
func (prt *PasswordResetToken) IsValid() bool {
	return !prt.IsExpired() && !prt.IsUsed()
}
