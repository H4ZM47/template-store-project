package services

import (
	"template-store/internal/models"

	"gorm.io/gorm"
)

// OrderService defines the interface for order-related operations.
type OrderService interface {
	CreateOrder(order *models.Order) error
}

// orderServiceImpl is the concrete implementation of the OrderService interface.
type orderServiceImpl struct {
	db *gorm.DB
}

// NewOrderService creates a new OrderService.
func NewOrderService(db *gorm.DB) OrderService {
	return &orderServiceImpl{db: db}
}

// CreateOrder creates a new order in the database.
func (s *orderServiceImpl) CreateOrder(order *models.Order) error {
	return s.db.Create(order).Error
}
