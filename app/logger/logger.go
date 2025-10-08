package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
	logger zerolog.Logger
}

type Config struct {
	Level      string
	Pretty     bool
	Output     io.Writer
	TimeFormat string
}

var globalLogger *Logger

func Init(config Config) *Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if config.TimeFormat != "" {
		zerolog.TimeFieldFormat = config.TimeFormat
	}

	level := parseLevel(config.Level)
	zerolog.SetGlobalLevel(level)

	var output io.Writer = os.Stdout
	if config.Output != nil {
		output = config.Output
	}

	if config.Pretty {
		output = zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: time.RFC3339,
		}
	}

	logger := zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()

	globalLogger = &Logger{logger: logger}
	log.Logger = logger

	return globalLogger
}

func parseLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

func Get() *Logger {
	if globalLogger == nil {
		return Init(Config{Level: "info"})
	}
	return globalLogger
}

func (l *Logger) Debug() *zerolog.Event {
	return l.logger.Debug()
}

func (l *Logger) Info() *zerolog.Event {
	return l.logger.Info()
}

func (l *Logger) Warn() *zerolog.Event {
	return l.logger.Warn()
}

func (l *Logger) Error() *zerolog.Event {
	return l.logger.Error()
}

func (l *Logger) Fatal() *zerolog.Event {
	return l.logger.Fatal()
}

func (l *Logger) Panic() *zerolog.Event {
	return l.logger.Panic()
}

func (l *Logger) With() zerolog.Context {
	return l.logger.With()
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	return &Logger{
		logger: l.logger.With().Logger(),
	}
}

func Debug() *zerolog.Event {
	return Get().Debug()
}

func Info() *zerolog.Event {
	return Get().Info()
}

func Warn() *zerolog.Event {
	return Get().Warn()
}

func Error() *zerolog.Event {
	return Get().Error()
}

func Fatal() *zerolog.Event {
	return Get().Fatal()
}

func Panic() *zerolog.Event {
	return Get().Panic()
}

type ContextKey string

const (
	CorrelationIDKey ContextKey = "correlation_id"
	UserIDKey        ContextKey = "user_id"
	RequestIDKey     ContextKey = "request_id"
)

func FromContext(ctx context.Context) *Logger {
	logger := Get().logger

	if correlationID, ok := ctx.Value(CorrelationIDKey).(string); ok {
		logger = logger.With().Str("correlation_id", correlationID).Logger()
	}

	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		logger = logger.With().Str("user_id", userID).Logger()
	}

	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		logger = logger.With().Str("request_id", requestID).Logger()
	}

	return &Logger{logger: logger}
}
