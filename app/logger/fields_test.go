package logger

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestFields_MarshalZerologObject(t *testing.T) {
	buf := &bytes.Buffer{}
	config := Config{
		Level:  "info",
		Pretty: false,
		Output: buf,
	}

	logger := Init(config)

	fields := Fields{
		"string_field": "test",
		"int_field":    42,
		"bool_field":   true,
		"float_field":  3.14,
	}

	logger.Info().Object("fields", fields).Msg("test fields")

	output := buf.String()
	if !strings.Contains(output, "test") {
		t.Error("Expected output to contain string field")
	}
	if !strings.Contains(output, "42") {
		t.Error("Expected output to contain int field")
	}
	if !strings.Contains(output, "true") {
		t.Error("Expected output to contain bool field")
	}
}

func TestHTTPFields_MarshalZerologObject(t *testing.T) {
	buf := &bytes.Buffer{}
	config := Config{
		Level:  "info",
		Pretty: false,
		Output: buf,
	}

	logger := Init(config)

	httpFields := HTTPFields{
		Method:    "GET",
		Path:      "/api/users",
		Status:    200,
		Duration:  100 * time.Millisecond,
		IP:        "127.0.0.1",
		UserAgent: "test-agent",
		RequestID: "req-123",
		UserID:    "user-456",
	}

	logger.Info().Object("http", httpFields).Msg("http request")

	output := buf.String()
	if !strings.Contains(output, "GET") {
		t.Error("Expected output to contain method")
	}
	if !strings.Contains(output, "/api/users") {
		t.Error("Expected output to contain path")
	}
	if !strings.Contains(output, "200") {
		t.Error("Expected output to contain status")
	}
	if !strings.Contains(output, "127.0.0.1") {
		t.Error("Expected output to contain IP")
	}
}

func TestHTTPFields_WithError(t *testing.T) {
	buf := &bytes.Buffer{}
	config := Config{
		Level:  "error",
		Pretty: false,
		Output: buf,
	}

	logger := Init(config)

	httpFields := HTTPFields{
		Method:   "POST",
		Path:     "/api/users",
		Status:   500,
		Duration: 50 * time.Millisecond,
		IP:       "127.0.0.1",
		Error:    errors.New("internal server error"),
	}

	logger.Error().Object("http", httpFields).Msg("http error")

	output := buf.String()
	if !strings.Contains(output, "internal server error") {
		t.Error("Expected output to contain error message")
	}
}

func TestDBFields_MarshalZerologObject(t *testing.T) {
	buf := &bytes.Buffer{}
	config := Config{
		Level:  "info",
		Pretty: false,
		Output: buf,
	}

	logger := Init(config)

	dbFields := DBFields{
		Query:    "SELECT * FROM users",
		Duration: 25 * time.Millisecond,
		Rows:     10,
	}

	logger.Info().Object("db", dbFields).Msg("database query")

	output := buf.String()
	if !strings.Contains(output, "SELECT * FROM users") {
		t.Error("Expected output to contain query")
	}
	if !strings.Contains(output, "10") {
		t.Error("Expected output to contain rows count")
	}
}

func TestErrorFields_MarshalZerologObject(t *testing.T) {
	buf := &bytes.Buffer{}
	config := Config{
		Level:  "error",
		Pretty: false,
		Output: buf,
	}

	logger := Init(config)

	errorFields := ErrorFields{
		Code:       "ERR_001",
		Message:    "Something went wrong",
		StackTrace: "line 1\nline 2",
		Context: map[string]interface{}{
			"user_id": "123",
			"action":  "create",
		},
	}

	logger.Error().Object("error", errorFields).Msg("error occurred")

	output := buf.String()
	if !strings.Contains(output, "ERR_001") {
		t.Error("Expected output to contain error code")
	}
	if !strings.Contains(output, "Something went wrong") {
		t.Error("Expected output to contain error message")
	}
}
