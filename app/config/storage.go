package config

type StorageConfig struct {
	Provider     string      `mapstructure:"provider"`
	Local        LocalConfig `mapstructure:"local"`
	GCS          GCSConfig   `mapstructure:"gcs"`
	MaxFileSize  int64       `mapstructure:"max_file_size"`
	AllowedTypes []string    `mapstructure:"allowed_types"`
}

type LocalConfig struct {
	Path      string `mapstructure:"path"`
	PublicURL string `mapstructure:"public_url"`
}

type GCSConfig struct {
	ProjectID           string `mapstructure:"project_id"`
	Bucket              string `mapstructure:"bucket"`
	CredentialsFile     string `mapstructure:"credentials_file"`
	CredentialsJSON     string `mapstructure:"credentials_json"`
	PublicURL           string `mapstructure:"public_url"`
	SignedURLExpiration int    `mapstructure:"signed_url_expiration"`
}
