package handlers

import (
	"net/http"
	"strings"
	"template-store/internal/services"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Register user
	response, err := h.authService.Register(req)
	if err != nil {
		// Check if it's a duplicate email error
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Login user
	response, err := h.authService.Login(req)
	if err != nil {
		// Don't reveal specific error for security
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProfile returns the current user's profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get user from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	email, _ := c.Get("user_email")
	name, _ := c.Get("user_name")
	role, _ := c.Get("user_role")

	c.JSON(http.StatusOK, gin.H{
		"id":    userID,
		"email": email,
		"name":  name,
		"role":  role,
	})
}

// RefreshToken generates a new token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Get token from header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
		return
	}

	// Extract token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
		return
	}

	// Refresh token
	newToken, err := h.authService.RefreshToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": newToken,
	})
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ChangePassword handles password change
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Change password
	if err := h.authService.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword); err != nil {
		if strings.Contains(err.Error(), "invalid old password") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid old password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}

// ResetPasswordRequestBody represents a password reset request
type ResetPasswordRequestBody struct {
	Email string `json:"email" binding:"required,email"`
}

// RequestPasswordReset initiates password reset process
func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var req ResetPasswordRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Request password reset
	if err := h.authService.ResetPasswordRequest(req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	// Always return success for security (don't reveal if email exists)
	c.JSON(http.StatusOK, gin.H{
		"message": "If the email exists, a password reset link has been sent",
	})
}

// ResetPasswordBody represents a password reset body
type ResetPasswordBody struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ResetPassword resets password using reset token
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Reset password
	if err := h.authService.ResetPassword(req.Token, req.NewPassword); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired reset token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully",
	})
}

// Logout handles user logout (client-side token removal, server can't revoke JWT)
func (h *AuthHandler) Logout(c *gin.Context) {
	// JWT tokens are stateless, so logout is handled client-side
	// This endpoint exists for consistency and future session management
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}
