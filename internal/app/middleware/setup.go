package middleware

import (
	"base-gin/internal/domain/models"

	"github.com/gin-gonic/gin"
)

func SetupGlobalMiddlewares(r *gin.Engine) {

	r.Use(models.MetricsMiddleware())
	// API Versioning
	r.Use(APIVersioning())

	// CORS
	r.Use(CORS())
}
