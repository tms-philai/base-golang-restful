package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	os.Setenv("ENVIRONMENT", "test")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("JWT_SECRET", "test-secret")

	config, err := Load("")
	if err != nil {
		t.Errorf("Load failed: %v", err)
	}

	if config.Environment != "test" {
		t.Errorf("Expected environment 'test', got '%s'", config.Environment)
	}

	if config.Server.Port != "9090" {
		t.Errorf("Expected port '9090', got '%s'", config.Server.Port)
	}

	if config.Database.Host != "testhost" {
		t.Errorf("Expected DB host 'testhost', got '%s'", config.Database.Host)
	}

	os.Unsetenv("ENVIRONMENT")
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("JWT_SECRET")
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "Valid config",
			config: &Config{
				Server: ServerConfig{
					Port: "8080",
				},
				Database: DatabaseConfig{
					Host: "localhost",
				},
				JWT: JWTConfig{
					Secret: "secret",
				},
			},
			wantErr: false,
		},
		{
			name: "Missing server port",
			config: &Config{
				Server: ServerConfig{
					Port: "",
				},
				Database: DatabaseConfig{
					Host: "localhost",
				},
				JWT: JWTConfig{
					Secret: "secret",
				},
			},
			wantErr: true,
		},
		{
			name: "Missing database host",
			config: &Config{
				Server: ServerConfig{
					Port: "8080",
				},
				Database: DatabaseConfig{
					Host: "",
				},
				JWT: JWTConfig{
					Secret: "secret",
				},
			},
			wantErr: true,
		},
		{
			name: "Missing JWT secret",
			config: &Config{
				Server: ServerConfig{
					Port: "8080",
				},
				Database: DatabaseConfig{
					Host: "localhost",
				},
				JWT: JWTConfig{
					Secret: "",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_IsDevelopment(t *testing.T) {
	config := &Config{Environment: "development"}
	if !config.IsDevelopment() {
		t.Error("Expected IsDevelopment() to return true")
	}

	config.Environment = "production"
	if config.IsDevelopment() {
		t.Error("Expected IsDevelopment() to return false")
	}
}

func TestConfig_IsProduction(t *testing.T) {
	config := &Config{Environment: "production"}
	if !config.IsProduction() {
		t.Error("Expected IsProduction() to return true")
	}

	config.Environment = "development"
	if config.IsProduction() {
		t.Error("Expected IsProduction() to return false")
	}
}

func TestConfig_IsStaging(t *testing.T) {
	config := &Config{Environment: "staging"}
	if !config.IsStaging() {
		t.Error("Expected IsStaging() to return true")
	}

	config.Environment = "production"
	if config.IsStaging() {
		t.Error("Expected IsStaging() to return false")
	}
}

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY")

	value := getEnv("TEST_KEY", "default")
	if value != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", value)
	}

	value = getEnv("NON_EXISTENT_KEY", "default")
	if value != "default" {
		t.Errorf("Expected 'default', got '%s'", value)
	}
}

func TestGetEnvInt(t *testing.T) {
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")

	value := getEnvInt("TEST_INT", 10)
	if value != 42 {
		t.Errorf("Expected 42, got %d", value)
	}

	value = getEnvInt("NON_EXISTENT_INT", 10)
	if value != 10 {
		t.Errorf("Expected 10, got %d", value)
	}
}

func TestGetEnvBool(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected bool
	}{
		{"true value", "true", true},
		{"1 value", "1", true},
		{"false value", "false", false},
		{"0 value", "0", false},
		{"other value", "other", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("TEST_BOOL", tt.envValue)
			defer os.Unsetenv("TEST_BOOL")

			value := getEnvBool("TEST_BOOL", false)
			if value != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, value)
			}
		})
	}

	value := getEnvBool("NON_EXISTENT_BOOL", true)
	if !value {
		t.Error("Expected true for default value")
	}
}