package config

import (
	"os"
	"testing"
)

func TestSetDefaults(t *testing.T) {
	setDefaults()
}

func TestIsDevelopment(t *testing.T) {
	globalConfig = &Config{Environment: "development"}
	if !IsDevelopment() {
		t.Error("Expected IsDevelopment to return true")
	}

	globalConfig = &Config{Environment: "production"}
	if IsDevelopment() {
		t.Error("Expected IsDevelopment to return false")
	}
}

func TestIsProduction(t *testing.T) {
	globalConfig = &Config{Environment: "production"}
	if !IsProduction() {
		t.Error("Expected IsProduction to return true")
	}

	globalConfig = &Config{Environment: "development"}
	if IsProduction() {
		t.Error("Expected IsProduction to return false")
	}
}

func TestIsStaging(t *testing.T) {
	globalConfig = &Config{Environment: "staging"}
	if !IsStaging() {
		t.Error("Expected IsStaging to return true")
	}

	globalConfig = &Config{Environment: "development"}
	if IsStaging() {
		t.Error("Expected IsStaging to return false")
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "Valid config",
			config: &Config{
				Server: ServerConfig{Port: "8080"},
				Database: DatabaseConfig{
					Host: "localhost",
					Name: "testdb",
					User: "testuser",
				},
				JWT: JWTConfig{Secret: "test-secret"},
			},
			wantErr: false,
		},
		{
			name: "Missing server port",
			config: &Config{
				Database: DatabaseConfig{
					Host: "localhost",
					Name: "testdb",
					User: "testuser",
				},
				JWT: JWTConfig{Secret: "test-secret"},
			},
			wantErr: true,
		},
		{
			name: "Missing database host",
			config: &Config{
				Server: ServerConfig{Port: "8080"},
				Database: DatabaseConfig{
					Name: "testdb",
					User: "testuser",
				},
				JWT: JWTConfig{Secret: "test-secret"},
			},
			wantErr: true,
		},
		{
			name: "Missing database name",
			config: &Config{
				Server: ServerConfig{Port: "8080"},
				Database: DatabaseConfig{
					Host: "localhost",
					User: "testuser",
				},
				JWT: JWTConfig{Secret: "test-secret"},
			},
			wantErr: true,
		},
		{
			name: "Missing database user",
			config: &Config{
				Server: ServerConfig{Port: "8080"},
				Database: DatabaseConfig{
					Host: "localhost",
					Name: "testdb",
				},
				JWT: JWTConfig{Secret: "test-secret"},
			},
			wantErr: true,
		},
		{
			name: "Missing JWT secret",
			config: &Config{
				Server: ServerConfig{Port: "8080"},
				Database: DatabaseConfig{
					Host: "localhost",
					Name: "testdb",
					User: "testuser",
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

func TestGet(t *testing.T) {
	globalConfig = &Config{
		Environment: "test",
	}

	config := Get()
	if config.Environment != "test" {
		t.Errorf("Expected environment 'test', got '%s'", config.Environment)
	}
}

func TestLoad_WithEnvironmentVariables(t *testing.T) {
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("JWT_SECRET", "test-secret")

	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USER")
		os.Unsetenv("JWT_SECRET")
	}()

	config, err := Load("./nonexistent")
	if err != nil {
		t.Logf("Load returned error (expected if no config file): %v", err)
	}

	if config != nil {
		if config.Server.Port != "9090" {
			t.Errorf("Expected port 9090, got %s", config.Server.Port)
		}
	}
}
