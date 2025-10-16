package middleware

import (
	"base-gin/internal/domain/services"
	"base-gin/internal/pkg/auth"
	"base-gin/internal/pkg/i18n"
	"strings"

	appErrors "base-gin/internal/pkg/errors"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtManager  *auth.JWTManager
	userService *services.UserService
}

func NewAuthMiddleware(jwtManager *auth.JWTManager, userService *services.UserService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager:  jwtManager,
		userService: userService,
	}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			lang := i18n.GetLanguage(c)
			_ = c.Error(appErrors.LocalizedUnauthorized("auth.token_invalid", nil).Localize(lang))
			c.Abort()
			return
		}

		claims, err := m.jwtManager.ValidateToken(token, auth.AccessToken)
		if err != nil {
			lang := i18n.GetLanguage(c)
			var localizedErr *appErrors.AppError

			// Log error for debugging
			c.Error(err)

			switch err {
			case auth.ErrExpiredToken:
				localizedErr = appErrors.LocalizedUnauthorized("auth.token_expired", nil).Localize(lang)
			case auth.ErrInvalidSignature, auth.ErrMissingClaims:
				localizedErr = appErrors.LocalizedUnauthorized("auth.token_invalid", nil).Localize(lang)
			default:
				localizedErr = appErrors.LocalizedUnauthorized("auth.token_invalid", nil).Localize(lang)
			}

			_ = c.Error(localizedErr)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID.String())
		c.Set("user_email", claims.Email)
		c.Set("token_claims", claims)

		// Load full user with roles and permissions for RBAC
		user, err := m.userService.GetByIDWithRoles(claims.UserID.String())
		if err == nil && user != nil {
			c.Set("user", user)
		}

		c.Next()
	}
}

func (m *AuthMiddleware) OptionalAuthenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.Next()
			return
		}

		claims, err := m.jwtManager.ValidateToken(token, auth.AccessToken)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_id", claims.UserID.String())
		c.Set("user_email", claims.Email)
		c.Set("token_claims", claims)

		// Load full user with roles and permissions (optional auth)
		user, err := m.userService.GetByIDWithRoles(claims.UserID.String())
		if err == nil && user != nil {
			c.Set("user", user)
		}

		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")
	if bearerToken != "" {
		parts := strings.SplitN(bearerToken, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	tokenQuery := c.Query("token")
	if tokenQuery != "" {
		return tokenQuery
	}

	tokenCookie, err := c.Cookie("access_token")
	if err == nil && tokenCookie != "" {
		return tokenCookie
	}

	return ""
}

func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", false
	}

	return userIDStr, true
}

func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("user_email")
	if !exists {
		return "", false
	}

	emailStr, ok := email.(string)
	if !ok {
		return "", false
	}

	return emailStr, true
}

func GetTokenClaims(c *gin.Context) (*auth.JWTClaims, bool) {
	claims, exists := c.Get("token_claims")
	if !exists {
		return nil, false
	}

	jwtClaims, ok := claims.(*auth.JWTClaims)
	if !ok {
		return nil, false
	}

	return jwtClaims, true
}

func GetCurrentUserID(c *gin.Context) string {
	userID, _ := GetUserID(c)
	return userID
}

func GetCurrentUser(c *gin.Context) (interface{}, bool) {
	user, exists := c.Get("user")
	return user, exists
}
