package config

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	TimeZone        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
	AutoMigrate     bool
	LogQueries      bool
}

func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnv("DB_PORT", "5432"),
		User:            getEnv("DB_USER", "postgres"),
		Password:        getEnv("DB_PASSWORD", ""),
		Name:            getEnv("DB_NAME", "app_db"),
		SSLMode:         getEnv("DB_SSL_MODE", "disable"),
		TimeZone:        getEnv("DB_TIMEZONE", "UTC"),
		MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
		ConnMaxLifetime: getEnvInt("DB_CONN_MAX_LIFETIME", 300),
		AutoMigrate:     getEnvBool("DB_AUTO_MIGRATE", false),
		LogQueries:      getEnvBool("DB_LOG_QUERIES", false),
	}
}
