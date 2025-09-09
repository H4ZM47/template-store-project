package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"template-store/internal/config"
	"template-store/internal/handlers"
	"template-store/internal/middleware"
	"template-store/internal/models"
	"template-store/internal/services"
)

func main() {
	// Declare all variables that will be initialized with a potential error.
	var db *gorm.DB
	var authService services.AuthService
	var storageService services.StorageService
	var paymentService services.PaymentService
	var err error

	// Load environment variables from .env file
	if err = godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Set up logging
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// Enable permissive CORS for development
	r.Use(cors.Default())

	// Connect to the database
	db, err = services.ConnectDB()
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	// Auto-migrate models
	if err = models.AutoMigrate(db); err != nil {
		logger.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize services
	if gin.Mode() == gin.DebugMode {
		authService, err = services.NewMockAuthService(cfg)
	} else {
		authService, err = services.NewAuthService(cfg)
	}
	if err != nil {
		logger.Fatalf("Failed to initialize auth service: %v", err)
	}

	storageService, err = services.NewStorageService(cfg)
	if err != nil {
		logger.Fatalf("Failed to initialize storage service: %v", err)
	}

	paymentService, err = services.NewPaymentService(cfg)
	if err != nil {
		logger.Fatalf("Failed to initialize payment service: %v", err)
	}
	templateService := services.NewTemplateService(db)
	categoryService := services.NewCategoryService(db)
	blogService := services.NewBlogService(db)
	userService := services.NewUserService(db)
	orderService := services.NewOrderService(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	templateHandler := handlers.NewTemplateHandler(templateService, storageService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	blogHandler := handlers.NewBlogHandler(blogService)
	userHandler := handlers.NewUserHandler(userService)
	paymentHandler := handlers.NewPaymentHandler(paymentService, templateService, orderService, userService, cfg.Stripe.WebhookSecret)

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
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Template Store API v1",
			})
		})

		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Webhook routes
		webhooks := api.Group("/webhooks")
		{
			webhooks.POST("/stripe", paymentHandler.StripeWebhook)
		}

		// Authenticated routes group
		authenticated := api.Group("/")
		if gin.Mode() != gin.DebugMode {
			authenticated.Use(middleware.AuthMiddleware(cfg))
		} else {
			// In debug mode, use a dummy middleware that does nothing
			authenticated.Use(func(c *gin.Context) {
				c.Next()
			})
		}
		{
			// Checkout route
			authenticated.POST("/checkout", paymentHandler.CreateCheckoutSession)

			// User routes
			users := authenticated.Group("/users")
			{
				users.GET("", userHandler.ListUsers)
				users.POST("", userHandler.CreateUser)
				users.GET("/:id", userHandler.GetUser)
				users.POST("/seed", userHandler.SeedUsers)
			}

			// Template routes
			templates := authenticated.Group("/templates")
			{
				templates.POST("", templateHandler.CreateTemplate)
				templates.PUT("/:id", templateHandler.UpdateTemplate)
				templates.DELETE("/:id", templateHandler.DeleteTemplate)
			}

			// Blog routes
			blog := authenticated.Group("/blog")
			{
				blog.POST("", blogHandler.CreateBlogPost)
				blog.PUT("/:id", blogHandler.UpdateBlogPost)
				blog.DELETE("/:id", blogHandler.DeleteBlogPost)
			}
		}

		// Public routes
		public := api.Group("/")
		{
			// Payment routes
			payments := public.Group("/payment")
			{
				payments.GET("/success", paymentHandler.PaymentSuccess)
				payments.GET("/cancel", paymentHandler.PaymentCancel)
			}

			// Template routes
			templates := public.Group("/templates")
			{
				templates.GET("", templateHandler.ListTemplates)
				templates.GET("/:id", templateHandler.GetTemplate)
				templates.GET("/category/:category_id", templateHandler.GetTemplatesByCategory)
				templates.POST("/seed", templateHandler.SeedTemplates)
			}

			// Category routes
			categories := public.Group("/categories")
			{
				categories.GET("", categoryHandler.ListCategories)
				categories.POST("", categoryHandler.CreateCategory)
				categories.GET("/:id", categoryHandler.GetCategory)
				categories.PUT("/:id", categoryHandler.UpdateCategory)
				categories.DELETE("/:id", categoryHandler.DeleteCategory)
				categories.POST("/seed", categoryHandler.SeedCategories)
			}

			// Blog routes
			blog := public.Group("/blog")
			{
				blog.GET("", blogHandler.ListBlogPosts)
				blog.GET("/:id", blogHandler.GetBlogPost)
				blog.GET("/category/:category_id", blogHandler.GetBlogPostsByCategory)
				blog.GET("/author/:author_id", blogHandler.GetBlogPostsByAuthor)
			}
		}
	}

	// Get port from environment or use default
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	logger.Infof("Starting server on port %s", port)
	if err = r.Run(":" + port); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
