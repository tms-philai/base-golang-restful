// Package test contains test dependencies and imports
// This file ensures all test dependencies are properly imported
package test

import (
	// Test dependencies
	_ "bytes"
	_ "context"
	_ "encoding/json"
	_ "errors"
	_ "mime/multipart"
	_ "net/http"
	_ "net/http/httptest"
	_ "testing"
	_ "time"

	_ "github.com/gin-gonic/gin"
	_ "github.com/google/uuid"
	_ "github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/mock"
	_ "github.com/stretchr/testify/suite"
)
