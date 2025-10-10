package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type TemplateVariable struct {
	Name     string   `json:"name"`
	Label    string   `json:"label"`
	Type     string   `json:"type"`     // text, date, select
	Default  string   `json:"default"`
	Required bool     `json:"required"`
	Options  []string `json:"options,omitempty"`
}

type TemplateVariables []TemplateVariable

// Scan implements sql.Scanner interface
func (tv *TemplateVariables) Scan(value interface{}) error {
	if value == nil {
		*tv = TemplateVariables{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, tv)
}

// Value implements driver.Valuer interface
func (tv TemplateVariables) Value() (driver.Value, error) {
	if len(tv) == 0 {
		return nil, nil
	}
	return json.Marshal(tv)
}

type Template struct {
	ID          uint              `gorm:"primaryKey" json:"id"`
	Name        string            `json:"name"`
	FileInfo    string            `json:"file_info"`
	CategoryID  uint              `json:"category_id"`
	Category    Category          `gorm:"foreignKey:CategoryID"`
	Price       float64           `json:"price"`
	PreviewData string            `json:"preview_data"`
	Variables   TemplateVariables `gorm:"type:jsonb" json:"variables,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	DeletedAt   gorm.DeletedAt    `gorm:"index" json:"-"`
}
