package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"template-store/internal/handlers"
	"template-store/internal/models"
	"template-store/internal/services"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Set up logging
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	r := gin.Default()

	// Add CORS middleware
	r.Use(gin.Recovery())

	// Connect to the database
	db, err := services.ConnectDB()
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	// Auto-migrate User model
	if err := models.AutoMigrate(db); err != nil {
		logger.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize services
	templateService := services.NewTemplateService(db)
	categoryService := services.NewCategoryService(db)
	blogService := services.NewBlogService(db)

	// Initialize handlers
	templateHandler := handlers.NewTemplateHandler(templateService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	blogHandler := handlers.NewBlogHandler(blogService)
	userHandler := handlers.NewUserHandler()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "template-store",
			"version": "1.0.0",
		})
	})

	// API routes group
	api := r.Group("/api/v1")
	{
		// Placeholder for future API endpoints
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Template Store API v1",
			})
		})
		api.GET("/users", userHandler.ListUsers)

		// Template routes
		templates := api.Group("/templates")
		{
			templates.GET("", templateHandler.ListTemplates)
			templates.POST("", templateHandler.CreateTemplate)
			templates.GET("/:id", templateHandler.GetTemplate)
			templates.PUT("/:id", templateHandler.UpdateTemplate)
			templates.DELETE("/:id", templateHandler.DeleteTemplate)
			templates.GET("/category/:category_id", templateHandler.GetTemplatesByCategory)
		}

		// Category routes
		categories := api.Group("/categories")
		{
			categories.GET("", categoryHandler.ListCategories)
			categories.POST("", categoryHandler.CreateCategory)
			categories.GET("/:id", categoryHandler.GetCategory)
			categories.PUT("/:id", categoryHandler.UpdateCategory)
			categories.DELETE("/:id", categoryHandler.DeleteCategory)
			categories.POST("/seed", categoryHandler.SeedCategories)
		}

		// Blog routes
		blog := api.Group("/blog")
		{
			blog.GET("", blogHandler.ListBlogPosts)
			blog.POST("", blogHandler.CreateBlogPost)
			blog.GET("/:id", blogHandler.GetBlogPost)
			blog.PUT("/:id", blogHandler.UpdateBlogPost)
			blog.DELETE("/:id", blogHandler.DeleteBlogPost)
			blog.GET("/category/:category_id", blogHandler.GetBlogPostsByCategory)
			blog.GET("/author/:author_id", blogHandler.GetBlogPostsByAuthor)
		}
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Infof("Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
