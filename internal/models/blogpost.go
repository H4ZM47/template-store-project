package models

import (
	"gorm.io/gorm"
	"time"
)

type BlogPost struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	AuthorID  uint           `json:"author_id"`
	Author    User           `gorm:"foreignKey:AuthorID"`
	CategoryID uint          `json:"category_id"`
	Category  Category       `gorm:"foreignKey:CategoryID"`
	SEO       string         `json:"seo"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
} 