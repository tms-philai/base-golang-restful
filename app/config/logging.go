package config

type LoggingConfig struct {
	Level      string
	Pretty     bool
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

func loadLoggingConfig() LoggingConfig {
	return LoggingConfig{
		Level:      getEnv("LOG_LEVEL", "info"),
		Pretty:     getEnvBool("LOG_PRETTY", true),
		FilePath:   getEnv("LOG_FILE_PATH", "./logs/app.log"),
		MaxSize:    getEnvInt("LOG_MAX_SIZE", 100),
		MaxBackups: getEnvInt("LOG_MAX_BACKUPS", 7),
		MaxAge:     getEnvInt("LOG_MAX_AGE", 30),
	}
}
