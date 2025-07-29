package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	// Add user service when implemented
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	// Placeholder: return empty list
	c.JSON(http.StatusOK, gin.H{"users": []interface{}{}})
}
