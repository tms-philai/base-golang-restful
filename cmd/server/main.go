package main

import (
	_ "base-gin/docs"
	"base-gin/internal/app/handlers"
	"base-gin/internal/app/middleware"
	"base-gin/internal/app/routes"
	"base-gin/internal/domain/repository"
	"base-gin/internal/domain/services"
	"base-gin/internal/pkg/auth"
	"base-gin/internal/pkg/config"
	"base-gin/internal/pkg/database"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

// @host localhost:8001
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	cfg, err := config.Load("./configs/.env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

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

	if err := database.ConnectWithRetry(dbConfig, 5, 2); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run auto migrations
	log.Println("Running database migrations...")
	if err := database.AutoMigrate(database.GetDB()); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Seed default roles and permissions
	log.Println("Seeding default data...")
	if err := database.SeedDefaultRoles(database.GetDB()); err != nil {
		log.Printf("Warning: Failed to seed default roles: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(database.GetDB())
	roleRepo := repository.NewRoleRepository(database.GetDB())
	permissionRepo := repository.NewPermissionRepository(database.GetDB())

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(auth.JWTConfig{
		SecretKey:            cfg.JWT.SecretKey,
		AccessTokenDuration:  cfg.JWT.AccessTokenDuration,
		RefreshTokenDuration: cfg.JWT.RefreshTokenDuration,
		Issuer:               "base-golang-restful",
	})

	// Initialize services
	userService := services.NewUserService(userRepo)
	roleService := services.NewRoleService(database.GetDB(), roleRepo, userRepo, permissionRepo)

	// Initialize middleware and handlers
	authMiddleware := middleware.NewAuthMiddleware(jwtManager, userService)
	authHandler := handlers.NewAuthHandler(userService, roleService, jwtManager)

	r := routes.SetupRouter(routes.RouterConfig{
		AuthHandler:    authHandler,
		AuthMiddleware: authMiddleware,
	})
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}
	go func() {
		log.Printf("🚀 Server is running on http://%s", serverAddr)
		log.Printf("📚 Swagger documentation: http://%s/swagger/index.html", serverAddr)
		log.Printf("🔐 Auth endpoints: http://%s/api/v1/auth/*", serverAddr)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Failed to gracefully shutdown server: %v", err)
	}
	log.Println("Server stopped gracefully")

}
