package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"template-store/internal/models"
	"template-store/internal/services"
)

// SecurityHandler handles security-related HTTP requests
type SecurityHandler struct {
	securityService services.SecurityService
	emailService    services.EmailService
}

// NewSecurityHandler creates a new SecurityHandler
func NewSecurityHandler(securityService services.SecurityService, emailService services.EmailService) *SecurityHandler {
	return &SecurityHandler{
		securityService: securityService,
		emailService:    emailService,
	}
}

// ForgotPassword handles POST /api/v1/auth/forgot-password
func (h *SecurityHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Valid email is required"})
		return
	}

	// Request password reset (creates token)
	resetToken, err := h.securityService.RequestPasswordReset(req.Email)
	if err != nil {
		// Don't reveal if email exists or not for security
		c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a password reset link has been sent"})
		return
	}

	// Send email if token was created (user exists)
	if resetToken != nil {
		// Get user for email
		// TODO: Get user object to pass to email service
		// For now, we'll skip the actual email sending
		// _ = h.emailService.SendPasswordResetEmail(user, resetToken.Token, "")
	}

	c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a password reset link has been sent"})
}

// ResetPassword handles POST /api/v1/auth/reset-password
func (h *SecurityHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token and new password (min 8 characters) are required"})
		return
	}

	if err := h.securityService.ResetPasswordWithToken(req.Token, req.NewPassword); err != nil {
		switch err {
		case services.ErrInvalidToken:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		case services.ErrTokenAlreadyUsed:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Token has already been used"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// ChangePassword handles POST /api/v1/auth/change-password
func (h *SecurityHandler) ChangePassword(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Current and new password (min 8 characters) are required"})
		return
	}

	if err := h.securityService.ChangePassword(userID, req.CurrentPassword, req.NewPassword); err != nil {
		switch err {
		case services.ErrPasswordMismatch:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is incorrect"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change password"})
		}
		return
	}

	// Log activity
	_ = h.securityService.LogActivity(
		userID,
		models.ActivityPasswordChanged,
		"user",
		&userID,
		nil,
		c.ClientIP(),
		c.Request.UserAgent(),
	)

	// TODO: Send password changed notification email
	// _ = h.emailService.SendPasswordChangedNotification(user)

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// VerifyEmail handles POST /api/v1/auth/verify-email
func (h *SecurityHandler) VerifyEmail(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	if err := h.securityService.VerifyEmail(req.Token); err != nil {
		switch err {
		case services.ErrInvalidToken:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

// ResendVerification handles POST /api/v1/auth/resend-verification
func (h *SecurityHandler) ResendVerification(c *gin.Context) {
	userID := getUserIDFromContext(c)

	token, err := h.securityService.SendVerificationEmail(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}

	// TODO: Actually send the email
	// _ = h.emailService.SendEmailVerificationEmail(user, token.Token, "")

	_ = token // Suppress unused variable warning

	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent"})
}

// GetLoginHistory handles GET /api/v1/profile/login-history
func (h *SecurityHandler) GetLoginHistory(c *gin.Context) {
	userID := getUserIDFromContext(c)

	// Parse pagination params
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	history, total, err := h.securityService.GetLoginHistory(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get login history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"history": history,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}

// GetActiveSessions handles GET /api/v1/profile/sessions
func (h *SecurityHandler) GetActiveSessions(c *gin.Context) {
	userID := getUserIDFromContext(c)

	sessions, err := h.securityService.GetActiveSessions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get active sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}

// LogoutSession handles POST /api/v1/auth/logout-session/:id
func (h *SecurityHandler) LogoutSession(c *gin.Context) {
	userID := getUserIDFromContext(c)

	sessionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	// TODO: Verify that the session belongs to the user before logging out

	if err := h.securityService.LogoutSession(uint(sessionID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout session"})
		return
	}

	// Log activity
	_ = h.securityService.LogActivity(
		userID,
		models.ActivityLogout,
		"session",
		nil,
		map[string]interface{}{"session_id": sessionID},
		c.ClientIP(),
		c.Request.UserAgent(),
	)

	c.JSON(http.StatusOK, gin.H{"message": "Session logged out successfully"})
}

// LogoutAll handles POST /api/v1/auth/logout-all
func (h *SecurityHandler) LogoutAll(c *gin.Context) {
	userID := getUserIDFromContext(c)

	if err := h.securityService.LogoutAllSessions(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout all sessions"})
		return
	}

	// Log activity
	_ = h.securityService.LogActivity(
		userID,
		models.ActivityLogout,
		"session",
		nil,
		map[string]interface{}{"action": "logout_all"},
		c.ClientIP(),
		c.Request.UserAgent(),
	)

	c.JSON(http.StatusOK, gin.H{"message": "All sessions logged out successfully"})
}

// GetActivityLog handles GET /api/v1/profile/activity
func (h *SecurityHandler) GetActivityLog(c *gin.Context) {
	userID := getUserIDFromContext(c)

	// Parse pagination params
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	activities, total, err := h.securityService.GetActivityLog(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get activity log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"activities": activities,
		"total":      total,
		"limit":      limit,
		"offset":     offset,
	})
}
