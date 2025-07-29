package services

import (
	"errors"
	"template-store/internal/models"
	"time"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(user *models.User) error {
	if user.Name == "" {
		return errors.New("user name is required")
	}
	if user.Email == "" {
		return errors.New("user email is required")
	}

	// Set default values
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = time.Now()
	}

	return s.db.Create(user).Error
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id uint) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ListUsers retrieves all users
func (s *UserService) ListUsers() ([]models.User, error) {
	var users []models.User
	err := s.db.Find(&users).Error
	return users, err
}

// SeedUsers seeds the database with initial users
func (s *UserService) SeedUsers() error {
	users := []models.User{
		{Name: "John Doe", Email: "john@example.com"},
		{Name: "Jane Smith", Email: "jane@example.com"},
		{Name: "Admin User", Email: "admin@example.com"},
	}

	for _, user := range users {
		if err := s.db.Create(&user).Error; err != nil {
			return err
		}
	}

	return nil
}
