package logger

import (
	"time"

	"github.com/rs/zerolog"
)

type Fields map[string]interface{}

func (f Fields) MarshalZerologObject(e *zerolog.Event) {
	for key, value := range f {
		switch v := value.(type) {
		case string:
			e.Str(key, v)
		case int:
			e.Int(key, v)
		case int64:
			e.Int64(key, v)
		case float64:
			e.Float64(key, v)
		case bool:
			e.Bool(key, v)
		case time.Time:
			e.Time(key, v)
		case time.Duration:
			e.Dur(key, v)
		case error:
			e.AnErr(key, v)
		default:
			e.Interface(key, v)
		}
	}
}

type HTTPFields struct {
	Method     string
	Path       string
	Status     int
	Duration   time.Duration
	IP         string
	UserAgent  string
	RequestID  string
	UserID     string
	Error      error
}

func (h HTTPFields) MarshalZerologObject(e *zerolog.Event) {
	e.Str("method", h.Method).
		Str("path", h.Path).
		Int("status", h.Status).
		Dur("duration_ms", h.Duration).
		Str("ip", h.IP).
		Str("user_agent", h.UserAgent)

	if h.RequestID != "" {
		e.Str("request_id", h.RequestID)
	}

	if h.UserID != "" {
		e.Str("user_id", h.UserID)
	}

	if h.Error != nil {
		e.Err(h.Error)
	}
}

type DBFields struct {
	Query    string
	Duration time.Duration
	Rows     int64
	Error    error
}

func (d DBFields) MarshalZerologObject(e *zerolog.Event) {
	e.Str("query", d.Query).
		Dur("duration_ms", d.Duration).
		Int64("rows", d.Rows)

	if d.Error != nil {
		e.Err(d.Error)
	}
}

type ErrorFields struct {
	Code       string
	Message    string
	StackTrace string
	Context    map[string]interface{}
}

func (ef ErrorFields) MarshalZerologObject(e *zerolog.Event) {
	e.Str("error_code", ef.Code).
		Str("error_message", ef.Message)

	if ef.StackTrace != "" {
		e.Str("stack_trace", ef.StackTrace)
	}

	if ef.Context != nil {
		e.Interface("error_context", ef.Context)
	}
}
