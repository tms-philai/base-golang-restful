package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockUser struct {
	roles       []string
	permissions []string
}

func (m *mockUser) HasRole(roleName string) bool {
	for _, role := range m.roles {
		if role == roleName {
			return true
		}
	}
	return false
}

func (m *mockUser) HasAnyRole(roleNames []string) bool {
	for _, role := range m.roles {
		for _, name := range roleNames {
			if role == name {
				return true
			}
		}
	}
	return false
}

func (m *mockUser) HasAllRoles(roleNames []string) bool {
	roleMap := make(map[string]bool)
	for _, role := range m.roles {
		roleMap[role] = true
	}
	for _, name := range roleNames {
		if !roleMap[name] {
			return false
		}
	}
	return true
}

func (m *mockUser) HasPermission(permissionName string) bool {
	for _, perm := range m.permissions {
		if perm == permissionName {
			return true
		}
	}
	return false
}

func (m *mockUser) HasAnyPermission(permissionNames []string) bool {
	for _, perm := range m.permissions {
		for _, name := range permissionNames {
			if perm == name {
				return true
			}
		}
	}
	return false
}

func (m *mockUser) HasAllPermissions(permissionNames []string) bool {
	permMap := make(map[string]bool)
	for _, perm := range m.permissions {
		permMap[perm] = true
	}
	for _, name := range permissionNames {
		if !permMap[name] {
			return false
		}
	}
	return true
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestRequireAuth(t *testing.T) {
	tests := []struct {
		name           string
		setUser        bool
		expectedStatus int
	}{
		{
			name:           "authenticated user",
			setUser:        true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "unauthenticated user",
			setUser:        false,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestRouter()
			router.Use(func(c *gin.Context) {
				if tt.setUser {
					c.Set(UserContextKey, &mockUser{})
				}
			})
			router.GET("/test", RequireAuth(), func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRequireRole(t *testing.T) {
	tests := []struct {
		name           string
		userRoles      []string
		requiredRole   string
		expectedStatus int
	}{
		{
			name:           "user has required role",
			userRoles:      []string{"admin", "editor"},
			requiredRole:   "admin",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "user does not have required role",
			userRoles:      []string{"editor"},
			requiredRole:   "admin",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestRouter()
			router.Use(func(c *gin.Context) {
				c.Set(UserContextKey, &mockUser{roles: tt.userRoles})
			})
			router.GET("/test", RequireRole(tt.requiredRole), func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRequireAnyRole(t *testing.T) {
	tests := []struct {
		name           string
		userRoles      []string
		requiredRoles  []string
		expectedStatus int
	}{
		{
			name:           "user has one of required roles",
			userRoles:      []string{"editor"},
			requiredRoles:  []string{"admin", "editor"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "user has none of required roles",
			userRoles:      []string{"viewer"},
			requiredRoles:  []string{"admin", "editor"},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestRouter()
			router.Use(func(c *gin.Context) {
				c.Set(UserContextKey, &mockUser{roles: tt.userRoles})
			})
			router.GET("/test", RequireAnyRole(tt.requiredRoles), func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRequireAllRoles(t *testing.T) {
	tests := []struct {
		name           string
		userRoles      []string
		requiredRoles  []string
		expectedStatus int
	}{
		{
			name:           "user has all required roles",
			userRoles:      []string{"admin", "editor"},
			requiredRoles:  []string{"admin", "editor"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "user missing one required role",
			userRoles:      []string{"admin"},
			requiredRoles:  []string{"admin", "editor"},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestRouter()
			router.Use(func(c *gin.Context) {
				c.Set(UserContextKey, &mockUser{roles: tt.userRoles})
			})
			router.GET("/test", RequireAllRoles(tt.requiredRoles), func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRequirePermission(t *testing.T) {
	tests := []struct {
		name               string
		userPermissions    []string
		requiredPermission string
		expectedStatus     int
	}{
		{
			name:               "user has required permission",
			userPermissions:    []string{"user.create", "user.read"},
			requiredPermission: "user.create",
			expectedStatus:     http.StatusOK,
		},
		{
			name:               "user does not have required permission",
			userPermissions:    []string{"user.read"},
			requiredPermission: "user.create",
			expectedStatus:     http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestRouter()
			router.Use(func(c *gin.Context) {
				c.Set(UserContextKey, &mockUser{permissions: tt.userPermissions})
			})
			router.GET("/test", RequirePermission(tt.requiredPermission), func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRequireAnyPermission(t *testing.T) {
	tests := []struct {
		name                string
		userPermissions     []string
		requiredPermissions []string
		expectedStatus      int
	}{
		{
			name:                "user has one of required permissions",
			userPermissions:     []string{"user.read"},
			requiredPermissions: []string{"user.create", "user.read"},
			expectedStatus:      http.StatusOK,
		},
		{
			name:                "user has none of required permissions",
			userPermissions:     []string{"user.update"},
			requiredPermissions: []string{"user.create", "user.read"},
			expectedStatus:      http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestRouter()
			router.Use(func(c *gin.Context) {
				c.Set(UserContextKey, &mockUser{permissions: tt.userPermissions})
			})
			router.GET("/test", RequireAnyPermission(tt.requiredPermissions), func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRequireAllPermissions(t *testing.T) {
	tests := []struct {
		name                string
		userPermissions     []string
		requiredPermissions []string
		expectedStatus      int
	}{
		{
			name:                "user has all required permissions",
			userPermissions:     []string{"user.create", "user.read"},
			requiredPermissions: []string{"user.create", "user.read"},
			expectedStatus:      http.StatusOK,
		},
		{
			name:                "user missing one required permission",
			userPermissions:     []string{"user.create"},
			requiredPermissions: []string{"user.create", "user.read"},
			expectedStatus:      http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestRouter()
			router.Use(func(c *gin.Context) {
				c.Set(UserContextKey, &mockUser{permissions: tt.userPermissions})
			})
			router.GET("/test", RequireAllPermissions(tt.requiredPermissions), func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
