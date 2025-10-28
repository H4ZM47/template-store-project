package models

import "time"

// LoginHistory tracks user login attempts and sessions
type LoginHistory struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"index;not null" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	// Request details
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	Device    string `json:"device"` // mobile, desktop, tablet, unknown
	Location  string `json:"location,omitempty"` // City, Country (from IP geolocation)

	// Session tracking
	LoginAt  time.Time  `json:"login_at"`
	LogoutAt *time.Time `json:"logout_at,omitempty"`

	// Login status
	Success       bool   `json:"success"`
	FailureReason string `json:"failure_reason,omitempty"` // invalid_credentials, account_suspended, etc.

	CreatedAt time.Time `json:"created_at"`
}

// IsActiveSession returns true if the session is still active (not logged out)
func (lh *LoginHistory) IsActiveSession() bool {
	return lh.Success && lh.LogoutAt == nil
}
