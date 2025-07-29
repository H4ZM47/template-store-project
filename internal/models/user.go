package models

import (
	"time"
)

type User struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	// Relationships
	BlogPosts []BlogPost `gorm:"foreignKey:AuthorID"`
	Orders    []Order    `gorm:"foreignKey:UserID"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
