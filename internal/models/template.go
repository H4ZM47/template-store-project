package models

import (
	"time"

	"gorm.io/gorm"
)

type Template struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `json:"name"`
	FileInfo    string         `json:"file_info"`
	CategoryID  uint           `json:"category_id"`
	Category    Category       `gorm:"foreignKey:CategoryID"`
	Price       float64        `json:"price"`
	PreviewData string         `json:"preview_data"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
