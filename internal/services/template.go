package services

import (
	"errors"
	"template-store/internal/models"

	"gorm.io/gorm"
)

// TemplateService defines the interface for template-related operations.
type TemplateService interface {
	CreateTemplate(template *models.Template) error
	GetTemplate(id uint) (*models.Template, error)
	ListTemplates(categoryID *uint, search *string, limit, offset int) ([]models.Template, int64, error)
	UpdateTemplate(id uint, updates map[string]interface{}) error
	DeleteTemplate(id uint) error
	GetTemplatesByCategory(categoryID uint) ([]models.Template, error)
	SeedTemplates() error
}

// templateServiceImpl is the concrete implementation of the TemplateService interface.
type templateServiceImpl struct {
	db *gorm.DB
}

// NewTemplateService creates a new TemplateService.
func NewTemplateService(db *gorm.DB) TemplateService {
	return &templateServiceImpl{db: db}
}

// CreateTemplate creates a new template
func (s *templateServiceImpl) CreateTemplate(template *models.Template) error {
	if template.Name == "" {
		return errors.New("template name is required")
	}
	if template.Price < 0 {
		return errors.New("template price cannot be negative")
	}

	return s.db.Create(template).Error
}

// GetTemplate retrieves a template by ID
func (s *templateServiceImpl) GetTemplate(id uint) (*models.Template, error) {
	var template models.Template
	err := s.db.Preload("Category").First(&template, id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// ListTemplates retrieves all templates with optional filtering
func (s *templateServiceImpl) ListTemplates(categoryID *uint, search *string, limit, offset int) ([]models.Template, int64, error) {
	var templates []models.Template
	var total int64

	query := s.db.Model(&models.Template{}).Preload("Category")

	// Apply category filter
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	// Apply search filter
	if search != nil && *search != "" {
		query = query.Where("name ILIKE ?", "%"+*search+"%")
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&templates).Error
	return templates, total, err
}

// UpdateTemplate updates an existing template
func (s *templateServiceImpl) UpdateTemplate(id uint, updates map[string]interface{}) error {
	if name, exists := updates["name"]; exists && name == "" {
		return errors.New("template name cannot be empty")
	}
	if price, exists := updates["price"]; exists {
		if p, ok := price.(float64); ok && p < 0 {
			return errors.New("template price cannot be negative")
		}
	}

	return s.db.Model(&models.Template{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteTemplate deletes a template by ID
func (s *templateServiceImpl) DeleteTemplate(id uint) error {
	return s.db.Delete(&models.Template{}, id).Error
}

// GetTemplatesByCategory retrieves templates filtered by category
func (s *templateServiceImpl) GetTemplatesByCategory(categoryID uint) ([]models.Template, error) {
	var templates []models.Template
	err := s.db.Where("category_id = ?", categoryID).Preload("Category").Find(&templates).Error
	return templates, err
}

// SeedTemplates seeds the database with initial templates
func (s *templateServiceImpl) SeedTemplates() error {
	templates := []models.Template{
		{Name: "Data Classification Standard", CategoryID: 1, Price: 49.99, FileInfo: "Comprehensive data classification standard document with guidelines for handling sensitive information."},
		{Name: "Vulnerability Management Standard", CategoryID: 1, Price: 49.99, FileInfo: "Complete vulnerability management standard with assessment procedures and remediation workflows."},
	}

	for _, template := range templates {
		if err := s.db.Create(&template).Error; err != nil {
			// Don't error out if a template already exists, just continue
			if !errors.Is(err, gorm.ErrDuplicatedKey) {
				return err
			}
		}
	}

	return nil
}
