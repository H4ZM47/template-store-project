package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"template-store/internal/handlers"
	"template-store/internal/middleware"
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
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// Enable permissive CORS for development
	r.Use(cors.Default())

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
	userService := services.NewUserService(db)
	orderService := services.NewOrderService(db)

	// Initialize JWT service
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-this-in-production"
		logger.Warn("JWT_SECRET not set, using default (unsafe for production)")
	}
	jwtIssuer := os.Getenv("JWT_ISSUER")
	if jwtIssuer == "" {
		jwtIssuer = "template-store"
	}
	jwtService := services.NewJWTService(jwtSecret, jwtIssuer, 24) // 24 hours

	// Initialize auth service
	authService := services.NewAuthService(db, jwtService)

	// Initialize Stripe service
	stripeAPIKey := os.Getenv("STRIPE_API_KEY")
	stripeWebhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	stripeSuccessURL := os.Getenv("STRIPE_SUCCESS_URL")
	stripeCancelURL := os.Getenv("STRIPE_CANCEL_URL")
	stripeService := services.NewStripeService(stripeAPIKey, stripeWebhookSecret, stripeSuccessURL, stripeCancelURL)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	templateHandler := handlers.NewTemplateHandler(templateService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	blogHandler := handlers.NewBlogHandler(blogService)
	userHandler := handlers.NewUserHandler(userService)
	paymentHandler := handlers.NewPaymentHandler(stripeService, orderService, templateService, userService)

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

		// Public auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/password/reset-request", authHandler.RequestPasswordReset)
			auth.POST("/password/reset", authHandler.ResetPassword)
		}

		// Protected auth routes (require authentication)
		authProtected := api.Group("/auth")
		authProtected.Use(middleware.AuthMiddleware(jwtService))
		{
			authProtected.GET("/profile", authHandler.GetProfile)
			authProtected.POST("/password/change", authHandler.ChangePassword)
			authProtected.POST("/logout", authHandler.Logout)
		}

		// User routes (admin only for management)
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(jwtService))
		users.Use(middleware.AdminMiddleware())
		{
			users.GET("", userHandler.ListUsers)
			users.POST("", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUser)
			users.POST("/seed", userHandler.SeedUsers)
		}

		// Template routes (public read, protected write)
		templates := api.Group("/templates")
		{
			// Public routes
			templates.GET("", templateHandler.ListTemplates)
			templates.GET("/:id", templateHandler.GetTemplate)
			templates.GET("/:id/view", templateHandler.ViewTemplate)
			templates.GET("/:id/thumbnail", templateHandler.GetTemplateThumbnail)
			templates.GET("/:id/variables", templateHandler.GetTemplateVariables)
			templates.GET("/category/:category_id", templateHandler.GetTemplatesByCategory)

			// Protected routes (authenticated users)
			templatesAuth := templates.Group("")
			templatesAuth.Use(middleware.AuthMiddleware(jwtService))
			{
				templatesAuth.POST("/:id/generate", templateHandler.GenerateCustomPDF)
				templatesAuth.GET("/:id/download", templateHandler.DownloadTemplatePDF)
			}

			// Admin only routes
			templatesAdmin := templates.Group("")
			templatesAdmin.Use(middleware.AuthMiddleware(jwtService))
			templatesAdmin.Use(middleware.AdminMiddleware())
			{
				templatesAdmin.POST("", templateHandler.CreateTemplate)
				templatesAdmin.PUT("/:id", templateHandler.UpdateTemplate)
				templatesAdmin.DELETE("/:id", templateHandler.DeleteTemplate)
			}
		}

		// Category routes (public read, admin write)
		categories := api.Group("/categories")
		{
			// Public routes
			categories.GET("", categoryHandler.ListCategories)
			categories.GET("/:id", categoryHandler.GetCategory)

			// Admin only routes
			categoriesAdmin := categories.Group("")
			categoriesAdmin.Use(middleware.AuthMiddleware(jwtService))
			categoriesAdmin.Use(middleware.AdminMiddleware())
			{
				categoriesAdmin.POST("", categoryHandler.CreateCategory)
				categoriesAdmin.PUT("/:id", categoryHandler.UpdateCategory)
				categoriesAdmin.DELETE("/:id", categoryHandler.DeleteCategory)
				categoriesAdmin.POST("/seed", categoryHandler.SeedCategories)
			}
		}

		// Blog routes (public read, admin write)
		blog := api.Group("/blog")
		{
			// Public routes
			blog.GET("", blogHandler.ListBlogPosts)
			blog.GET("/:id", blogHandler.GetBlogPost)
			blog.GET("/category/:category_id", blogHandler.GetBlogPostsByCategory)
			blog.GET("/author/:author_id", blogHandler.GetBlogPostsByAuthor)

			// Admin only routes
			blogAdmin := blog.Group("")
			blogAdmin.Use(middleware.AuthMiddleware(jwtService))
			blogAdmin.Use(middleware.AdminMiddleware())
			{
				blogAdmin.POST("", blogHandler.CreateBlogPost)
				blogAdmin.PUT("/:id", blogHandler.UpdateBlogPost)
				blogAdmin.DELETE("/:id", blogHandler.DeleteBlogPost)
			}
		}

		// Payment routes (authenticated users only)
		payment := api.Group("/payment")
		payment.Use(middleware.AuthMiddleware(jwtService))
		{
			payment.POST("/checkout", paymentHandler.CreateCheckoutSession)
			payment.POST("/intent", paymentHandler.CreatePaymentIntent)
			payment.GET("/success", paymentHandler.GetCheckoutSessionSuccess)
		}

		// Order routes (authenticated users for own orders, admin for all)
		orders := api.Group("/orders")
		orders.Use(middleware.AuthMiddleware(jwtService))
		{
			orders.GET("/user/:user_id", paymentHandler.GetUserOrders) // Users can view their own
			orders.GET("/:id", paymentHandler.GetOrderByID)

			// Admin only
			ordersAdmin := orders.Group("")
			ordersAdmin.Use(middleware.AdminMiddleware())
			{
				ordersAdmin.GET("", paymentHandler.GetAllOrders)
			}
		}
	}

	// Webhook endpoint (outside of API group, no middleware)
	r.POST("/webhook/stripe", paymentHandler.HandleWebhook)

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
