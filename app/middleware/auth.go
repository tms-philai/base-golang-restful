package middleware

import (
	"errors"
	"net/http"

	"base-golang-restful-app/models"
	"base-golang-restful-app/services"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates authentication middleware
func AuthMiddleware(jwtService *services.JWTService, userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		
		// Extract token from header
		token, err := jwtService.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "unauthorized",
				Message: err.Error(),
			})
			c.Abort()
			return
		}

		// Validate token
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "invalid or expired token",
			})
			c.Abort()
			return
		}

		// Get user from database to ensure user still exists and is active
		user, err := userService.GetByID(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "user not found",
			})
			c.Abort()
			return
		}

		if !user.IsActive {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "user account is inactive",
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", user.ID)
		c.Set("user", user)
		c.Set("claims", claims)

		c.Next()
	}
}

// OptionalAuthMiddleware creates optional authentication middleware
// This middleware will set user context if token is provided and valid,
// but won't abort if no token is provided
func OptionalAuthMiddleware(jwtService *services.JWTService, userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		
		if authHeader == "" {
			c.Next()
			return
		}

		// Extract token from header
		token, err := jwtService.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.Next()
			return
		}

		// Validate token
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			c.Next()
			return
		}

		// Get user from database
		user, err := userService.GetByID(claims.UserID)
		if err != nil || !user.IsActive {
			c.Next()
			return
		}

		// Set user information in context
		c.Set("user_id", user.ID)
		c.Set("user", user)
		c.Set("claims", claims)

		c.Next()
	}
}

// RequireRole creates role-based authorization middleware
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "authentication required",
			})
			c.Abort()
			return
		}

		userModel, ok := user.(*models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "internal_error",
				Message: "invalid user context",
			})
			c.Abort()
			return
		}

		// Check if user has required role
		hasRole := false
		for _, role := range roles {
			if userModel.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, models.ErrorResponse{
				Error:   "forbidden",
				Message: "insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetCurrentUser gets the current authenticated user from context
func GetCurrentUser(c *gin.Context) (*models.User, error) {
	user, exists := c.Get("user")
	if !exists {
		return nil, errors.New("user not found in context")
	}

	userModel, ok := user.(*models.User)
	if !ok {
		return nil, errors.New("invalid user type in context")
	}

	return userModel, nil
}

// GetCurrentUserID gets the current authenticated user ID from context
func GetCurrentUserID(c *gin.Context) (string, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", errors.New("user ID not found in context")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", errors.New("invalid user ID type in context")
	}

	return userIDStr, nil
}
