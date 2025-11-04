package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"template-store/internal/models"
)

// RequireRole creates middleware that requires user to have one of the specified roles
func RequireRole(db *gorm.DB, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by auth middleware)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// Get user from database
		var user models.User
		if err := db.Select("role").First(&user, userID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify permissions"})
			c.Abort()
			return
		}

		// Check if user has one of the required roles
		hasRole := false
		for _, role := range roles {
			if user.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		// Store role in context for handlers to use
		c.Set("userRole", user.Role)
		c.Next()
	}
}

// RequireAdmin is a shorthand for RequireRole("admin")
func RequireAdmin(db *gorm.DB) gin.HandlerFunc {
	return RequireRole(db, "admin")
}

// RequireAdminOrAuthor requires user to be admin or author
func RequireAdminOrAuthor(db *gorm.DB) gin.HandlerFunc {
	return RequireRole(db, "admin", "author")
}

// RequireActiveAccount ensures user account is active (not suspended or deleted)
func RequireActiveAccount(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// Get user status from database
		var user models.User
		if err := db.Select("status").First(&user, userID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify account status"})
			c.Abort()
			return
		}

		// Check account status
		if user.Status != "active" {
			var message string
			switch user.Status {
			case "suspended":
				message = "Your account has been suspended"
			case "deleted":
				message = "Your account has been deleted"
			default:
				message = "Your account is not active"
			}
			c.JSON(http.StatusForbidden, gin.H{"error": message})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireEmailVerified ensures user has verified their email
func RequireEmailVerified(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// Get user email verification status
		var user models.User
		if err := db.Select("email_verified").First(&user, userID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email status"})
			c.Abort()
			return
		}

		// Check if email is verified
		if !user.EmailVerified {
			c.JSON(http.StatusForbidden, gin.H{"error": "Email verification required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth is middleware that sets user context if authenticated, but doesn't require it
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is authenticated
		_, exists := c.Get("userID")

		// Set a flag to indicate whether user is authenticated
		c.Set("isAuthenticated", exists)

		// Always continue, even if not authenticated
		c.Next()
	}
}
