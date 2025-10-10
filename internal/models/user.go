package models

import (
	"time"
)

type User struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	CognitoSubject string `gorm:"unique;not null" json:"-"` // Cognito's sub claim
	Name           string `json:"name"`
	Email          string `gorm:"unique" json:"email"`
	// Relationships
	BlogPosts []BlogPost `gorm:"foreignKey:AuthorID"`
	Orders    []Order    `gorm:"foreignKey:UserID"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
