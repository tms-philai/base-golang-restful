package handlers

import (
	"net/http"
	"runtime"
	"time"

	"base-golang-restful-app/middleware"
	"base-golang-restful-app/models"

	"github.com/gin-gonic/gin"
)

type APIHandler struct{}

func NewAPIHandler() *APIHandler {
	return &APIHandler{}
}

// GetAPIInfo godoc
// @Summary Get API information
// @Description Get comprehensive API information including version, status, and capabilities
// @Tags API Info
// @Produce json
// @Success 200 {object} models.APIInfoResponse
// @Router /info [get]
func (h *APIHandler) GetAPIInfo(c *gin.Context) {
	version := middleware.GetAPIVersion(c)

	c.JSON(http.StatusOK, models.APIInfoResponse{
		Name:        "Base Golang RESTful API",
		Version:     version,
		Description: "A comprehensive RESTful API with authentication, RBAC, email, and notifications",
		Status:      "operational",
		Timestamp:   time.Now(),
		Features: []string{
			"JWT Authentication",
			"Role-Based Access Control (RBAC)",
			"Email Service",
			"Notification System",
			"File Upload & Storage",
			"Caching with Redis",
			"API Versioning",
			"Swagger Documentation",
			"Internationalization (i18n)",
		},
		Endpoints: models.APIEndpoints{
			Auth:          "/api/v1/auth",
			Users:         "/api/v1/users",
			Products:      "/api/v1/products",
			Email:         "/api/v1/email",
			Notifications: "/api/v1/notifications",
			Documentation: "/swagger/index.html",
		},
	})
}

// GetAPIVersion godoc
// @Summary Get current API version
// @Description Get the current API version being used (from header, query, or path)
// @Tags API Info
// @Produce json
// @Param API-Version header string false "API Version in header"
// @Param version query string false "API Version in query"
// @Success 200 {object} models.APIVersionResponse
// @Router /version [get]
func (h *APIHandler) GetAPIVersion(c *gin.Context) {
	version := middleware.GetAPIVersion(c)

	extractedFrom := "default"
	if c.GetHeader(middleware.VersionHeader) != "" {
		extractedFrom = "header"
	} else if c.Query(middleware.VersionParam) != "" {
		extractedFrom = "query"
	} else if version != middleware.DefaultAPIVersion {
		extractedFrom = "path"
	}

	c.JSON(http.StatusOK, models.APIVersionResponse{
		Version:       version,
		ExtractedFrom: extractedFrom,
		SupportedVersions: []string{
			"v1",
			"v2",
		},
		DeprecatedVersions: []string{},
	})
}

// GetHealthDetailed godoc
// @Summary Get detailed health check
// @Description Get detailed health information including system metrics and service status
// @Tags API Info
// @Produce json
// @Success 200 {object} models.HealthDetailedResponse
// @Router /health/detailed [get]
func (h *APIHandler) GetHealthDetailed(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	c.JSON(http.StatusOK, models.HealthDetailedResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   middleware.GetAPIVersion(c),
		Uptime:    time.Since(startTime).String(),
		System: models.SystemMetrics{
			GoVersion:    runtime.Version(),
			NumGoroutine: runtime.NumGoroutine(),
			NumCPU:       runtime.NumCPU(),
			MemoryAlloc:  m.Alloc,
			MemoryTotal:  m.TotalAlloc,
			MemorySys:    m.Sys,
			NumGC:        m.NumGC,
		},
		Services: models.ServiceStatus{
			Database:      true,
			Redis:         false,
			Email:         true,
			Notifications: true,
		},
	})
}

// V1Example godoc
// @Summary API V1 example endpoint
// @Description Example endpoint that only works with API v1
// @Tags API Versioning
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /api/v1/api-version/v1-only [get]
func (h *APIHandler) V1Example(c *gin.Context) {
	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "This endpoint works with API v1 only",
		Data: gin.H{
			"version": "v1",
			"feature": "Legacy authentication system",
		},
	})
}

// V2Example godoc
// @Summary API V2 example endpoint
// @Description Example endpoint that only works with API v2
// @Tags API Versioning
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /api/v1/api-version/v2-only [get]
func (h *APIHandler) V2Example(c *gin.Context) {
	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "This endpoint works with API v2 only",
		Data: gin.H{
			"version": "v2",
			"feature": "Enhanced security with OAuth2",
		},
	})
}

// GetServerTime godoc
// @Summary Get server time
// @Description Get current server time in various formats
// @Tags API Info
// @Produce json
// @Success 200 {object} models.ServerTimeResponse
// @Router /time [get]
func (h *APIHandler) GetServerTime(c *gin.Context) {
	now := time.Now()

	c.JSON(http.StatusOK, models.ServerTimeResponse{
		UTC:       now.UTC(),
		Local:     now,
		Unix:      now.Unix(),
		UnixMilli: now.UnixMilli(),
		Timezone:  now.Location().String(),
		Formatted: now.Format(time.RFC3339),
	})
}

var startTime = time.Now()
