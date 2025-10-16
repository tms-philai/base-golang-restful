package routes

import (
	"base-gin/internal/app/handlers"
	"base-gin/internal/app/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouterConfig struct {
	AuthHandler    *handlers.AuthHandler
	AuthMiddleware *middleware.AuthMiddleware
	UserHandler    *handlers.UserHandler
	ProductHandler *handlers.ProductHandler
	EmailHandler   *handlers.EmailHandler
	FileHandler    *handlers.FileHandler
}

func SetupRouter(cfg RouterConfig) *gin.Engine {
	r := gin.Default()

	middleware.SetupGlobalMiddlewares(r)
	// Setup public routes
	setupPublicRoutes(r)

	// Setup API v1 routes
	v1 := r.Group("/api/v1")
	{
		setupAuthRoutes(v1, cfg)
		setupUserRoutes(v1, cfg)
		setupProductRoutes(v1, cfg)
		setupEmailRoutes(v1, cfg)
		setupFileRoutes(v1, cfg)
	}

	return r
}

func setupPublicRoutes(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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
		users.GET("/:id", cfg.AuthMiddleware.OptionalAuthenticate(), cfg.UserHandler.GetUser)

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

func setupFileRoutes(v1 *gin.RouterGroup, cfg RouterConfig) {
	files := v1.Group("/files")
	{
		// Public routes
		files.GET("/:id", cfg.FileHandler.GetFile)
		files.GET("/:id/download", cfg.FileHandler.DownloadFile)
		files.GET("/:id/url", cfg.FileHandler.GetFileURL)

		// Protected routes
		filesProtected := files.Group("")
		filesProtected.Use(cfg.AuthMiddleware.Authenticate())
		{
			filesProtected.POST("/upload", cfg.FileHandler.UploadFile)
			filesProtected.POST("/upload/multiple", cfg.FileHandler.UploadMultipleFiles)
			filesProtected.GET("", cfg.FileHandler.ListFiles)
			filesProtected.DELETE("/:id", cfg.FileHandler.DeleteFile)
		}
	}
}
