package models

import "time"

type Category struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `json:"name"`
	// Relationships
	BlogPosts  []BlogPost  `gorm:"foreignKey:CategoryID"`
	Templates  []Template  `gorm:"foreignKey:CategoryID"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
} 