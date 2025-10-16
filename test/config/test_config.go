package config

import (
	"base-gin/internal/pkg/config"
	"os"
	"path/filepath"
	"time"
)

// TestConfig holds test configuration
type TestConfig struct {
	Database config.DatabaseConfig
	Server   config.ServerConfig
	JWT      config.JWTConfig
}

// LoadTestConfig loads test configuration
func LoadTestConfig() (*TestConfig, error) {
	// Create test environment file if it doesn't exist
	testEnvPath := filepath.Join("..", "..", "configs", ".env.test")
	if _, err := os.Stat(testEnvPath); os.IsNotExist(err) {
		if err := createTestEnvFile(testEnvPath); err != nil {
			return nil, err
		}
	}

	// Load configuration
	cfg, err := config.Load(testEnvPath)
	if err != nil {
		return nil, err
	}

	return &TestConfig{
		Database: cfg.Database,
		Server:   cfg.Server,
		JWT:      cfg.JWT,
	}, nil
}

// createTestEnvFile creates a test environment file
func createTestEnvFile(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Create test environment content
	content := `# Test Environment Configuration

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=test_user
DB_PASSWORD=test_password
DB_NAME=test_base_gin
DB_SSLMODE=disable
DB_TIMEZONE=UTC

# Server Configuration
SERVER_HOST=localhost
SERVER_PORT=8001
SERVER_READ_TIMEOUT=30
SERVER_WRITE_TIMEOUT=30

# JWT Configuration
JWT_SECRET_KEY=test-secret-key-for-testing-only
JWT_ACCESS_TOKEN_DURATION=15m
JWT_REFRESH_TOKEN_DURATION=7d

# Email Configuration
SMTP_HOST=localhost
SMTP_PORT=587
SMTP_USERNAME=test@example.com
SMTP_PASSWORD=test_password
SMTP_FROM_EMAIL=test@example.com
SMTP_FROM_NAME=Test Sender

# Storage Configuration
STORAGE_TYPE=local
LOCAL_STORAGE_PATH=/tmp/test_files
S3_BUCKET=test-bucket
S3_REGION=us-east-1
S3_ACCESS_KEY=test-access-key
S3_SECRET_KEY=test-secret-key

# File Upload Configuration
MAX_FILE_SIZE=10485760
ALLOWED_FILE_TYPES=image/jpeg,image/png,text/plain,application/pdf

# Logging Configuration
LOG_LEVEL=debug
LOG_FORMAT=json
`

	return os.WriteFile(path, []byte(content), 0644)
}

// GetTestDatabaseConfig returns test database configuration
func GetTestDatabaseConfig() config.DatabaseConfig {
	return config.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "test_user",
		Password: "test_password",
		Name:     "test_base_gin",
		SSLMode:  "disable",
		TimeZone: "UTC",
	}
}

// GetTestServerConfig returns test server configuration
func GetTestServerConfig() config.ServerConfig {
	return config.ServerConfig{
		Host:         "localhost",
		Port:         "8001",
		ReadTimeout:  30,
		WriteTimeout: 30,
	}
}

// GetTestJWTConfig returns test JWT configuration
func GetTestJWTConfig() config.JWTConfig {
	return config.JWTConfig{
		SecretKey:            "test-secret-key-for-testing-only",
		AccessTokenDuration:  time.Duration(15) * time.Minute,
		RefreshTokenDuration: time.Duration(7) * 24 * time.Hour,
	}
}

// CleanupTestConfig cleans up test configuration files
func CleanupTestConfig() error {
	testEnvPath := filepath.Join("..", "..", "configs", ".env.test")
	if _, err := os.Stat(testEnvPath); err == nil {
		return os.Remove(testEnvPath)
	}
	return nil
}
