package handlers

import (
	"net/http"

	"base-golang-restful-app/health"
	"base-golang-restful-app/metrics"

	"github.com/gin-gonic/gin"
)

type MonitoringHandler struct {
	healthChecker *health.HealthChecker
}

func NewMonitoringHandler(healthChecker *health.HealthChecker) *MonitoringHandler {
	return &MonitoringHandler{
		healthChecker: healthChecker,
	}
}

// GetHealthCheck godoc
// @Summary Get comprehensive health check
// @Description Get detailed health status of all system components
// @Tags Monitoring
// @Produce json
// @Success 200 {object} health.HealthCheck
// @Router /health [get]
func (h *MonitoringHandler) GetHealthCheck(c *gin.Context) {
	healthStatus := h.healthChecker.Check(c.Request.Context())

	statusCode := http.StatusOK
	if healthStatus.Status == health.HealthStatusUnhealthy {
		statusCode = http.StatusServiceUnavailable
	} else if healthStatus.Status == health.HealthStatusDegraded {
		statusCode = http.StatusOK
	}

	c.JSON(statusCode, healthStatus)
}

// GetMetrics godoc
// @Summary Get application metrics
// @Description Get detailed metrics about API performance and usage
// @Tags Monitoring
// @Produce json
// @Success 200 {object} metrics.Metrics
// @Router /metrics [get]
func (h *MonitoringHandler) GetMetrics(c *gin.Context) {
	metricsData := metrics.GetMetrics()
	c.JSON(http.StatusOK, metricsData)
}

// ResetMetrics godoc
// @Summary Reset metrics
// @Description Reset all collected metrics (admin only)
// @Tags Monitoring
// @Success 200 {object} models.SuccessResponse
// @Security BearerAuth
// @Router /metrics/reset [post]
func (h *MonitoringHandler) ResetMetrics(c *gin.Context) {
	metrics.ResetMetrics()
	c.JSON(http.StatusOK, gin.H{
		"message": "Metrics reset successfully",
	})
}
