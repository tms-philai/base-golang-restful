package config

import "time"

type JWTConfig struct {
	Secret               string
	SecretKey            string
	ExpirationHours      int
	Issuer               string
	RefreshEnabled       bool
	RefreshExpHours      int
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func loadJWTConfig() JWTConfig {
	secret := getEnv("JWT_SECRET", "your-secret-key-change-in-production")
	accessDuration := getEnvDuration("JWT_ACCESS_TOKEN_DURATION", 15*time.Minute)
	refreshDuration := getEnvDuration("JWT_REFRESH_TOKEN_DURATION", 7*24*time.Hour)

	return JWTConfig{
		Secret:               secret,
		SecretKey:            secret,
		ExpirationHours:      getEnvInt("JWT_EXPIRATION_HOURS", 24),
		Issuer:               getEnv("JWT_ISSUER", "base-golang-restful"),
		RefreshEnabled:       getEnvBool("JWT_REFRESH_ENABLED", true),
		RefreshExpHours:      getEnvInt("JWT_REFRESH_EXP_HOURS", 168),
		AccessTokenDuration:  accessDuration,
		RefreshTokenDuration: refreshDuration,
	}
}

func (j *JWTConfig) GetExpiration() time.Duration {
	if j.AccessTokenDuration > 0 {
		return j.AccessTokenDuration
	}
	return time.Duration(j.ExpirationHours) * time.Hour
}

func (j *JWTConfig) GetRefreshExpiration() time.Duration {
	if j.RefreshTokenDuration > 0 {
		return j.RefreshTokenDuration
	}
	return time.Duration(j.RefreshExpHours) * time.Hour
}
