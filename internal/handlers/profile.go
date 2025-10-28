package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"template-store/internal/models"
	"template-store/internal/services"
)

// ProfileHandler handles profile-related HTTP requests
type ProfileHandler struct {
	profileService  services.ProfileService
	securityService services.SecurityService
}

// NewProfileHandler creates a new ProfileHandler
func NewProfileHandler(profileService services.ProfileService, securityService services.SecurityService) *ProfileHandler {
	return &ProfileHandler{
		profileService:  profileService,
		securityService: securityService,
	}
}

// GetProfile handles GET /api/v1/profile
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID := getUserIDFromContext(c)

	user, err := h.profileService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile handles PUT /api/v1/profile
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := h.profileService.UpdateProfile(userID, updates)
	if err != nil {
		if err == services.ErrInvalidProfileData {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	// Log activity
	_ = h.securityService.LogActivity(
		userID,
		models.ActivityProfileUpdated,
		"user",
		&userID,
		updates,
		c.ClientIP(),
		c.Request.UserAgent(),
	)

	c.JSON(http.StatusOK, user)
}

// UploadAvatar handles POST /api/v1/profile/avatar
func (h *ProfileHandler) UploadAvatar(c *gin.Context) {
	userID := getUserIDFromContext(c)

	// Get file from form
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Avatar file is required"})
		return
	}

	// Upload avatar
	avatarURL, err := h.profileService.UploadAvatar(c.Request.Context(), userID, file)
	if err != nil {
		switch err {
		case services.ErrAvatarTooLarge:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Avatar file too large (max 5MB)"})
		case services.ErrInvalidAvatarFormat:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid avatar format (allowed: jpg, jpeg, png, webp)"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload avatar"})
		}
		return
	}

	// Log activity
	_ = h.securityService.LogActivity(
		userID,
		models.ActivityAvatarUploaded,
		"user",
		&userID,
		map[string]interface{}{"avatar_url": avatarURL},
		c.ClientIP(),
		c.Request.UserAgent(),
	)

	c.JSON(http.StatusOK, gin.H{"avatar_url": avatarURL})
}

// DeleteAvatar handles DELETE /api/v1/profile/avatar
func (h *ProfileHandler) DeleteAvatar(c *gin.Context) {
	userID := getUserIDFromContext(c)

	if err := h.profileService.DeleteAvatar(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete avatar"})
		return
	}

	// Log activity
	_ = h.securityService.LogActivity(
		userID,
		models.ActivityAvatarDeleted,
		"user",
		&userID,
		nil,
		c.ClientIP(),
		c.Request.UserAgent(),
	)

	c.Status(http.StatusNoContent)
}

// GetPublicProfile handles GET /api/v1/users/:id/profile
func (h *ProfileHandler) GetPublicProfile(c *gin.Context) {
	// Get target user ID from URL
	targetUserID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get viewer ID (may be nil if not authenticated)
	var viewerID *uint
	if uid := getUserIDFromContextOptional(c); uid != 0 {
		viewerID = &uid
	}

	user, err := h.profileService.GetPublicProfile(uint(targetUserID), viewerID)
	if err != nil {
		if err == services.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetPreferences handles GET /api/v1/profile/preferences
func (h *ProfileHandler) GetPreferences(c *gin.Context) {
	userID := getUserIDFromContext(c)

	prefs, err := h.profileService.GetPreferences(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get preferences"})
		return
	}

	c.JSON(http.StatusOK, prefs)
}

// UpdatePreferences handles PUT /api/v1/profile/preferences
func (h *ProfileHandler) UpdatePreferences(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var prefs models.UserPreferences
	if err := c.ShouldBindJSON(&prefs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.profileService.UpdatePreferences(userID, &prefs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update preferences"})
		return
	}

	// Log activity
	_ = h.securityService.LogActivity(
		userID,
		models.ActivityPreferencesUpdated,
		"user",
		&userID,
		nil,
		c.ClientIP(),
		c.Request.UserAgent(),
	)

	c.JSON(http.StatusOK, prefs)
}

// DeactivateAccount handles POST /api/v1/profile/deactivate
func (h *ProfileHandler) DeactivateAccount(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req struct {
		Password string `json:"password" binding:"required"`
		Reason   string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
		return
	}

	// TODO: Verify password with Cognito before deactivating

	if err := h.profileService.DeactivateAccount(userID, req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate account"})
		return
	}

	// Log activity
	_ = h.securityService.LogActivity(
		userID,
		models.ActivityAccountDeactivated,
		"user",
		&userID,
		map[string]interface{}{"reason": req.Reason},
		c.ClientIP(),
		c.Request.UserAgent(),
	)

	c.JSON(http.StatusOK, gin.H{"message": "Account deactivated successfully"})
}

// DeleteAccount handles DELETE /api/v1/profile
func (h *ProfileHandler) DeleteAccount(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req struct {
		Password     string `json:"password" binding:"required"`
		Confirmation string `json:"confirmation" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password and confirmation are required"})
		return
	}

	// Verify confirmation text
	if req.Confirmation != "DELETE MY ACCOUNT" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid confirmation text"})
		return
	}

	// TODO: Verify password with Cognito before deleting

	// Log activity before deletion
	_ = h.securityService.LogActivity(
		userID,
		models.ActivityAccountDeleted,
		"user",
		&userID,
		nil,
		c.ClientIP(),
		c.Request.UserAgent(),
	)

	if err := h.profileService.DeleteAccount(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

// Helper functions

// getUserIDFromContext extracts user ID from the request context
// This assumes the auth middleware sets the user ID in the context
func getUserIDFromContext(c *gin.Context) uint {
	// This will need to be implemented based on your auth middleware
	// For now, return a placeholder
	userID, exists := c.Get("userID")
	if !exists {
		return 0
	}
	if id, ok := userID.(uint); ok {
		return id
	}
	return 0
}

// getUserIDFromContextOptional gets user ID if authenticated, returns 0 if not
func getUserIDFromContextOptional(c *gin.Context) uint {
	userID, exists := c.Get("userID")
	if !exists {
		return 0
	}
	if id, ok := userID.(uint); ok {
		return id
	}
	return 0
}
