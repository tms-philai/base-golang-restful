package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"base-golang-restful/app/health"
	"base-golang-restful/app/i18n"
	"base-golang-restful/app/logger"
	"base-golang-restful/app/metrics"
	"base-golang-restful/app/middleware"

	"github.com/gin-gonic/gin"
)

// @title Base Golang RESTful API
// @version 1.0
// @description A comprehensive RESTful API built with Go and Gin framework
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Initialize logger
	err := logger.InitLogger(logger.LogConfig{
		Level:      "info",
		OutputPath: "./logs/app.log",
	})
	if err != nil {
		log.Printf("Failed to initialize logger: %v", err)
	}

	// Initialize i18n
	err = i18n.InitI18n(i18n.I18nConfig{
		DefaultLanguage: "en",
		LocalesPath:     "./app/locales",
		SupportedLangs:  []string{"en", "vi"},
	})
	if err != nil {
		log.Printf("Failed to initialize i18n: %v", err)
	}

	// Initialize health checker
	healthChecker := health.NewHealthChecker(nil, "1.0.0")
	
	// Register health checkers
	healthChecker.RegisterChecker("api", func(ctx context.Context) health.ComponentHealth {
		return health.ComponentHealth{
			Name:      "api",
			Status:    health.HealthStatusHealthy,
			Message:   "API is running",
			Timestamp: time.Now(),
		}
	})

	// Create Gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Add middleware
	r.Use(gin.Recovery())
	r.Use(logger.RequestLogger())
	r.Use(i18n.LanguageMiddleware())
	r.Use(middleware.APIVersioning())
	r.Use(metrics.MetricsMiddleware())
	r.Use(middleware.ErrorHandler())

	// CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, API-Version")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	})

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Base Golang RESTful API",
			"version": "1.0.0",
			"status":  "running",
			"docs":    "/swagger/index.html",
		})
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		result := healthChecker.Check(c.Request.Context())
		
		statusCode := http.StatusOK
		if result.Status == health.HealthStatusUnhealthy {
			statusCode = http.StatusServiceUnavailable
		} else if result.Status == health.HealthStatusDegraded {
			statusCode = http.StatusOK
		}
		
		c.JSON(statusCode, result)
	})

	// Metrics endpoint
	r.GET("/metrics", func(c *gin.Context) {
		metricsData := metrics.GetMetrics()
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    metricsData,
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Test endpoint
		v1.GET("/test", func(c *gin.Context) {
			lang := i18n.GetLanguage(c)
			version := middleware.GetAPIVersion(c)
			
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": i18n.T(c, "common.success", nil),
				"data": gin.H{
					"language":    lang,
					"api_version": version,
					"timestamp":   time.Now(),
				},
			})
		})

		// Echo endpoint for testing
		v1.POST("/echo", func(c *gin.Context) {
			var body map[string]interface{}
			if err := c.ShouldBindJSON(&body); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"error": gin.H{
						"code":    "INVALID_REQUEST",
						"message": "Invalid JSON body",
					},
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    body,
			})
		})

		// Localization test endpoint
		v1.GET("/i18n-test", func(c *gin.Context) {
			lang := i18n.GetLanguage(c)
			
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data": gin.H{
					"language": lang,
					"messages": gin.H{
						"success":       i18n.Translate(lang, "common.success", nil),
						"error":         i18n.Translate(lang, "common.error", nil),
						"created":       i18n.Translate(lang, "common.created", nil),
						"not_found":     i18n.Translate(lang, "common.not_found", nil),
						"unauthorized":  i18n.Translate(lang, "common.unauthorized", nil),
					},
				},
			})
		})
	}

	// 404 handler
	r.NoRoute(middleware.HandleNotFound())

	// 405 handler
	r.NoMethod(middleware.HandleMethodNotAllowed())

	// Start server
	port := "8080"
	logger.Info().
		Str("port", port).
		Str("version", "1.0.0").
		Msg("Starting server")

	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║   🚀 Base Golang RESTful API                              ║")
	fmt.Println("╠════════════════════════════════════════════════════════════╣")
	fmt.Printf("║   📍 Server:     http://localhost:%s                     ║\n", port)
	fmt.Printf("║   🏥 Health:     http://localhost:%s/health              ║\n", port)
	fmt.Printf("║   📊 Metrics:    http://localhost:%s/metrics             ║\n", port)
	fmt.Printf("║   🧪 Test API:   http://localhost:%s/api/v1/test         ║\n", port)
	fmt.Printf("║   🌍 i18n Test:  http://localhost:%s/api/v1/i18n-test    ║\n", port)
	fmt.Println("╠════════════════════════════════════════════════════════════╣")
	fmt.Println("║   Features:                                                ║")
	fmt.Println("║   ✅ Structured Logging (Zerolog)                          ║")
	fmt.Println("║   ✅ Internationalization (EN/VI)                          ║")
	fmt.Println("║   ✅ API Versioning                                        ║")
	fmt.Println("║   ✅ Health Checks                                         ║")
	fmt.Println("║   ✅ Metrics Collection                                    ║")
	fmt.Println("║   ✅ Error Handling                                        ║")
	fmt.Println("║   ✅ CORS Support                                          ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")

	if err := r.Run(":" + port); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
