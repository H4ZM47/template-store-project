package models

import "gorm.io/gorm"

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// Core models
		&User{},
		&BlogPost{},
		&Template{},
		&Order{},
		&Category{},

		// User account management models
		&LoginHistory{},
		&ActivityLog{},
		&PasswordResetToken{},
		&EmailVerificationToken{},
	)
} 