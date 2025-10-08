package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	DefaultAPIVersion = "v1"
	VersionHeader     = "API-Version"
	VersionParam      = "version"
)

func APIVersioning() gin.HandlerFunc {
	return func(c *gin.Context) {
		version := extractVersion(c)
		
		if version == "" {
			version = DefaultAPIVersion
		}

		c.Set("api_version", version)
		c.Header(VersionHeader, version)

		c.Next()
	}
}

func extractVersion(c *gin.Context) string {
	version := c.GetHeader(VersionHeader)
	if version != "" {
		return normalizeVersion(version)
	}

	version = c.Query(VersionParam)
	if version != "" {
		return normalizeVersion(version)
	}

	path := c.Request.URL.Path
	parts := strings.Split(strings.Trim(path, "/"), "/")
	
	for _, part := range parts {
		if strings.HasPrefix(part, "v") && len(part) > 1 {
			return normalizeVersion(part)
		}
	}

	return ""
}

func normalizeVersion(version string) string {
	version = strings.ToLower(strings.TrimSpace(version))
	
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	return version
}

func GetAPIVersion(c *gin.Context) string {
	version, exists := c.Get("api_version")
	if !exists {
		return DefaultAPIVersion
	}

	versionStr, ok := version.(string)
	if !ok {
		return DefaultAPIVersion
	}

	return versionStr
}

func RequireVersion(requiredVersion string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentVersion := GetAPIVersion(c)
		
		if currentVersion != requiredVersion {
			c.JSON(400, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNSUPPORTED_API_VERSION",
					"message": "API version not supported",
					"details": gin.H{
						"current":  currentVersion,
						"required": requiredVersion,
					},
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
