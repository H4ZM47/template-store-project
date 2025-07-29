package models

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	UserID         uint           `json:"user_id"`
	User           User           `gorm:"foreignKey:UserID"`
	TemplateID     uint           `json:"template_id"`
	Template       Template       `gorm:"foreignKey:TemplateID"`
	PurchaseHistory string        `json:"purchase_history"`
	DeliveryStatus string         `json:"delivery_status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
} 