package config

import "time"

type JWTConfig struct {
	Secret          string
	ExpirationHours int
	Issuer          string
	RefreshEnabled  bool
	RefreshExpHours int
}

func loadJWTConfig() JWTConfig {
	return JWTConfig{
		Secret:          getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		ExpirationHours: getEnvInt("JWT_EXPIRATION_HOURS", 24),
		Issuer:          getEnv("JWT_ISSUER", "base-golang-restful"),
		RefreshEnabled:  getEnvBool("JWT_REFRESH_ENABLED", true),
		RefreshExpHours: getEnvInt("JWT_REFRESH_EXP_HOURS", 168),
	}
}

func (j *JWTConfig) GetExpiration() time.Duration {
	return time.Duration(j.ExpirationHours) * time.Hour
}

func (j *JWTConfig) GetRefreshExpiration() time.Duration {
	return time.Duration(j.RefreshExpHours) * time.Hour
}