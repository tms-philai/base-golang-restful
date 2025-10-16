package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	JWT         JWTConfig
}

var globalConfig *Config

func Load(envFile string) (*Config, error) {
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			log.Printf("Warning: Error loading %s file: %v", envFile, err)
		}
	}
	viper.AutomaticEnv()

	config := &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Server:      loadServerConfig(),
		Database:    loadDatabaseConfig(),
		JWT:         loadJWTConfig(),
	}
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	globalConfig = config
	return config, nil
}

func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}

	return nil
}

/*
- LoadConfig initializes the global configuration by reading environment variables.
- It sets up the server configuration and other application settings.
*/
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

/*
	- getEnvInt retrieves an integer value from environment variables.
	- If the variable is not set or cannot be converted, it returns the provided default value.
*/

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		fmt.Sscanf(value, "%d", &intValue)
		return intValue
	}
	return defaultValue
}

/*
	- getEnvDuration retrieves a duration value from environment variables.
	- If the variable is not set or cannot be parsed, it returns the provided default value.
*/

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		duration, err := time.ParseDuration(value)
		if err == nil {
			return duration
		}
	}
	return defaultValue
}

/*
- getEnvBool retrieves a boolean value from environment variables.
- If the variable is not set, it returns the provided default value.
*/
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1"
	}
	return defaultValue
}
