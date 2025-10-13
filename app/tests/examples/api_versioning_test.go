package examples

import (
	"net/http"
	"testing"

	"base-golang-restful-app/middleware"
	helperPkg "base-golang-restful-app/testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAPIVersioning(t *testing.T) {
	router := setupTestRouter()
	helper := helperPkg.NewTestHelper(t, router)

	t.Run("Default version should be v1", func(t *testing.T) {
		w := helper.GET("/test-version")
		helper.AssertStatus(w, http.StatusOK)

		var response map[string]interface{}
		helper.ParseJSON(w, &response)
		assert.Equal(t, "v1", response["version"])
	})

	t.Run("Version from header", func(t *testing.T) {
		w := helper.GET("/test-version", map[string]string{
			"API-Version": "v2",
		})
		helper.AssertStatus(w, http.StatusOK)

		var response map[string]interface{}
		helper.ParseJSON(w, &response)
		assert.Equal(t, "v2", response["version"])
	})

	t.Run("Version from query parameter", func(t *testing.T) {
		w := helper.GET("/test-version?version=v2")
		helper.AssertStatus(w, http.StatusOK)

		var response map[string]interface{}
		helper.ParseJSON(w, &response)
		assert.Equal(t, "v2", response["version"])
	})

	t.Run("Header takes precedence over query", func(t *testing.T) {
		w := helper.GET("/test-version?version=v1", map[string]string{
			"API-Version": "v2",
		})
		helper.AssertStatus(w, http.StatusOK)

		var response map[string]interface{}
		helper.ParseJSON(w, &response)
		assert.Equal(t, "v2", response["version"])
	})
}

func TestVersionRequirement(t *testing.T) {
	router := setupTestRouter()
	helper := helperPkg.NewTestHelper(t, router)

	t.Run("V1 endpoint with v1 version should work", func(t *testing.T) {
		w := helper.GET("/v1-only", map[string]string{
			"API-Version": "v1",
		})
		helper.AssertStatus(w, http.StatusOK)
	})

	t.Run("V1 endpoint with v2 version should fail", func(t *testing.T) {
		w := helper.GET("/v1-only", map[string]string{
			"API-Version": "v2",
		})
		helper.AssertStatus(w, http.StatusBadRequest)
		helper.AssertContains(w, "UNSUPPORTED_API_VERSION")
	})

	t.Run("V2 endpoint with v2 version should work", func(t *testing.T) {
		w := helper.GET("/v2-only", map[string]string{
			"API-Version": "v2",
		})
		helper.AssertStatus(w, http.StatusOK)
	})

	t.Run("V2 endpoint with v1 version should fail", func(t *testing.T) {
		w := helper.GET("/v2-only", map[string]string{
			"API-Version": "v1",
		})
		helper.AssertStatus(w, http.StatusBadRequest)
		helper.AssertContains(w, "UNSUPPORTED_API_VERSION")
	})
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.APIVersioning())

	r.GET("/test-version", func(c *gin.Context) {
		version := middleware.GetAPIVersion(c)
		c.JSON(http.StatusOK, gin.H{
			"version": version,
		})
	})

	r.GET("/v1-only", middleware.RequireVersion("v1"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "V1 endpoint",
		})
	})

	r.GET("/v2-only", middleware.RequireVersion("v2"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "V2 endpoint",
		})
	})

	return r
}
