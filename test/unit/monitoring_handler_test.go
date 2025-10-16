package unit

import (
	"base-gin/internal/app/handlers"
	"base-gin/internal/domain/models"
	"base-gin/internal/pkg/health"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHealthChecker is a mock implementation of HealthChecker
type MockHealthChecker struct {
	mock.Mock
}

func (m *MockHealthChecker) Check(ctx context.Context) health.HealthCheck {
	args := m.Called(ctx)
	return args.Get(0).(health.HealthCheck)
}

func TestMonitoringHandler_GetHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockSetup      func(*MockHealthChecker)
		expectedStatus int
	}{
		{
			name: "healthy status",
			mockSetup: func(healthChecker *MockHealthChecker) {
				healthCheck := health.HealthCheck{
					Status:    health.HealthStatusHealthy,
					Timestamp: time.Now(),
					Components: []health.ComponentHealth{
						{
							Name:      "database",
							Status:    health.HealthStatusHealthy,
							Message:   "Database connection is healthy",
							Timestamp: time.Now(),
							Duration:  time.Since(time.Now()).String(),
						},
						{
							Name:      "redis",
							Status:    health.HealthStatusHealthy,
							Message:   "Redis connection is healthy",
							Timestamp: time.Now(),
							Duration:  time.Since(time.Now()).String(),
						},
					},
				}
				healthChecker.On("Check", mock.Anything).Return(healthCheck)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "degraded status",
			mockSetup: func(healthChecker *MockHealthChecker) {
				healthCheck := health.HealthCheck{
					Status:    health.HealthStatusDegraded,
					Timestamp: time.Now(),
					Components: []health.ComponentHealth{
						{
							Name:      "database",
							Status:    health.HealthStatusHealthy,
							Message:   "Database connection is healthy",
							Timestamp: time.Now(),
							Duration:  time.Since(time.Now()).String(),
						},
					},
				}
				healthChecker.On("Check", mock.Anything).Return(healthCheck)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "unhealthy status",
			mockSetup: func(healthChecker *MockHealthChecker) {
				healthCheck := health.HealthCheck{
					Status:    health.HealthStatusUnhealthy,
					Timestamp: time.Now(),
					Components: []health.ComponentHealth{
						{
							Name:      "database",
							Status:    health.HealthStatusUnhealthy,
							Message:   "Database connection failed",
							Timestamp: time.Now(),
							Duration:  time.Since(time.Now()).String(),
						},
						{
							Name:      "redis",
							Status:    health.HealthStatusUnhealthy,
							Message:   "Redis connection failed",
							Timestamp: time.Now(),
							Duration:  time.Since(time.Now()).String(),
						},
					},
				}
				healthChecker.On("Check", mock.Anything).Return(healthCheck)
			},
			expectedStatus: http.StatusServiceUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockHealthChecker := new(MockHealthChecker)
			tt.mockSetup(mockHealthChecker)

			// Create handler
			healthChecker := health.NewHealthChecker(nil, "")
			handler := handlers.NewMonitoringHandler(healthChecker)

			// Setup router
			router := gin.New()
			router.GET("/health", handler.GetHealthCheck)

			// Create request
			req, _ := http.NewRequest("GET", "/health", nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestMonitoringHandler_GetMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "successful get metrics",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create handler
			handler := handlers.NewMonitoringHandler(nil)

			// Setup router
			router := gin.New()
			router.GET("/metrics", handler.GetMetrics)

			// Create request
			req, _ := http.NewRequest("GET", "/metrics", nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response models.Metrics
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.NotNil(t, response)
		})
	}
}

func TestMonitoringHandler_ResetMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "successful reset metrics",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create handler
			handler := handlers.NewMonitoringHandler(nil)

			// Setup router
			router := gin.New()
			router.POST("/metrics/reset", handler.ResetMetrics)

			// Create request
			req, _ := http.NewRequest("POST", "/metrics/reset", nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Metrics reset successfully", response["message"])
		})
	}
}
