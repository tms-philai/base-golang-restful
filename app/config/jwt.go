package config

type JWTConfig struct {
	Secret            string `mapstructure:"secret"`
	Expiration        int    `mapstructure:"expiration"`
	RefreshExpiration int    `mapstructure:"refresh_expiration"`
	Issuer            string `mapstructure:"issuer"`
}
