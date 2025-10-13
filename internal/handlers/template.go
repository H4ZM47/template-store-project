package handlers

import (
	"net/http"
	"strconv"
	"template-store/internal/models"
	"template-store/internal/services"

	"github.com/gin-gonic/gin"
)

type TemplateHandler struct {
	templateService services.TemplateService
	storageService  services.StorageService
}

func NewTemplateHandler(templateService services.TemplateService, storageService services.StorageService) *TemplateHandler {
	return &TemplateHandler{
		templateService: templateService,
		storageService:  storageService,
	}
}

// CreateTemplate handles POST /api/v1/templates
func (h *TemplateHandler) CreateTemplate(c *gin.Context) {
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing form"})
		return
	}

	// Get form values
	name := c.Request.FormValue("name")
	priceStr := c.Request.FormValue("price")
	categoryIDStr := c.Request.FormValue("category_id")
	previewData := c.Request.FormValue("preview_data")

	if name == "" || priceStr == "" || categoryIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields: name, price, category_id"})
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price format"})
		return
	}

	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category_id format"})
		return
	}

	// Get file from form
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	// Upload file to S3
	fileURL, err := h.storageService.UploadFile(c.Request.Context(), fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	// Create template model
	template := models.Template{
		Name:        name,
		Price:       price,
		CategoryID:  uint(categoryID),
		FileInfo:    fileURL,
		PreviewData: previewData,
	}

	if err := h.templateService.CreateTemplate(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Template created successfully",
		"template": template,
	})
}

// GetTemplate handles GET /api/v1/templates/:id
func (h *TemplateHandler) GetTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	template, err := h.templateService.GetTemplate(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"template": template})
}

// ListTemplates handles GET /api/v1/templates
func (h *TemplateHandler) ListTemplates(c *gin.Context) {
	// Parse query parameters
	search := c.Query("search")
	categoryIDStr := c.Query("category_id")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	var categoryID *uint
	if categoryIDStr != "" {
		if id, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			catID := uint(id)
			categoryID = &catID
		}
	}

	var searchPtr *string
	if search != "" {
		searchPtr = &search
	}

	templates, total, err := h.templateService.ListTemplates(categoryID, searchPtr, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve templates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"templates": templates,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	})
}

// UpdateTemplate handles PUT /api/v1/templates/:id
func (h *TemplateHandler) UpdateTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.templateService.UpdateTemplate(uint(id), updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template updated successfully"})
}

// DeleteTemplate handles DELETE /api/v1/templates/:id
func (h *TemplateHandler) DeleteTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	if err := h.templateService.DeleteTemplate(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete template"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})
}

// GetTemplatesByCategory handles GET /api/v1/templates/category/:category_id
func (h *TemplateHandler) GetTemplatesByCategory(c *gin.Context) {
	categoryIDStr := c.Param("category_id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	templates, err := h.templateService.GetTemplatesByCategory(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve templates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

// SeedTemplates handles POST /api/v1/templates/seed
func (h *TemplateHandler) SeedTemplates(c *gin.Context) {
	if err := h.templateService.SeedTemplates(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to seed templates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Templates seeded successfully"})
}

// ViewTemplate handles GET /api/v1/templates/:id/view
func (h *TemplateHandler) ViewTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	template, err := h.templateService.GetTemplate(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	// Map template names to HTML files
	templateFiles := map[string]string{
		"Data Classification Standard":       "web/templates/data-classification-standard.html",
		"Vulnerability Management Standard":  "web/templates/vulnerability-management-standard.html",
	}

	templateFile, exists := templateFiles[template.Name]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template file not found"})
		return
	}

	// Serve the HTML file
	c.File(templateFile)
}
