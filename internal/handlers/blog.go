package handlers

import (
	"net/http"
	"strconv"
	"template-store/internal/models"
	"template-store/internal/services"

	"github.com/gin-gonic/gin"
)

type BlogHandler struct {
	blogService *services.BlogService
}

func NewBlogHandler(blogService *services.BlogService) *BlogHandler {
	return &BlogHandler{
		blogService: blogService,
	}
}

// CreateBlogPost handles POST /api/v1/blog
func (h *BlogHandler) CreateBlogPost(c *gin.Context) {
	var post models.BlogPost
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.blogService.CreateBlogPost(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Blog post created successfully",
		"post":    post,
	})
}

// GetBlogPost handles GET /api/v1/blog/:id
func (h *BlogHandler) GetBlogPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog post ID"})
		return
	}

	post, err := h.blogService.GetBlogPost(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog post not found"})
		return
	}

	// Process markdown content
	processedContent := h.blogService.ProcessMarkdown(post.Content)
	
	response := gin.H{
		"post": gin.H{
			"id":         post.ID,
			"title":      post.Title,
			"content":    post.Content,
			"html_content": processedContent,
			"author_id":  post.AuthorID,
			"author":     post.Author,
			"category_id": post.CategoryID,
			"category":   post.Category,
			"seo":        post.SEO,
			"created_at": post.CreatedAt,
			"updated_at": post.UpdatedAt,
		},
	}

	c.JSON(http.StatusOK, response)
}

// ListBlogPosts handles GET /api/v1/blog
func (h *BlogHandler) ListBlogPosts(c *gin.Context) {
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

	posts, total, err := h.blogService.ListBlogPosts(categoryID, searchPtr, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve blog posts"})
		return
	}

	// Process posts to include excerpts and HTML content
	var processedPosts []gin.H
	for _, post := range posts {
		excerpt := h.blogService.GenerateExcerpt(post.Content, 150)
		processedPosts = append(processedPosts, gin.H{
			"id":         post.ID,
			"title":      post.Title,
			"excerpt":    excerpt,
			"author_id":  post.AuthorID,
			"author":     post.Author,
			"category_id": post.CategoryID,
			"category":   post.Category,
			"seo":        post.SEO,
			"created_at": post.CreatedAt,
			"updated_at": post.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"posts":  processedPosts,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// UpdateBlogPost handles PUT /api/v1/blog/:id
func (h *BlogHandler) UpdateBlogPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog post ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.blogService.UpdateBlogPost(uint(id), updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog post updated successfully"})
}

// DeleteBlogPost handles DELETE /api/v1/blog/:id
func (h *BlogHandler) DeleteBlogPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog post ID"})
		return
	}

	if err := h.blogService.DeleteBlogPost(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete blog post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog post deleted successfully"})
}

// GetBlogPostsByCategory handles GET /api/v1/blog/category/:category_id
func (h *BlogHandler) GetBlogPostsByCategory(c *gin.Context) {
	categoryIDStr := c.Param("category_id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	posts, err := h.blogService.GetBlogPostsByCategory(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve blog posts"})
		return
	}

	// Process posts to include excerpts
	var processedPosts []gin.H
	for _, post := range posts {
		excerpt := h.blogService.GenerateExcerpt(post.Content, 150)
		processedPosts = append(processedPosts, gin.H{
			"id":         post.ID,
			"title":      post.Title,
			"excerpt":    excerpt,
			"author_id":  post.AuthorID,
			"author":     post.Author,
			"category_id": post.CategoryID,
			"category":   post.Category,
			"seo":        post.SEO,
			"created_at": post.CreatedAt,
			"updated_at": post.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"posts": processedPosts})
}

// GetBlogPostsByAuthor handles GET /api/v1/blog/author/:author_id
func (h *BlogHandler) GetBlogPostsByAuthor(c *gin.Context) {
	authorIDStr := c.Param("author_id")
	authorID, err := strconv.ParseUint(authorIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	posts, err := h.blogService.GetBlogPostsByAuthor(uint(authorID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve blog posts"})
		return
	}

	// Process posts to include excerpts
	var processedPosts []gin.H
	for _, post := range posts {
		excerpt := h.blogService.GenerateExcerpt(post.Content, 150)
		processedPosts = append(processedPosts, gin.H{
			"id":         post.ID,
			"title":      post.Title,
			"excerpt":    excerpt,
			"author_id":  post.AuthorID,
			"author":     post.Author,
			"category_id": post.CategoryID,
			"category":   post.Category,
			"seo":        post.SEO,
			"created_at": post.CreatedAt,
			"updated_at": post.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"posts": processedPosts})
} 