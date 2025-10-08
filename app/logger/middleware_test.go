package logger

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	buf := &bytes.Buffer{}
	Init(Config{
		Level:  "info",
		Pretty: false,
		Output: buf,
	})

	router := gin.New()
	router.Use(Middleware())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Correlation-ID", "test-correlation")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	correlationID := w.Header().Get("X-Correlation-ID")
	if correlationID != "test-correlation" {
		t.Errorf("Expected correlation ID 'test-correlation', got '%s'", correlationID)
	}

	requestID := w.Header().Get("X-Request-ID")
	if requestID == "" {
		t.Error("Expected request ID to be set")
	}

	output := buf.String()
	if !strings.Contains(output, "GET") {
		t.Error("Expected log to contain method")
	}
	if !strings.Contains(output, "/test") {
		t.Error("Expected log to contain path")
	}
	if !strings.Contains(output, "200") {
		t.Error("Expected log to contain status")
	}
}

func TestMiddleware_WithoutCorrelationID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	buf := &bytes.Buffer{}
	Init(Config{
		Level:  "info",
		Pretty: false,
		Output: buf,
	})

	router := gin.New()
	router.Use(Middleware())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	correlationID := w.Header().Get("X-Correlation-ID")
	if correlationID == "" {
		t.Error("Expected correlation ID to be generated")
	}
}

func TestMiddleware_WithError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	buf := &bytes.Buffer{}
	Init(Config{
		Level:  "error",
		Pretty: false,
		Output: buf,
	})

	router := gin.New()
	router.Use(Middleware())

	router.GET("/error", func(c *gin.Context) {
		c.Error(gin.Error{Err: http.ErrAbortHandler, Type: gin.ErrorTypePrivate})
		c.JSON(500, gin.H{"error": "internal error"})
	})

	req := httptest.NewRequest("GET", "/error", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	output := buf.String()
	if !strings.Contains(output, "error") {
		t.Error("Expected log to contain error level")
	}
}

func TestRecovery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	buf := &bytes.Buffer{}
	Init(Config{
		Level:  "error",
		Pretty: false,
		Output: buf,
	})

	router := gin.New()
	router.Use(Recovery())

	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	req := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	output := buf.String()
	if !strings.Contains(output, "panic") {
		t.Error("Expected log to contain panic message")
	}
}
