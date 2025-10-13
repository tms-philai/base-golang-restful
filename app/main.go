package main

import (
	"log"
	"time"

	"base-golang-restful-app/auth"
	"base-golang-restful-app/config"
	"base-golang-restful-app/database"
	"base-golang-restful-app/email"
	"base-golang-restful-app/handlers"
	"base-golang-restful-app/health"
	"base-golang-restful-app/i18n"
	"base-golang-restful-app/middleware"
	"base-golang-restful-app/models"
	"base-golang-restful-app/repository"
	"base-golang-restful-app/routes"
	"base-golang-restful-app/services"

	_ "base-golang-restful-app/docs"
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
// @BasePath /
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

	// Initialize i18n
	log.Println("Initializing i18n...")
	if err := i18n.InitI18n(i18n.I18nConfig{
		DefaultLanguage: "en",
		LocalesPath:     "./locales",
		SupportedLangs:  []string{"en", "vi", "ja"},
	}); err != nil {
		log.Fatal("Failed to initialize i18n:", err)
	}

	// Initialize database connection
	log.Println("Connecting to database...")
	dbConfig := database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.Name,
		SSLMode:  cfg.Database.SSLMode,
		TimeZone: cfg.Database.TimeZone,
	}

	if err := database.ConnectWithRetry(dbConfig, 5, 2*time.Second); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Run database migrations (only if enabled)
	if cfg.Database.AutoMigrate {
		log.Println("Running database migrations...")
		migrator := database.NewMigrator(database.GetDB())
		if err := migrator.AutoMigrate(
			&models.User{},
			&models.Role{},
			&models.Permission{},
			&models.Product{},
			&models.RefreshToken{},
			&models.File{},
		); err != nil {
			log.Fatal("Failed to run migrations:", err)
		}
		log.Println("Database migrations completed successfully")
	} else {
		log.Println("Database auto-migration is disabled. Set DB_AUTO_MIGRATE=true to enable.")
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(database.GetDB())
	productRepo := repository.NewProductRepository(database.GetDB())

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(auth.JWTConfig{
		SecretKey:            cfg.JWT.SecretKey,
		AccessTokenDuration:  cfg.JWT.AccessTokenDuration,
		RefreshTokenDuration: cfg.JWT.RefreshTokenDuration,
		Issuer:               "base-golang-restful",
	})

	// Initialize services with repositories
	userService := services.NewUserService(userRepo)
	productService := services.NewProductService(productRepo)

	// Initialize email service
	var emailService *email.EmailService
	var notificationService *email.NotificationService
	var emailHandler *handlers.EmailHandler
	var notificationHandler *handlers.NotificationHandler
	var inAppChannel *email.InAppNotificationChannel

	if cfg.Email.Enabled {
		log.Println("Initializing email service...")
		emailClient := email.NewEmailClient(email.SMTPConfig{
			Host:     cfg.Email.SMTPHost,
			Port:     cfg.Email.SMTPPort,
			Username: cfg.Email.SMTPUser,
			Password: cfg.Email.SMTPPass,
			From:     cfg.Email.From,
			UseTLS:   true,
		})

		emailService = email.NewEmailService(email.EmailServiceConfig{
			Client:    emailClient,
			Workers:   5,
			QueueSize: 100,
		})

		notificationService = email.NewNotificationService(emailService)

		emailChannel := email.NewEmailNotificationChannel(emailService, cfg.Email.From)
		notificationService.RegisterChannel(emailChannel)

		inAppChannel = email.NewInAppNotificationChannel()
		notificationService.RegisterChannel(inAppChannel)

		emailHandler = handlers.NewEmailHandler(emailService)
		notificationHandler = handlers.NewNotificationHandler(notificationService, inAppChannel)

		log.Println("Email and Notification services initialized successfully")
	} else {
		log.Println("Email service is disabled. Set EMAIL_ENABLED=true to enable.")
	}

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	// Initialize health checker
	healthChecker := health.NewHealthChecker(database.GetDB(), "2.0.0")
	healthChecker.RegisterChecker("database", healthChecker.DatabaseChecker())

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService, jwtManager)
	userHandler := handlers.NewUserHandler(userService)
	productHandler := handlers.NewProductHandler(productService)
	apiHandler := handlers.NewAPIHandler()
	monitoringHandler := handlers.NewMonitoringHandler(healthChecker)

	// Setup router with all routes
	r := routes.SetupRouter(routes.RouterConfig{
		AuthHandler:         authHandler,
		UserHandler:         userHandler,
		ProductHandler:      productHandler,
		EmailHandler:        emailHandler,
		NotificationHandler: notificationHandler,
		APIHandler:          apiHandler,
		MonitoringHandler:   monitoringHandler,
		AuthMiddleware:      authMiddleware,
		EmailEnabled:        cfg.Email.Enabled && emailHandler != nil,
	})

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
	// Check if admin user already exists
	_, err := userService.GetByEmail("admin@example.com")
	if err == nil {
		// Admin already exists
		log.Println("Default admin user already exists")
		return
	}

	// Create admin user
	adminReq := models.UserCreateRequest{
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

	log.Println("Default admin user created successfully!")
	log.Println("   Email: admin@example.com")
	log.Println("   Password: admin123")
	log.Printf("   User ID: %s", admin.ID)
	log.Println("IMPORTANT: Please change the default password after first login!")
}
