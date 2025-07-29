package services

import (
	"errors"
	"strings"
	"template-store/internal/models"
	"time"

	"github.com/gomarkdown/markdown"
	"gorm.io/gorm"
)

type BlogService struct {
	db *gorm.DB
}

func NewBlogService(db *gorm.DB) *BlogService {
	return &BlogService{db: db}
}

// CreateBlogPost creates a new blog post
func (s *BlogService) CreateBlogPost(post *models.BlogPost) error {
	if post.Title == "" {
		return errors.New("blog post title is required")
	}
	if post.Content == "" {
		return errors.New("blog post content is required")
	}
	if post.AuthorID == 0 {
		return errors.New("author ID is required")
	}
	
	// Set default values
	if post.CreatedAt.IsZero() {
		post.CreatedAt = time.Now()
	}
	if post.UpdatedAt.IsZero() {
		post.UpdatedAt = time.Now()
	}
	
	return s.db.Create(post).Error
}

// GetBlogPost retrieves a blog post by ID
func (s *BlogService) GetBlogPost(id uint) (*models.BlogPost, error) {
	var post models.BlogPost
	err := s.db.Preload("Author").Preload("Category").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// ListBlogPosts retrieves all blog posts with optional filtering
func (s *BlogService) ListBlogPosts(categoryID *uint, search *string, limit, offset int) ([]models.BlogPost, int64, error) {
	var posts []models.BlogPost
	var total int64
	
	query := s.db.Model(&models.BlogPost{}).Preload("Author").Preload("Category")
	
	// Apply category filter
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	
	// Apply search filter
	if search != nil && *search != "" {
		query = query.Where("title ILIKE ? OR content ILIKE ?", "%"+*search+"%", "%"+*search+"%")
	}
	
	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination and ordering (newest first)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	err := query.Order("created_at DESC").Find(&posts).Error
	return posts, total, err
}

// UpdateBlogPost updates an existing blog post
func (s *BlogService) UpdateBlogPost(id uint, updates map[string]interface{}) error {
	if title, exists := updates["title"]; exists && title == "" {
		return errors.New("blog post title cannot be empty")
	}
	if content, exists := updates["content"]; exists && content == "" {
		return errors.New("blog post content cannot be empty")
	}
	
	// Update the updated_at timestamp
	updates["updated_at"] = time.Now()
	
	return s.db.Model(&models.BlogPost{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteBlogPost deletes a blog post by ID
func (s *BlogService) DeleteBlogPost(id uint) error {
	return s.db.Delete(&models.BlogPost{}, id).Error
}

// GetBlogPostsByCategory retrieves blog posts filtered by category
func (s *BlogService) GetBlogPostsByCategory(categoryID uint) ([]models.BlogPost, error) {
	var posts []models.BlogPost
	err := s.db.Where("category_id = ?", categoryID).Preload("Author").Preload("Category").Order("created_at DESC").Find(&posts).Error
	return posts, err
}

// GetBlogPostsByAuthor retrieves blog posts filtered by author
func (s *BlogService) GetBlogPostsByAuthor(authorID uint) ([]models.BlogPost, error) {
	var posts []models.BlogPost
	err := s.db.Where("author_id = ?", authorID).Preload("Author").Preload("Category").Order("created_at DESC").Find(&posts).Error
	return posts, err
}

// ProcessMarkdown converts markdown content to HTML
func (s *BlogService) ProcessMarkdown(content string) string {
	md := []byte(content)
	html := markdown.ToHTML(md, nil, nil)
	return string(html)
}

// GenerateExcerpt generates a short excerpt from blog content
func (s *BlogService) GenerateExcerpt(content string, maxLength int) string {
	// Remove markdown formatting for excerpt
	plainText := s.stripMarkdown(content)
	
	if len(plainText) <= maxLength {
		return plainText
	}
	
	// Truncate and add ellipsis
	excerpt := plainText[:maxLength]
	lastSpace := strings.LastIndex(excerpt, " ")
	if lastSpace > 0 {
		excerpt = excerpt[:lastSpace]
	}
	
	return excerpt + "..."
}

// stripMarkdown removes markdown formatting from text
func (s *BlogService) stripMarkdown(content string) string {
	// Simple markdown stripping - in production, you might want a more robust solution
	content = strings.ReplaceAll(content, "#", "")
	content = strings.ReplaceAll(content, "*", "")
	content = strings.ReplaceAll(content, "**", "")
	content = strings.ReplaceAll(content, "`", "")
	content = strings.ReplaceAll(content, "```", "")
	content = strings.ReplaceAll(content, ">", "")
	content = strings.ReplaceAll(content, "-", "")
	content = strings.ReplaceAll(content, "+", "")
	
	// Remove extra whitespace
	content = strings.Join(strings.Fields(content), " ")
	
	return content
} 