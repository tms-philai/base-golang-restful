package routes

import (
	"base-golang-restful-app/handlers"
	"base-golang-restful-app/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouterConfig struct {
	AuthHandler         *handlers.AuthHandler
	UserHandler         *handlers.UserHandler
	ProductHandler      *handlers.ProductHandler
	EmailHandler        *handlers.EmailHandler
	NotificationHandler *handlers.NotificationHandler
	APIHandler          *handlers.APIHandler
	MonitoringHandler   *handlers.MonitoringHandler
	AuthMiddleware      *middleware.AuthMiddleware
	EmailEnabled        bool
}

func SetupRouter(cfg RouterConfig) *gin.Engine {
	r := gin.Default()

	// Apply global middlewares
	middleware.SetupGlobalMiddlewares(r)

	// Setup public routes
	setupPublicRoutes(r, cfg)

	// Setup API routes
	setupAPIRoutes(r, cfg)

	return r
}

func setupPublicRoutes(r *gin.Engine, cfg RouterConfig) {
	// Health check & Monitoring
	r.GET("/health", cfg.MonitoringHandler.GetHealthCheck)
	r.GET("/metrics", cfg.MonitoringHandler.GetMetrics)
	r.POST("/metrics/reset", cfg.AuthMiddleware.Authenticate(), middleware.RequireRole("admin"), cfg.MonitoringHandler.ResetMetrics)

	// API Info
	r.GET("/info", cfg.APIHandler.GetAPIInfo)
	r.GET("/version", cfg.APIHandler.GetAPIVersion)
	r.GET("/time", cfg.APIHandler.GetServerTime)
	r.GET("/health/detailed", cfg.APIHandler.GetHealthDetailed)

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func setupAPIRoutes(r *gin.Engine, cfg RouterConfig) {
	// API Versioning demo
	setupVersioningRoutes(r, cfg)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		setupAuthRoutes(v1, cfg)
		setupUserRoutes(v1, cfg)
		setupProductRoutes(v1, cfg)

		if cfg.EmailEnabled {
			setupEmailRoutes(v1, cfg)
			setupNotificationRoutes(v1, cfg)
		}
	}
}

func setupVersioningRoutes(r *gin.Engine, cfg RouterConfig) {
	apiVersionDemo := r.Group("/api/v1/api-version")
	{
		apiVersionDemo.GET("/v1-only", middleware.RequireVersion("v1"), cfg.APIHandler.V1Example)
		apiVersionDemo.GET("/v2-only", middleware.RequireVersion("v2"), cfg.APIHandler.V2Example)
	}
}

func setupAuthRoutes(v1 *gin.RouterGroup, cfg RouterConfig) {
	auth := v1.Group("/auth")
	{
		// Public routes
		auth.POST("/register", cfg.AuthHandler.Register)
		auth.POST("/login", cfg.AuthHandler.Login)
		auth.POST("/refresh", cfg.AuthHandler.RefreshToken)

		// Protected routes
		authProtected := auth.Group("")
		authProtected.Use(cfg.AuthMiddleware.Authenticate())
		{
			authProtected.GET("/profile", cfg.AuthHandler.GetProfile)
			authProtected.POST("/change-password", cfg.AuthHandler.ChangePassword)
		}
	}
}

func setupUserRoutes(v1 *gin.RouterGroup, cfg RouterConfig) {
	users := v1.Group("/users")
	{
		// Public routes (with optional auth)
		users.GET("/:id", cfg.AuthMiddleware.OptionalAuthenticate(), cfg.UserHandler.GetUser)

		// Protected routes
		usersProtected := users.Group("")
		usersProtected.Use(cfg.AuthMiddleware.Authenticate())
		{
			// Admin only routes
			usersProtected.GET("", middleware.RequireRole("admin"), cfg.UserHandler.ListUsers)
			usersProtected.POST("", middleware.RequireRole("admin"), cfg.UserHandler.CreateUser)
			usersProtected.DELETE("/:id", middleware.RequireRole("admin"), cfg.UserHandler.DeleteUser)

			// User can update their own profile, admin can update any
			usersProtected.PUT("/:id", cfg.UserHandler.UpdateUser)
		}
	}
}

func setupProductRoutes(v1 *gin.RouterGroup, cfg RouterConfig) {
	products := v1.Group("/products")
	{
		// Public routes
		products.GET("", cfg.ProductHandler.ListProducts)
		products.GET("/categories", cfg.ProductHandler.GetCategories)
		products.GET("/:id", cfg.ProductHandler.GetProduct)

		// Protected routes
		productsProtected := products.Group("")
		productsProtected.Use(cfg.AuthMiddleware.Authenticate())
		{
			productsProtected.POST("", cfg.ProductHandler.CreateProduct)
			productsProtected.PUT("/:id", cfg.ProductHandler.UpdateProduct)
			productsProtected.DELETE("/:id", cfg.ProductHandler.DeleteProduct)
			productsProtected.PATCH("/:id/stock", cfg.ProductHandler.UpdateStock)
		}
	}
}

func setupEmailRoutes(v1 *gin.RouterGroup, cfg RouterConfig) {
	emailRoutes := v1.Group("/email")
	emailRoutes.Use(cfg.AuthMiddleware.Authenticate())
	{
		emailRoutes.POST("/send", cfg.EmailHandler.SendEmail)
		emailRoutes.POST("/send-sync", cfg.EmailHandler.SendEmailSync)
		emailRoutes.POST("/send-bulk", cfg.EmailHandler.SendBulkEmail)
		emailRoutes.GET("/status", cfg.EmailHandler.GetEmailStatus)
		emailRoutes.GET("/test-connection", cfg.EmailHandler.TestEmailConnection)
	}
}

func setupNotificationRoutes(v1 *gin.RouterGroup, cfg RouterConfig) {
	notificationRoutes := v1.Group("/notifications")
	notificationRoutes.Use(cfg.AuthMiddleware.Authenticate())
	{
		notificationRoutes.POST("/send", cfg.NotificationHandler.SendNotification)
		notificationRoutes.POST("/send-bulk", cfg.NotificationHandler.SendBulkNotification)
		notificationRoutes.GET("/user/:user_id", cfg.NotificationHandler.GetNotifications)
		notificationRoutes.POST("/mark-read", cfg.NotificationHandler.MarkAsRead)
		notificationRoutes.GET("/status", cfg.NotificationHandler.GetServiceStatus)
	}
}
