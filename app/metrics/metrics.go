package metrics

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Metrics struct {
	RequestCount      int64                  `json:"request_count"`
	ErrorCount        int64                  `json:"error_count"`
	TotalDuration     time.Duration          `json:"total_duration"`
	AverageDuration   time.Duration          `json:"average_duration"`
	EndpointMetrics   map[string]*EndpointMetric `json:"endpoint_metrics"`
	StatusCodeCounts  map[int]int64          `json:"status_code_counts"`
	mu                sync.RWMutex
}

type EndpointMetric struct {
	Path          string        `json:"path"`
	Method        string        `json:"method"`
	Count         int64         `json:"count"`
	ErrorCount    int64         `json:"error_count"`
	TotalDuration time.Duration `json:"total_duration"`
	AvgDuration   time.Duration `json:"avg_duration"`
	MinDuration   time.Duration `json:"min_duration"`
	MaxDuration   time.Duration `json:"max_duration"`
}

var globalMetrics = &Metrics{
	EndpointMetrics:  make(map[string]*EndpointMetric),
	StatusCodeCounts: make(map[int]int64),
}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		globalMetrics.mu.Lock()
		defer globalMetrics.mu.Unlock()

		globalMetrics.RequestCount++
		globalMetrics.TotalDuration += duration
		globalMetrics.AverageDuration = globalMetrics.TotalDuration / time.Duration(globalMetrics.RequestCount)

		if statusCode >= 400 {
			globalMetrics.ErrorCount++
		}

		globalMetrics.StatusCodeCounts[statusCode]++

		endpointKey := c.Request.Method + ":" + c.FullPath()
		metric, exists := globalMetrics.EndpointMetrics[endpointKey]
		
		if !exists {
			metric = &EndpointMetric{
				Path:        c.FullPath(),
				Method:      c.Request.Method,
				MinDuration: duration,
				MaxDuration: duration,
			}
			globalMetrics.EndpointMetrics[endpointKey] = metric
		}

		metric.Count++
		metric.TotalDuration += duration
		metric.AvgDuration = metric.TotalDuration / time.Duration(metric.Count)

		if duration < metric.MinDuration {
			metric.MinDuration = duration
		}
		if duration > metric.MaxDuration {
			metric.MaxDuration = duration
		}

		if statusCode >= 400 {
			metric.ErrorCount++
		}
	}
}

func GetMetrics() Metrics {
	globalMetrics.mu.RLock()
	defer globalMetrics.mu.RUnlock()

	metricsCopy := Metrics{
		RequestCount:     globalMetrics.RequestCount,
		ErrorCount:       globalMetrics.ErrorCount,
		TotalDuration:    globalMetrics.TotalDuration,
		AverageDuration:  globalMetrics.AverageDuration,
		EndpointMetrics:  make(map[string]*EndpointMetric),
		StatusCodeCounts: make(map[int]int64),
	}

	for k, v := range globalMetrics.EndpointMetrics {
		metricCopy := *v
		metricsCopy.EndpointMetrics[k] = &metricCopy
	}

	for k, v := range globalMetrics.StatusCodeCounts {
		metricsCopy.StatusCodeCounts[k] = v
	}

	return metricsCopy
}

func ResetMetrics() {
	globalMetrics.mu.Lock()
	defer globalMetrics.mu.Unlock()

	globalMetrics.RequestCount = 0
	globalMetrics.ErrorCount = 0
	globalMetrics.TotalDuration = 0
	globalMetrics.AverageDuration = 0
	globalMetrics.EndpointMetrics = make(map[string]*EndpointMetric)
	globalMetrics.StatusCodeCounts = make(map[int]int64)
}
