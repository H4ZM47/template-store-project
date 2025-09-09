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
		{Name: "Modern Corporate Website", CategoryID: 1, Price: 49, FileInfo: "Responsive HTML5/CSS3 template for corporate sites."},
		{Name: "Creative Portfolio", CategoryID: 8, Price: 29, FileInfo: "A stylish portfolio template for designers and photographers."},
		{Name: "E-commerce Storefront", CategoryID: 7, Price: 99, FileInfo: "Feature-rich template for online stores."},
		{Name: "Minimalist Blog", CategoryID: 1, Price: 19, FileInfo: "A clean and simple blog template."},
		{Name: "Product Launch Page", CategoryID: 6, Price: 25, FileInfo: "A one-page template for launching new products."},
		{Name: "Real Estate Listing", CategoryID: 1, Price: 39, FileInfo: "Template for real estate agencies and listings."},
		{Name: "SaaS Landing Page", CategoryID: 6, Price: 35, FileInfo: "High-converting landing page for SaaS products."},
		{Name: "Restaurant Website", CategoryID: 1, Price: 39, FileInfo: "Template for restaurants, cafes, and bars."},
		{Name: "Mobile App Showcase", CategoryID: 6, Price: 29, FileInfo: "Showcase your mobile app with this modern template."},
		{Name: "Startup Pitch Deck", CategoryID: 5, Price: 15, FileInfo: "Professional presentation template for startups."},
		{Name: "Weekly Newsletter", CategoryID: 2, Price: 9, FileInfo: "A responsive email template for newsletters."},
		{Name: "Flash Sale Promotion", CategoryID: 2, Price: 5, FileInfo: "Email template for flash sales and promotions."},
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
