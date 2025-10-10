package models

import (
	"time"
)

type User struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Name        string     `json:"name"`
	Email       string     `gorm:"uniqueIndex;not null" json:"email"`
	Password    string     `gorm:"not null" json:"-"` // Never expose password in JSON
	Role        string     `gorm:"default:'customer'" json:"role"` // customer, admin
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	// Relationships
	BlogPosts []BlogPost `gorm:"foreignKey:AuthorID" json:"blog_posts,omitempty"`
	Orders    []Order    `gorm:"foreignKey:UserID" json:"orders,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
