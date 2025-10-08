package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAPIVersioning_FromHeader(t *testing.T) {
	router := setupTestRouter()
	router.Use(APIVersioning())
	router.GET("/test", func(c *gin.Context) {
		version := GetAPIVersion(c)
		c.JSON(200, gin.H{"version": version})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set(VersionHeader, "v2")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "v2")
}

func TestAPIVersioning_FromQuery(t *testing.T) {
	router := setupTestRouter()
	router.Use(APIVersioning())
	router.GET("/test", func(c *gin.Context) {
		version := GetAPIVersion(c)
		c.JSON(200, gin.H{"version": version})
	})

	req := httptest.NewRequest("GET", "/test?version=v3", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "v3")
}

func TestAPIVersioning_FromPath(t *testing.T) {
	router := setupTestRouter()
	router.Use(APIVersioning())
	router.GET("/api/v2/test", func(c *gin.Context) {
		version := GetAPIVersion(c)
		c.JSON(200, gin.H{"version": version})
	})

	req := httptest.NewRequest("GET", "/api/v2/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "v2")
}

func TestAPIVersioning_Default(t *testing.T) {
	router := setupTestRouter()
	router.Use(APIVersioning())
	router.GET("/test", func(c *gin.Context) {
		version := GetAPIVersion(c)
		c.JSON(200, gin.H{"version": version})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), DefaultAPIVersion)
}

func TestNormalizeVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "with v prefix",
			input:    "v1",
			expected: "v1",
		},
		{
			name:     "without v prefix",
			input:    "2",
			expected: "v2",
		},
		{
			name:     "uppercase",
			input:    "V3",
			expected: "v3",
		},
		{
			name:     "with spaces",
			input:    " v4 ",
			expected: "v4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeVersion(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRequireVersion(t *testing.T) {
	router := setupTestRouter()
	router.Use(APIVersioning())
	router.GET("/test", RequireVersion("v2"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	tests := []struct {
		name           string
		version        string
		expectedStatus int
	}{
		{
			name:           "correct version",
			version:        "v2",
			expectedStatus: 200,
		},
		{
			name:           "wrong version",
			version:        "v1",
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set(VersionHeader, tt.version)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
