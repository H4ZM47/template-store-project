package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"template-store/internal/services"
)

// DashboardHandler handles dashboard-related HTTP requests
type DashboardHandler struct {
	dashboardService services.DashboardService
}

// NewDashboardHandler creates a new DashboardHandler
func NewDashboardHandler(dashboardService services.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetDashboard handles GET /api/v1/profile/dashboard
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	userID := getUserIDFromContext(c)

	dashboard, err := h.dashboardService.GetDashboard(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dashboard data"})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}

// GetOrders handles GET /api/v1/profile/orders
func (h *DashboardHandler) GetOrders(c *gin.Context) {
	userID := getUserIDFromContext(c)

	// Parse query params
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	status := c.Query("status") // Optional filter

	if limit > 100 {
		limit = 100
	}

	orders, total, err := h.dashboardService.GetOrderHistory(userID, status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
		return
	}

	// Calculate total spent
	var totalSpent float64
	for range orders {
		// TODO: Add price calculation based on order/template price
		// totalSpent += order.Amount
	}

	c.JSON(http.StatusOK, gin.H{
		"orders":      orders,
		"total":       total,
		"limit":       limit,
		"offset":      offset,
		"total_spent": totalSpent,
	})
}

// GetOrder handles GET /api/v1/profile/orders/:id
func (h *DashboardHandler) GetOrder(c *gin.Context) {
	userID := getUserIDFromContext(c)

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.dashboardService.GetOrder(userID, uint(orderID))
	if err != nil {
		if err.Error() == "order not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get order"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetPurchasedTemplates handles GET /api/v1/profile/purchased-templates
func (h *DashboardHandler) GetPurchasedTemplates(c *gin.Context) {
	userID := getUserIDFromContext(c)

	templates, err := h.dashboardService.GetPurchasedTemplates(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get purchased templates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

// GetTemplateDownload handles GET /api/v1/profile/templates/:id/download
func (h *DashboardHandler) GetTemplateDownload(c *gin.Context) {
	userID := getUserIDFromContext(c)

	templateID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	downloadURL, err := h.dashboardService.GetTemplateDownloadURL(userID, uint(templateID))
	if err != nil {
		if err.Error() == "template not purchased" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Template not purchased"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get download URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"download_url": downloadURL})
}

// GetBlogPosts handles GET /api/v1/profile/blog-posts
func (h *DashboardHandler) GetBlogPosts(c *gin.Context) {
	userID := getUserIDFromContext(c)

	// Parse query params
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	status := c.Query("status") // Optional filter (published, draft, etc.)

	if limit > 100 {
		limit = 100
	}

	posts, total, err := h.dashboardService.GetAuthoredBlogPosts(userID, status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get blog posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts":  posts,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}
