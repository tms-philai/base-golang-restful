package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	ExposeHeaders    []string
	MaxAge           int
}

func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"API-Version",
		},
		AllowCredentials: false,
		MaxAge:           12 * 3600,
	}
}

func CORS() gin.HandlerFunc {
	return CORSWithConfig(DefaultCORSConfig())
}

func CORSWithConfig(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if len(config.AllowOrigins) > 0 {
			if config.AllowOrigins[0] == "*" {
				c.Header("Access-Control-Allow-Origin", "*")
			} else {
				for _, allowedOrigin := range config.AllowOrigins {
					if origin == allowedOrigin {
						c.Header("Access-Control-Allow-Origin", origin)
						break
					}
				}
			}
		}

		if len(config.AllowMethods) > 0 {
			c.Header("Access-Control-Allow-Methods", joinStrings(config.AllowMethods, ", "))
		}

		if len(config.AllowHeaders) > 0 {
			c.Header("Access-Control-Allow-Headers", joinStrings(config.AllowHeaders, ", "))
		}

		if config.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if len(config.ExposeHeaders) > 0 {
			c.Header("Access-Control-Expose-Headers", joinStrings(config.ExposeHeaders, ", "))
		}

		if config.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", string(rune(config.MaxAge)))
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
