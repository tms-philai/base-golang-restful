package main

import (
	"log"
	"net/http"

	"base-golang-restful-app/config"
	"base-golang-restful-app/handlers"
	"base-golang-restful-app/middleware"
	"base-golang-restful-app/models"
	"base-golang-restful-app/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "base-golang-restful-app/docs" // Import generated docs
)

// @title Base Golang RESTful API with Authentication
// @version 2.0
// @description A comprehensive RESTful API built with Go, Gin framework, JWT authentication, and full CRUD operations
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize services
	userService := services.NewUserService()
	productService := services.NewProductService()
	jwtService := services.NewJWTService(&cfg.JWT)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService, jwtService)
	userHandler := handlers.NewUserHandler(userService)
	productHandler := handlers.NewProductHandler(productService)

	// Create Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "base-golang-restful-app",
			"version": "2.0.0",
		})
	})

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Authentication routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			
			// Protected auth routes
			authProtected := auth.Group("")
			authProtected.Use(middleware.AuthMiddleware(jwtService, userService))
			{
				authProtected.GET("/profile", authHandler.GetProfile)
				authProtected.POST("/change-password", authHandler.ChangePassword)
			}
		}

		// User routes
		users := v1.Group("/users")
		{
			// Public user routes (with optional auth)
			users.GET("/:id", middleware.OptionalAuthMiddleware(jwtService, userService), userHandler.GetUser)
			
			// Protected user routes
			usersProtected := users.Group("")
			usersProtected.Use(middleware.AuthMiddleware(jwtService, userService))
			{
				// Admin only routes
				usersProtected.GET("", middleware.RequireRole("admin"), userHandler.ListUsers)
				usersProtected.POST("", middleware.RequireRole("admin"), userHandler.CreateUser)
				usersProtected.DELETE("/:id", middleware.RequireRole("admin"), userHandler.DeleteUser)
				
				// User can update their own profile, admin can update any
				usersProtected.PUT("/:id", userHandler.UpdateUser)
			}
		}

		// Product routes
		products := v1.Group("/products")
		{
			// Public product routes
			products.GET("", productHandler.ListProducts)
			products.GET("/categories", productHandler.GetCategories)
			products.GET("/:id", productHandler.GetProduct)
			
			// Protected product routes
			productsProtected := products.Group("")
			productsProtected.Use(middleware.AuthMiddleware(jwtService, userService))
			{
				productsProtected.POST("", productHandler.CreateProduct)
				productsProtected.PUT("/:id", productHandler.UpdateProduct)
				productsProtected.DELETE("/:id", productHandler.DeleteProduct)
				productsProtected.PATCH("/:id/stock", productHandler.UpdateStock)
			}
		}
	}

	// Create a default admin user if none exists
	createDefaultAdmin(userService)

	// Start server
	log.Printf("Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Swagger documentation available at: http://localhost:%s/swagger/index.html", cfg.Server.Port)
	log.Printf("API base URL: http://localhost:%s/api/v1", cfg.Server.Port)
	
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// createDefaultAdmin creates a default admin user for testing purposes
func createDefaultAdmin(userService *services.UserService) {
	// Check if any admin user exists
	users, _, err := userService.List(1, 100)
	if err != nil {
		log.Println("Warning: Could not check for existing admin users")
		return
	}

	hasAdmin := false
	for _, user := range users {
		if user.Role == "admin" {
			hasAdmin = true
			break
		}
	}

	if !hasAdmin {
		adminReq := models.UserCreateRequest{
			Username:  "admin",
			Email:     "admin@example.com",
			Password:  "admin123",
			FirstName: "System",
			LastName:  "Administrator",
		}

		admin, err := userService.Create(adminReq)
		if err != nil {
			log.Printf("Warning: Could not create default admin user: %v", err)
			return
		}

		// Set admin role
		admin.Role = "admin"
		log.Println("Default admin user created:")
		log.Println("  Username: admin")
		log.Println("  Password: admin123")
		log.Println("  Email: admin@example.com")
		log.Println("Please change the default password after first login!")
	}
}
