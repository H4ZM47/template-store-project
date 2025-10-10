package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"template-store/internal/models"
	"template-store/internal/services"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
)

type TemplateHandler struct {
	templateService *services.TemplateService
}

func NewTemplateHandler(templateService *services.TemplateService) *TemplateHandler {
	return &TemplateHandler{
		templateService: templateService,
	}
}

// CreateTemplate handles POST /api/v1/templates
func (h *TemplateHandler) CreateTemplate(c *gin.Context) {
	var template models.Template
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
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

// ViewTemplate handles GET /api/v1/templates/:id/view - serves the HTML template
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

	// Get the web directory path
	cwd, _ := os.Getwd()
	htmlPath := filepath.Join(cwd, "web", template.FileInfo)

	// Check if file exists
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template file not found"})
		return
	}

	// Serve the HTML file
	c.File(htmlPath)
}

// DownloadTemplatePDF handles GET /api/v1/templates/:id/download - generates and serves PDF
func (h *TemplateHandler) DownloadTemplatePDF(c *gin.Context) {
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

	// Get the web directory path
	cwd, _ := os.Getwd()
	htmlPath := filepath.Join(cwd, "web", template.FileInfo)

	// Check if file exists
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template file not found"})
		return
	}

	// Read the HTML file
	htmlContent, err := os.ReadFile(htmlPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read template"})
		return
	}

	// Generate PDF from HTML
	pdfBytes, err := h.generatePDFFromHTML(string(htmlContent))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate PDF: %v", err)})
		return
	}

	// Create safe filename
	safeFilename := strings.ReplaceAll(template.Name, " ", "_")
	safeFilename = strings.ToLower(safeFilename)

	// Set headers for PDF download
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.pdf\"", safeFilename))
	c.Header("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))

	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// generatePDFFromHTML converts HTML to PDF using chromedp
func (h *TemplateHandler) generatePDFFromHTML(htmlContent string) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var pdfBuffer []byte

	// URL encode the HTML content and use data URL
	// Navigate to data URL with HTML content and print to PDF
	err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Set the page content
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}

			return page.SetDocumentContent(frameTree.Frame.ID, htmlContent).Do(ctx)
		}),
		// Wait for page to load
		chromedp.Sleep(2*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			// Print to PDF with proper settings
			pdfBuffer, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPreferCSSPageSize(true).
				Do(ctx)
			return err
		}),
	)

	if err != nil {
		return nil, err
	}

	return pdfBuffer, nil
}

// GetTemplateThumbnail handles GET /api/v1/templates/:id/thumbnail - serves thumbnail image
func (h *TemplateHandler) GetTemplateThumbnail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	// Get the thumbnail HTML path based on template ID or file name
	cwd, _ := os.Getwd()

	// Get template to determine thumbnail path
	template, err := h.templateService.GetTemplate(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	// Extract base filename from template file_info
	baseFilename := filepath.Base(template.FileInfo)
	ext := filepath.Ext(baseFilename)
	thumbFilename := baseFilename[:len(baseFilename)-len(ext)] + "-thumb.html"

	thumbPath := filepath.Join(cwd, "web", "templates", "thumbnails", thumbFilename)

	// Check if thumbnail file exists
	if _, err := os.Stat(thumbPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Thumbnail not found"})
		return
	}

	// Read the HTML file
	htmlContent, err := os.ReadFile(thumbPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read thumbnail"})
		return
	}

	// Generate PNG from HTML
	pngBytes, err := h.generateThumbnailFromHTML(string(htmlContent))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate thumbnail: %v", err)})
		return
	}

	// Set headers for image
	c.Header("Content-Type", "image/png")
	c.Header("Cache-Control", "public, max-age=86400")

	c.Data(http.StatusOK, "image/png", pngBytes)
}

// generateThumbnailFromHTML converts HTML to PNG using chromedp
func (h *TemplateHandler) generateThumbnailFromHTML(htmlContent string) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var pngBuffer []byte

	// Navigate and capture screenshot
	err := chromedp.Run(ctx,
		chromedp.EmulateViewport(400, 300),
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Set the page content
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}

			return page.SetDocumentContent(frameTree.Frame.ID, htmlContent).Do(ctx)
		}),
		// Wait for page to load
		chromedp.Sleep(1*time.Second),
		chromedp.FullScreenshot(&pngBuffer, 100),
	)

	if err != nil {
		return nil, err
	}

	return pngBuffer, nil
}

// GenerateCustomPDF handles POST /api/v1/templates/:id/generate - generates PDF with custom variables
func (h *TemplateHandler) GenerateCustomPDF(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	// Get template
	template, err := h.templateService.GetTemplate(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	// Parse variables from request body
	var variables map[string]string
	if err := c.ShouldBindJSON(&variables); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid variables format"})
		return
	}

	// Validate required variables
	for _, varDef := range template.Variables {
		if varDef.Required {
			if val, ok := variables[varDef.Name]; !ok || val == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Required variable '%s' is missing", varDef.Label)})
				return
			}
		}
	}

	// Get the web directory path
	cwd, _ := os.Getwd()
	htmlPath := filepath.Join(cwd, "web", template.FileInfo)

	// Check if file exists
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template file not found"})
		return
	}

	// Read the HTML file
	htmlContent, err := os.ReadFile(htmlPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read template"})
		return
	}

	// Replace variables in HTML
	processedHTML := h.replaceVariables(string(htmlContent), variables)

	// Generate PDF from processed HTML
	pdfBytes, err := h.generatePDFFromHTML(processedHTML)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate PDF: %v", err)})
		return
	}

	// Create safe filename
	safeFilename := strings.ReplaceAll(template.Name, " ", "_")
	safeFilename = strings.ToLower(safeFilename)

	// Set headers for PDF download
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s_custom.pdf\"", safeFilename))
	c.Header("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))

	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// replaceVariables replaces all {{VARIABLE}} placeholders with actual values
func (h *TemplateHandler) replaceVariables(html string, variables map[string]string) string {
	result := html
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// GetTemplateVariables handles GET /api/v1/templates/:id/variables - returns variable definitions
func (h *TemplateHandler) GetTemplateVariables(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"variables": template.Variables})
}
