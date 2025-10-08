package logger

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	buf := &bytes.Buffer{}
	config := Config{
		Level:  "info",
		Pretty: false,
		Output: buf,
	}

	logger := Init(config)
	if logger == nil {
		t.Error("Expected Init to return non-nil logger")
	}

	logger.Info().Msg("test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Error("Expected log output to contain 'test message'")
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		expected string
	}{
		{"Debug level", "debug", "debug"},
		{"Info level", "info", "info"},
		{"Warn level", "warn", "warn"},
		{"Error level", "error", "error"},
		{"Fatal level", "fatal", "fatal"},
		{"Panic level", "panic", "panic"},
		{"Unknown level", "unknown", "info"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := parseLevel(tt.level)
			if level.String() != tt.expected {
				t.Errorf("Expected level %s, got %s", tt.expected, level.String())
			}
		})
	}
}

func TestGet(t *testing.T) {
	globalLogger = nil

	logger := Get()
	if logger == nil {
		t.Error("Expected Get to return non-nil logger")
	}

	logger2 := Get()
	if logger != logger2 {
		t.Error("Expected Get to return the same logger instance")
	}
}

func TestLoggerMethods(t *testing.T) {
	buf := &bytes.Buffer{}
	config := Config{
		Level:  "debug",
		Pretty: false,
		Output: buf,
	}

	logger := Init(config)

	tests := []struct {
		name    string
		logFunc func()
		message string
	}{
		{
			name: "Debug log",
			logFunc: func() {
				logger.Debug().Msg("debug message")
			},
			message: "debug message",
		},
		{
			name: "Info log",
			logFunc: func() {
				logger.Info().Msg("info message")
			},
			message: "info message",
		},
		{
			name: "Warn log",
			logFunc: func() {
				logger.Warn().Msg("warn message")
			},
			message: "warn message",
		},
		{
			name: "Error log",
			logFunc: func() {
				logger.Error().Msg("error message")
			},
			message: "error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc()

			output := buf.String()
			if !strings.Contains(output, tt.message) {
				t.Errorf("Expected log output to contain '%s', got: %s", tt.message, output)
			}
		})
	}
}

func TestFromContext(t *testing.T) {
	buf := &bytes.Buffer{}
	config := Config{
		Level:  "info",
		Pretty: false,
		Output: buf,
	}

	Init(config)

	ctx := context.Background()
	ctx = context.WithValue(ctx, CorrelationIDKey, "test-correlation-id")
	ctx = context.WithValue(ctx, UserIDKey, "test-user-id")
	ctx = context.WithValue(ctx, RequestIDKey, "test-request-id")

	logger := FromContext(ctx)
	logger.Info().Msg("test with context")

	output := buf.String()
	if !strings.Contains(output, "test-correlation-id") {
		t.Error("Expected log to contain correlation_id")
	}
	if !strings.Contains(output, "test-user-id") {
		t.Error("Expected log to contain user_id")
	}
	if !strings.Contains(output, "test-request-id") {
		t.Error("Expected log to contain request_id")
	}
}

func TestGlobalLoggerFunctions(t *testing.T) {
	buf := &bytes.Buffer{}
	config := Config{
		Level:  "debug",
		Pretty: false,
		Output: buf,
	}

	Init(config)

	tests := []struct {
		name    string
		logFunc func()
		message string
	}{
		{
			name: "Global Debug",
			logFunc: func() {
				Debug().Msg("global debug")
			},
			message: "global debug",
		},
		{
			name: "Global Info",
			logFunc: func() {
				Info().Msg("global info")
			},
			message: "global info",
		},
		{
			name: "Global Warn",
			logFunc: func() {
				Warn().Msg("global warn")
			},
			message: "global warn",
		},
		{
			name: "Global Error",
			logFunc: func() {
				Error().Msg("global error")
			},
			message: "global error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc()

			output := buf.String()
			if !strings.Contains(output, tt.message) {
				t.Errorf("Expected log output to contain '%s'", tt.message)
			}
		})
	}
}
