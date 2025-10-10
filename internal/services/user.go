package services

import (
	"errors"
	"template-store/internal/models"
	"time"

	"gorm.io/gorm"
)

// UserService defines the interface for user-related operations.
type UserService interface {
	CreateUser(user *models.User) error
	GetUser(id uint) (*models.User, error)
	GetUserByCognitoSub(sub string) (*models.User, error)
	ListUsers() ([]models.User, error)
	SeedUsers() error
}

// userServiceImpl is the concrete implementation of the UserService interface.
type userServiceImpl struct {
	db *gorm.DB
}

// NewUserService creates a new UserService.
func NewUserService(db *gorm.DB) UserService {
	return &userServiceImpl{db: db}
}

// CreateUser creates a new user
func (s *userServiceImpl) CreateUser(user *models.User) error {
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
func (s *userServiceImpl) GetUser(id uint) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByCognitoSub retrieves a user by their Cognito subject (sub) claim.
func (s *userServiceImpl) GetUserByCognitoSub(sub string) (*models.User, error) {
	var user models.User
	err := s.db.Where("cognito_subject = ?", sub).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ListUsers retrieves all users
func (s *userServiceImpl) ListUsers() ([]models.User, error) {
	var users []models.User
	err := s.db.Find(&users).Error
	return users, err
}

// SeedUsers seeds the database with initial users
func (s *userServiceImpl) SeedUsers() error {
	// This function needs to be updated to include CognitoSubject
	users := []models.User{
		{Name: "John Doe", Email: "john@example.com", CognitoSubject: "sub_john"},
		{Name: "Jane Smith", Email: "jane@example.com", CognitoSubject: "sub_jane"},
		{Name: "Admin User", Email: "admin@example.com", CognitoSubject: "sub_admin"},
	}

	for _, user := range users {
		// Use a transaction to avoid partial seeding
		tx := s.db.Begin()
		var existingUser models.User
		if err := tx.Where("email = ?", user.Email).First(&existingUser).Error; err == gorm.ErrRecordNotFound {
			if err := tx.Create(&user).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		tx.Commit()
	}

	return nil
}
