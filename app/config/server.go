package config

type ServerConfig struct {
	Port            string
	Host            string
	ReadTimeout     int
	WriteTimeout    int
	ShutdownTimeout int
	TrustedProxies  []string
	EnableCORS      bool
	CORSOrigins     []string
}

func loadServerConfig() ServerConfig {
	return ServerConfig{
		Port:            getEnv("SERVER_PORT", "8080"),
		Host:            getEnv("SERVER_HOST", "0.0.0.0"),
		ReadTimeout:     getEnvInt("SERVER_READ_TIMEOUT", 60),
		WriteTimeout:    getEnvInt("SERVER_WRITE_TIMEOUT", 60),
		ShutdownTimeout: getEnvInt("SERVER_SHUTDOWN_TIMEOUT", 30),
		TrustedProxies:  []string{},
		EnableCORS:      getEnvBool("SERVER_ENABLE_CORS", true),
		CORSOrigins:     []string{"*"},
	}
}