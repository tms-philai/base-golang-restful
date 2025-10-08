package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	UserContextKey = "user"
)

type RBACConfig struct {
	UnauthorizedHandler func(*gin.Context)
	ForbiddenHandler    func(*gin.Context)
}

func defaultUnauthorizedHandler(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "UNAUTHORIZED",
			"message": "Authentication required",
		},
	})
	c.Abort()
}

func defaultForbiddenHandler(c *gin.Context) {
	c.JSON(http.StatusForbidden, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "FORBIDDEN",
			"message": "You don't have permission to access this resource",
		},
	})
	c.Abort()
}

func RequireAuth(config ...RBACConfig) gin.HandlerFunc {
	cfg := getConfig(config)

	return func(c *gin.Context) {
		user, exists := c.Get(UserContextKey)
		if !exists || user == nil {
			cfg.UnauthorizedHandler(c)
			return
		}

		c.Next()
	}
}

func RequireRole(roleName string, config ...RBACConfig) gin.HandlerFunc {
	cfg := getConfig(config)

	return func(c *gin.Context) {
		user, exists := c.Get(UserContextKey)
		if !exists || user == nil {
			cfg.UnauthorizedHandler(c)
			return
		}

		userWithRoles, ok := user.(interface{ HasRole(string) bool })
		if !ok {
			cfg.ForbiddenHandler(c)
			return
		}

		if !userWithRoles.HasRole(roleName) {
			cfg.ForbiddenHandler(c)
			return
		}

		c.Next()
	}
}

func RequireAnyRole(roleNames []string, config ...RBACConfig) gin.HandlerFunc {
	cfg := getConfig(config)

	return func(c *gin.Context) {
		user, exists := c.Get(UserContextKey)
		if !exists || user == nil {
			cfg.UnauthorizedHandler(c)
			return
		}

		userWithRoles, ok := user.(interface{ HasAnyRole([]string) bool })
		if !ok {
			cfg.ForbiddenHandler(c)
			return
		}

		if !userWithRoles.HasAnyRole(roleNames) {
			cfg.ForbiddenHandler(c)
			return
		}

		c.Next()
	}
}

func RequireAllRoles(roleNames []string, config ...RBACConfig) gin.HandlerFunc {
	cfg := getConfig(config)

	return func(c *gin.Context) {
		user, exists := c.Get(UserContextKey)
		if !exists || user == nil {
			cfg.UnauthorizedHandler(c)
			return
		}

		userWithRoles, ok := user.(interface{ HasAllRoles([]string) bool })
		if !ok {
			cfg.ForbiddenHandler(c)
			return
		}

		if !userWithRoles.HasAllRoles(roleNames) {
			cfg.ForbiddenHandler(c)
			return
		}

		c.Next()
	}
}

func RequirePermission(permissionName string, config ...RBACConfig) gin.HandlerFunc {
	cfg := getConfig(config)

	return func(c *gin.Context) {
		user, exists := c.Get(UserContextKey)
		if !exists || user == nil {
			cfg.UnauthorizedHandler(c)
			return
		}

		userWithPerms, ok := user.(interface{ HasPermission(string) bool })
		if !ok {
			cfg.ForbiddenHandler(c)
			return
		}

		if !userWithPerms.HasPermission(permissionName) {
			cfg.ForbiddenHandler(c)
			return
		}

		c.Next()
	}
}

func RequireAnyPermission(permissionNames []string, config ...RBACConfig) gin.HandlerFunc {
	cfg := getConfig(config)

	return func(c *gin.Context) {
		user, exists := c.Get(UserContextKey)
		if !exists || user == nil {
			cfg.UnauthorizedHandler(c)
			return
		}

		userWithPerms, ok := user.(interface{ HasAnyPermission([]string) bool })
		if !ok {
			cfg.ForbiddenHandler(c)
			return
		}

		if !userWithPerms.HasAnyPermission(permissionNames) {
			cfg.ForbiddenHandler(c)
			return
		}

		c.Next()
	}
}

func RequireAllPermissions(permissionNames []string, config ...RBACConfig) gin.HandlerFunc {
	cfg := getConfig(config)

	return func(c *gin.Context) {
		user, exists := c.Get(UserContextKey)
		if !exists || user == nil {
			cfg.UnauthorizedHandler(c)
			return
		}

		userWithPerms, ok := user.(interface{ HasAllPermissions([]string) bool })
		if !ok {
			cfg.ForbiddenHandler(c)
			return
		}

		if !userWithPerms.HasAllPermissions(permissionNames) {
			cfg.ForbiddenHandler(c)
			return
		}

		c.Next()
	}
}

func getConfig(configs []RBACConfig) RBACConfig {
	if len(configs) > 0 {
		cfg := configs[0]
		if cfg.UnauthorizedHandler == nil {
			cfg.UnauthorizedHandler = defaultUnauthorizedHandler
		}
		if cfg.ForbiddenHandler == nil {
			cfg.ForbiddenHandler = defaultForbiddenHandler
		}
		return cfg
	}

	return RBACConfig{
		UnauthorizedHandler: defaultUnauthorizedHandler,
		ForbiddenHandler:    defaultForbiddenHandler,
	}
}
