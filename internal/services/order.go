package services

import (
	"fmt"
	"template-store/internal/models"
	"time"

	"gorm.io/gorm"
)

// OrderService handles order-related business logic
type OrderService struct {
	db *gorm.DB
}

// NewOrderService creates a new order service
func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusFailed    OrderStatus = "failed"
	OrderStatusRefunded  OrderStatus = "refunded"
)

// CreateOrder creates a new order in the database
func (s *OrderService) CreateOrder(userID, templateID uint, stripeSessionID string) (*models.Order, error) {
	order := &models.Order{
		UserID:          userID,
		TemplateID:      templateID,
		PurchaseHistory: fmt.Sprintf("Stripe Session: %s", stripeSessionID),
		DeliveryStatus:  string(OrderStatusPending),
	}

	if err := s.db.Create(order).Error; err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Preload relationships
	if err := s.db.Preload("User").Preload("Template").First(order, order.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load order relationships: %w", err)
	}

	return order, nil
}

// GetOrderByID retrieves an order by ID
func (s *OrderService) GetOrderByID(id uint) (*models.Order, error) {
	var order models.Order
	if err := s.db.Preload("User").Preload("Template").First(&order, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return &order, nil
}

// GetOrdersByUserID retrieves all orders for a specific user
func (s *OrderService) GetOrdersByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	if err := s.db.Preload("User").Preload("Template").Where("user_id = ?", userID).Order("created_at desc").Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}
	return orders, nil
}

// GetAllOrders retrieves all orders with pagination
func (s *OrderService) GetAllOrders(limit, offset int) ([]models.Order, error) {
	var orders []models.Order
	query := s.db.Preload("User").Preload("Template").Order("created_at desc")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	return orders, nil
}

// UpdateOrderStatus updates the delivery status of an order
func (s *OrderService) UpdateOrderStatus(orderID uint, status OrderStatus) error {
	result := s.db.Model(&models.Order{}).Where("id = ?", orderID).Update("delivery_status", string(status))
	if result.Error != nil {
		return fmt.Errorf("failed to update order status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("order not found")
	}
	return nil
}

// MarkOrderAsPaid marks an order as paid and updates purchase history
func (s *OrderService) MarkOrderAsPaid(orderID uint, paymentDetails string) error {
	updates := map[string]interface{}{
		"delivery_status":  string(OrderStatusPaid),
		"purchase_history": paymentDetails,
	}

	result := s.db.Model(&models.Order{}).Where("id = ?", orderID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to mark order as paid: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("order not found")
	}
	return nil
}

// MarkOrderAsDelivered marks an order as delivered
func (s *OrderService) MarkOrderAsDelivered(orderID uint) error {
	return s.UpdateOrderStatus(orderID, OrderStatusDelivered)
}

// MarkOrderAsFailed marks an order as failed
func (s *OrderService) MarkOrderAsFailed(orderID uint, reason string) error {
	updates := map[string]interface{}{
		"delivery_status":  string(OrderStatusFailed),
		"purchase_history": reason,
	}

	result := s.db.Model(&models.Order{}).Where("id = ?", orderID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to mark order as failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("order not found")
	}
	return nil
}

// GetOrderByStripeSession retrieves an order by Stripe session ID
func (s *OrderService) GetOrderByStripeSession(sessionID string) (*models.Order, error) {
	var order models.Order
	searchPattern := fmt.Sprintf("%%Stripe Session: %s%%", sessionID)

	if err := s.db.Preload("User").Preload("Template").Where("purchase_history LIKE ?", searchPattern).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("order not found for session")
		}
		return nil, fmt.Errorf("failed to get order by session: %w", err)
	}
	return &order, nil
}

// OrderResponse represents the JSON response for an order
type OrderResponse struct {
	ID              uint      `json:"id"`
	UserID          uint      `json:"user_id"`
	UserName        string    `json:"user_name"`
	UserEmail       string    `json:"user_email"`
	TemplateID      uint      `json:"template_id"`
	TemplateName    string    `json:"template_name"`
	TemplatePrice   float64   `json:"template_price"`
	PurchaseHistory string    `json:"purchase_history"`
	DeliveryStatus  string    `json:"delivery_status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ToOrderResponse converts an Order model to OrderResponse
func ToOrderResponse(order *models.Order) OrderResponse {
	response := OrderResponse{
		ID:              order.ID,
		UserID:          order.UserID,
		TemplateID:      order.TemplateID,
		PurchaseHistory: order.PurchaseHistory,
		DeliveryStatus:  order.DeliveryStatus,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}

	if order.User.ID != 0 {
		response.UserName = order.User.Name
		response.UserEmail = order.User.Email
	}

	if order.Template.ID != 0 {
		response.TemplateName = order.Template.Name
		response.TemplatePrice = order.Template.Price
	}

	return response
}
