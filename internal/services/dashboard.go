package services

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"template-store/internal/models"
)

// DashboardStats represents dashboard statistics
type DashboardStats struct {
	TotalOrders        int64   `json:"total_orders"`
	TotalSpent         float64 `json:"total_spent"`
	TemplatesPurchased int64   `json:"templates_purchased"`
	BlogPostsAuthored  int64   `json:"blog_posts_authored"`
	AccountAgeDays     int     `json:"account_age_days"`
}

// OrderWithTemplate represents an order with template details
type OrderWithTemplate struct {
	models.Order
	Template models.Template `json:"template"`
}

// PurchasedTemplate represents a template purchase with download info
type PurchasedTemplate struct {
	Template      models.Template `json:"template"`
	PurchasedAt   time.Time       `json:"purchased_at"`
	OrderID       uint            `json:"order_id"`
	DownloadURL   string          `json:"download_url"`
	DownloadCount int             `json:"download_count"`
}

// DashboardService defines the interface for user dashboard operations
type DashboardService interface {
	// Dashboard summary
	GetDashboardStats(userID uint) (*DashboardStats, error)
	GetDashboard(userID uint) (map[string]interface{}, error)

	// Order management
	GetOrderHistory(userID uint, status string, limit, offset int) ([]OrderWithTemplate, int64, error)
	GetOrder(userID uint, orderID uint) (*OrderWithTemplate, error)

	// Template access
	GetPurchasedTemplates(userID uint) ([]PurchasedTemplate, error)
	GetTemplateDownloadURL(userID uint, templateID uint) (string, error)

	// Blog posts
	GetAuthoredBlogPosts(userID uint, status string, limit, offset int) ([]models.BlogPost, int64, error)
}

// DashboardServiceImpl implements the DashboardService interface
type DashboardServiceImpl struct {
	db *gorm.DB
}

// NewDashboardService creates a new DashboardService instance
func NewDashboardService(db *gorm.DB) DashboardService {
	return &DashboardServiceImpl{
		db: db,
	}
}

// GetDashboardStats retrieves dashboard statistics for a user
func (s *DashboardServiceImpl) GetDashboardStats(userID uint) (*DashboardStats, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	stats := &DashboardStats{}

	// Count total orders
	if err := s.db.Model(&models.Order{}).Where("user_id = ?", userID).Count(&stats.TotalOrders).Error; err != nil {
		return nil, fmt.Errorf("failed to count orders: %w", err)
	}

	// Calculate total spent (assuming Order has an Amount field - adjust based on actual schema)
	// For now, we'll leave it at 0 as the Order model structure needs to be checked
	// TODO: Implement when Order model has price/amount field
	stats.TotalSpent = 0

	// Count unique templates purchased
	if err := s.db.Model(&models.Order{}).
		Where("user_id = ?", userID).
		Distinct("template_id").
		Count(&stats.TemplatesPurchased).Error; err != nil {
		return nil, fmt.Errorf("failed to count purchased templates: %w", err)
	}

	// Count blog posts authored
	if err := s.db.Model(&models.BlogPost{}).Where("author_id = ?", userID).Count(&stats.BlogPostsAuthored).Error; err != nil {
		return nil, fmt.Errorf("failed to count blog posts: %w", err)
	}

	// Calculate account age in days
	accountAge := time.Since(user.CreatedAt)
	stats.AccountAgeDays = int(accountAge.Hours() / 24)

	return stats, nil
}

// GetDashboard retrieves complete dashboard data
func (s *DashboardServiceImpl) GetDashboard(userID uint) (map[string]interface{}, error) {
	stats, err := s.GetDashboardStats(userID)
	if err != nil {
		return nil, err
	}

	// Get recent orders (last 5)
	recentOrders, _, err := s.GetOrderHistory(userID, "", 5, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent orders: %w", err)
	}

	// Get recent blog posts (last 5)
	recentBlogPosts, _, err := s.GetAuthoredBlogPosts(userID, "published", 5, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent blog posts: %w", err)
	}

	dashboard := map[string]interface{}{
		"stats":              stats,
		"recent_orders":      recentOrders,
		"recent_blog_posts":  recentBlogPosts,
	}

	return dashboard, nil
}

// GetOrderHistory retrieves a user's order history
func (s *DashboardServiceImpl) GetOrderHistory(userID uint, status string, limit, offset int) ([]OrderWithTemplate, int64, error) {
	var orders []models.Order
	var total int64

	query := s.db.Model(&models.Order{}).Where("user_id = ?", userID)

	// Filter by status if provided
	if status != "" {
		query = query.Where("delivery_status = ?", status)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	// Get paginated records with template preloaded
	if err := query.
		Preload("Template").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get orders: %w", err)
	}

	// Convert to OrderWithTemplate
	ordersWithTemplate := make([]OrderWithTemplate, len(orders))
	for i, order := range orders {
		ordersWithTemplate[i] = OrderWithTemplate{
			Order:    order,
			Template: order.Template,
		}
	}

	return ordersWithTemplate, total, nil
}

// GetOrder retrieves a specific order
func (s *DashboardServiceImpl) GetOrder(userID uint, orderID uint) (*OrderWithTemplate, error) {
	var order models.Order

	if err := s.db.Preload("Template").
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	orderWithTemplate := &OrderWithTemplate{
		Order:    order,
		Template: order.Template,
	}

	return orderWithTemplate, nil
}

// GetPurchasedTemplates retrieves all templates purchased by a user
func (s *DashboardServiceImpl) GetPurchasedTemplates(userID uint) ([]PurchasedTemplate, error) {
	var orders []models.Order

	// Get all orders with templates
	if err := s.db.Preload("Template").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to get purchased templates: %w", err)
	}

	// Convert to PurchasedTemplate
	purchasedTemplates := make([]PurchasedTemplate, 0)
	seenTemplates := make(map[uint]bool) // Track unique templates

	for _, order := range orders {
		// Skip duplicates (user may have purchased same template multiple times)
		if seenTemplates[order.TemplateID] {
			continue
		}
		seenTemplates[order.TemplateID] = true

		purchased := PurchasedTemplate{
			Template:      order.Template,
			PurchasedAt:   order.CreatedAt,
			OrderID:       order.ID,
			DownloadURL:   "", // TODO: Generate download URL
			DownloadCount: 0,  // TODO: Track download count
		}
		purchasedTemplates = append(purchasedTemplates, purchased)
	}

	return purchasedTemplates, nil
}

// GetTemplateDownloadURL generates a download URL for a purchased template
func (s *DashboardServiceImpl) GetTemplateDownloadURL(userID uint, templateID uint) (string, error) {
	// Verify user has purchased this template
	var order models.Order
	if err := s.db.Where("user_id = ? AND template_id = ?", userID, templateID).
		First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("template not purchased")
		}
		return "", fmt.Errorf("failed to verify purchase: %w", err)
	}

	// Get template details
	var template models.Template
	if err := s.db.First(&template, templateID).Error; err != nil {
		return "", fmt.Errorf("failed to get template: %w", err)
	}

	// TODO: Generate signed/temporary download URL from S3
	// For now, return the template's file info (assuming it contains URL)
	downloadURL := template.FileInfo // Adjust based on actual Template schema

	return downloadURL, nil
}

// GetAuthoredBlogPosts retrieves blog posts authored by a user
func (s *DashboardServiceImpl) GetAuthoredBlogPosts(userID uint, status string, limit, offset int) ([]models.BlogPost, int64, error) {
	var posts []models.BlogPost
	var total int64

	query := s.db.Model(&models.BlogPost{}).Where("author_id = ?", userID)

	// Filter by status if provided (e.g., "published", "draft")
	if status != "" {
		// TODO: Add status field to BlogPost model if needed
		// query = query.Where("status = ?", status)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count blog posts: %w", err)
	}

	// Get paginated records
	if err := query.
		Preload("Category").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get blog posts: %w", err)
	}

	return posts, total, nil
}
