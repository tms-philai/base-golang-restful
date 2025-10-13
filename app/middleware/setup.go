package middleware

import (
	"base-golang-restful-app/metrics"

	"github.com/gin-gonic/gin"
)

func SetupGlobalMiddlewares(r *gin.Engine) {
	// Metrics collection
	r.Use(metrics.MetricsMiddleware())

	// API Versioning
	r.Use(APIVersioning())

	// CORS
	r.Use(CORS())
}
