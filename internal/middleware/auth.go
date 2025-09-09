package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	cognitojwt "github.com/jhosan7/cognito-jwt-verify"
	"template-store/internal/config"
)

// AuthMiddleware creates a Gin middleware for JWT authentication.
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	// Initialize the Cognito JWT validator
	cognitoConfig := cognitojwt.Config{
		UserPoolId: cfg.AWS.CognitoPoolID,
		ClientId:   cfg.AWS.CognitoAppClientID,
		TokenUse:   "access", // Can be "id" or "access"
	}

	validator, err := cognitojwt.Create(cognitoConfig)
	if err != nil {
		// This would be a configuration error, so it's okay to panic at startup
		panic("Failed to initialize Cognito JWT validator: " + err.Error())
	}

	return func(c *gin.Context) {
		// Get token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		// The token is expected to be in the format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}
		token := parts[1]

		// Verify the token
		claims, err := validator.Verify(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			return
		}

		// Set the user claims in the context for downstream handlers
		c.Set("user_claims", claims)

		c.Next()
	}
}
