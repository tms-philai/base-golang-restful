package health

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
)

type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
)

type ComponentHealth struct {
	Name      string       `json:"name"`
	Status    HealthStatus `json:"status"`
	Message   string       `json:"message,omitempty"`
	Timestamp time.Time    `json:"timestamp"`
	Duration  string       `json:"duration,omitempty"`
}

type HealthCheck struct {
	Status     HealthStatus      `json:"status"`
	Version    string            `json:"version"`
	Uptime     string            `json:"uptime"`
	Timestamp  time.Time         `json:"timestamp"`
	Components []ComponentHealth `json:"components"`
}

type HealthChecker struct {
	db        *gorm.DB
	startTime time.Time
	version   string
	checkers  map[string]HealthCheckFunc
	mu        sync.RWMutex
}

type HealthCheckFunc func(ctx context.Context) ComponentHealth

func NewHealthChecker(db *gorm.DB, version string) *HealthChecker {
	return &HealthChecker{
		db:        db,
		startTime: time.Now(),
		version:   version,
		checkers:  make(map[string]HealthCheckFunc),
	}
}

func (hc *HealthChecker) RegisterChecker(name string, checker HealthCheckFunc) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.checkers[name] = checker
}

func (hc *HealthChecker) Check(ctx context.Context) HealthCheck {
	hc.mu.RLock()
	checkers := make(map[string]HealthCheckFunc, len(hc.checkers))
	for name, checker := range hc.checkers {
		checkers[name] = checker
	}
	hc.mu.RUnlock()

	components := make([]ComponentHealth, 0, len(checkers))
	overallStatus := HealthStatusHealthy

	for _, checker := range checkers {
		component := checker(ctx)
		components = append(components, component)

		if component.Status == HealthStatusUnhealthy {
			overallStatus = HealthStatusUnhealthy
		} else if component.Status == HealthStatusDegraded && overallStatus == HealthStatusHealthy {
			overallStatus = HealthStatusDegraded
		}
	}

	uptime := time.Since(hc.startTime)

	return HealthCheck{
		Status:     overallStatus,
		Version:    hc.version,
		Uptime:     formatDuration(uptime),
		Timestamp:  time.Now(),
		Components: components,
	}
}

func (hc *HealthChecker) DatabaseChecker() HealthCheckFunc {
	return func(ctx context.Context) ComponentHealth {
		start := time.Now()

		sqlDB, err := hc.db.DB()
		if err != nil {
			return ComponentHealth{
				Name:      "database",
				Status:    HealthStatusUnhealthy,
				Message:   fmt.Sprintf("Failed to get database: %v", err),
				Timestamp: time.Now(),
				Duration:  time.Since(start).String(),
			}
		}

		if err := sqlDB.PingContext(ctx); err != nil {
			return ComponentHealth{
				Name:      "database",
				Status:    HealthStatusUnhealthy,
				Message:   fmt.Sprintf("Database ping failed: %v", err),
				Timestamp: time.Now(),
				Duration:  time.Since(start).String(),
			}
		}

		return ComponentHealth{
			Name:      "database",
			Status:    HealthStatusHealthy,
			Message:   "Database connection is healthy",
			Timestamp: time.Now(),
			Duration:  time.Since(start).String(),
		}
	}
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
