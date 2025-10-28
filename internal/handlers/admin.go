package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"template-store/internal/services"
)

// AdminHandler handles admin-related HTTP requests
type AdminHandler struct {
	adminService services.AdminService
}

// NewAdminHandler creates a new AdminHandler
func NewAdminHandler(adminService services.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

// ListUsers handles GET /api/v1/admin/users
func (h *AdminHandler) ListUsers(c *gin.Context) {
	// Parse query params
	filters := services.UserFilters{
		Search: c.Query("search"),
		Role:   c.Query("role"),
		Status: c.Query("status"),
	}

	sortBy := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	users, total, err := h.adminService.ListUsers(filters, sortBy, order, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":  users,
		"total":  total,
		"limit":  limit,
		"offset": offset,
		"page":   (offset / limit) + 1,
	})
}

// GetUser handles GET /api/v1/admin/users/:id
func (h *AdminHandler) GetUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.adminService.GetUserDetails(uint(userID))
	if err != nil {
		if err == services.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	// Get additional stats
	stats, err := h.adminService.GetUserStats(uint(userID))
	if err != nil {
		stats = map[string]interface{}{}
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"stats": stats,
	})
}

// UpdateUserRole handles PUT /api/v1/admin/users/:id/role
func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	adminID := getUserIDFromContext(c)

	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role is required"})
		return
	}

	if err := h.adminService.UpdateUserRole(adminID, uint(userID), req.Role); err != nil {
		switch err {
		case services.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		case services.ErrInvalidRole:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		case services.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role updated successfully"})
}

// SuspendUser handles POST /api/v1/admin/users/:id/suspend
func (h *AdminHandler) SuspendUser(c *gin.Context) {
	adminID := getUserIDFromContext(c)

	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Reason       string `json:"reason" binding:"required"`
		DurationDays int    `json:"duration_days"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reason is required"})
		return
	}

	if err := h.adminService.SuspendUser(adminID, uint(userID), req.Reason, req.DurationDays); err != nil {
		switch err {
		case services.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		case services.ErrCannotSuspendSelf:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot suspend your own account"})
		case services.ErrUserAlreadySuspended:
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is already suspended"})
		case services.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to suspend user"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User suspended successfully"})
}

// UnsuspendUser handles POST /api/v1/admin/users/:id/unsuspend
func (h *AdminHandler) UnsuspendUser(c *gin.Context) {
	adminID := getUserIDFromContext(c)

	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.adminService.UnsuspendUser(adminID, uint(userID)); err != nil {
		switch err {
		case services.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		case services.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsuspend user"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User unsuspended successfully"})
}

// DeleteUser handles DELETE /api/v1/admin/users/:id
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	adminID := getUserIDFromContext(c)

	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.adminService.DeleteUser(adminID, uint(userID)); err != nil {
		switch err {
		case services.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		case services.ErrCannotDeleteSelf:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete your own account"})
		case services.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetDashboard handles GET /api/v1/admin/dashboard
func (h *AdminHandler) GetDashboard(c *gin.Context) {
	stats, additional, err := h.adminService.GetAdminDashboard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dashboard data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats":         stats,
		"recent_users":  additional["recent_users"],
		"recent_orders": additional["recent_orders"],
	})
}

// GetUserActivity handles GET /api/v1/admin/users/:id/activity
func (h *AdminHandler) GetUserActivity(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Parse pagination params
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	activities, total, err := h.adminService.GetUserActivity(uint(userID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user activity"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"activities": activities,
		"total":      total,
		"limit":      limit,
		"offset":     offset,
	})
}
