package services

import (
	"errors"
	"template-store/internal/models"
	"gorm.io/gorm"
)

type CategoryService struct {
	db *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{db: db}
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(category *models.Category) error {
	if category.Name == "" {
		return errors.New("category name is required")
	}
	
	return s.db.Create(category).Error
}

// GetCategory retrieves a category by ID
func (s *CategoryService) GetCategory(id uint) (*models.Category, error) {
	var category models.Category
	err := s.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// ListCategories retrieves all categories
func (s *CategoryService) ListCategories() ([]models.Category, error) {
	var categories []models.Category
	err := s.db.Find(&categories).Error
	return categories, err
}

// UpdateCategory updates an existing category
func (s *CategoryService) UpdateCategory(id uint, updates map[string]interface{}) error {
	if name, exists := updates["name"]; exists && name == "" {
		return errors.New("category name cannot be empty")
	}
	
	return s.db.Model(&models.Category{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteCategory deletes a category by ID
func (s *CategoryService) DeleteCategory(id uint) error {
	return s.db.Delete(&models.Category{}, id).Error
}

// SeedCategories seeds the database with initial categories
func (s *CategoryService) SeedCategories() error {
	categories := []models.Category{
		{Name: "Security & Compliance"},
	}

	for _, category := range categories {
		if err := s.db.Create(&category).Error; err != nil {
			return err
		}
	}

	return nil
} 