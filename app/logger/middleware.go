package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		requestID := uuid.New().String()

		c.Set(string(CorrelationIDKey), correlationID)
		c.Set(string(RequestIDKey), requestID)

		c.Header("X-Correlation-ID", correlationID)
		c.Header("X-Request-ID", requestID)

		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		fields := HTTPFields{
			Method:    method,
			Path:      path,
			Status:    status,
			Duration:  duration,
			IP:        c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			RequestID: requestID,
		}

		if userID, exists := c.Get("user_id"); exists {
			if uid, ok := userID.(string); ok {
				fields.UserID = uid
			}
		}

		logger := Get()

		if len(c.Errors) > 0 {
			fields.Error = c.Errors.Last()
			logger.Error().Object("http", fields).Msg("Request completed with errors")
		} else if status >= 500 {
			logger.Error().Object("http", fields).Msg("Request failed")
		} else if status >= 400 {
			logger.Warn().Object("http", fields).Msg("Request completed with client error")
		} else {
			logger.Info().Object("http", fields).Msg("Request completed")
		}
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger := FromContext(c.Request.Context())

				logger.Error().
					Interface("panic", err).
					Str("path", c.Request.URL.Path).
					Str("method", c.Request.Method).
					Msg("Panic recovered")

				c.AbortWithStatus(500)
			}
		}()

		c.Next()
	}
}
