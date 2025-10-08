package health

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewHealthChecker(t *testing.T) {
	hc := NewHealthChecker(nil, "1.0.0")

	assert.NotNil(t, hc)
	assert.Equal(t, "1.0.0", hc.version)
	assert.NotNil(t, hc.checkers)
}

func TestHealthChecker_RegisterChecker(t *testing.T) {
	hc := NewHealthChecker(nil, "1.0.0")

	checker := func(ctx context.Context) ComponentHealth {
		return ComponentHealth{
			Name:   "test",
			Status: HealthStatusHealthy,
		}
	}

	hc.RegisterChecker("test", checker)

	assert.Len(t, hc.checkers, 1)
}

func TestHealthChecker_Check(t *testing.T) {
	hc := NewHealthChecker(nil, "1.0.0")

	hc.RegisterChecker("component1", func(ctx context.Context) ComponentHealth {
		return ComponentHealth{
			Name:      "component1",
			Status:    HealthStatusHealthy,
			Message:   "All good",
			Timestamp: time.Now(),
		}
	})

	hc.RegisterChecker("component2", func(ctx context.Context) ComponentHealth {
		return ComponentHealth{
			Name:      "component2",
			Status:    HealthStatusHealthy,
			Message:   "Working fine",
			Timestamp: time.Now(),
		}
	})

	result := hc.Check(context.Background())

	assert.Equal(t, HealthStatusHealthy, result.Status)
	assert.Equal(t, "1.0.0", result.Version)
	assert.Len(t, result.Components, 2)
	assert.NotEmpty(t, result.Uptime)
}

func TestHealthChecker_Check_Degraded(t *testing.T) {
	hc := NewHealthChecker(nil, "1.0.0")

	hc.RegisterChecker("healthy", func(ctx context.Context) ComponentHealth {
		return ComponentHealth{
			Name:   "healthy",
			Status: HealthStatusHealthy,
		}
	})

	hc.RegisterChecker("degraded", func(ctx context.Context) ComponentHealth {
		return ComponentHealth{
			Name:   "degraded",
			Status: HealthStatusDegraded,
		}
	})

	result := hc.Check(context.Background())

	assert.Equal(t, HealthStatusDegraded, result.Status)
	assert.Len(t, result.Components, 2)
}

func TestHealthChecker_Check_Unhealthy(t *testing.T) {
	hc := NewHealthChecker(nil, "1.0.0")

	hc.RegisterChecker("healthy", func(ctx context.Context) ComponentHealth {
		return ComponentHealth{
			Name:   "healthy",
			Status: HealthStatusHealthy,
		}
	})

	hc.RegisterChecker("unhealthy", func(ctx context.Context) ComponentHealth {
		return ComponentHealth{
			Name:   "unhealthy",
			Status: HealthStatusUnhealthy,
		}
	})

	result := hc.Check(context.Background())

	assert.Equal(t, HealthStatusUnhealthy, result.Status)
	assert.Len(t, result.Components, 2)
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "seconds only",
			duration: 45 * time.Second,
			expected: "45s",
		},
		{
			name:     "minutes and seconds",
			duration: 5*time.Minute + 30*time.Second,
			expected: "5m 30s",
		},
		{
			name:     "hours, minutes, seconds",
			duration: 2*time.Hour + 15*time.Minute + 45*time.Second,
			expected: "2h 15m 45s",
		},
		{
			name:     "days, hours, minutes, seconds",
			duration: 3*24*time.Hour + 5*time.Hour + 30*time.Minute + 20*time.Second,
			expected: "3d 5h 30m 20s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDuration(tt.duration)
			assert.Equal(t, tt.expected, result)
		})
	}
}
