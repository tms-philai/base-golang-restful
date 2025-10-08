package config

type StorageConfig struct {
	Provider    string
	LocalPath   string
	GCSBucket   string
	GCSProject  string
	S3Bucket    string
	S3Region    string
	MaxFileSize int64
}

func loadStorageConfig() StorageConfig {
	return StorageConfig{
		Provider:    getEnv("STORAGE_PROVIDER", "local"),
		LocalPath:   getEnv("STORAGE_LOCAL_PATH", "./uploads"),
		GCSBucket:   getEnv("STORAGE_GCS_BUCKET", ""),
		GCSProject:  getEnv("STORAGE_GCS_PROJECT", ""),
		S3Bucket:    getEnv("STORAGE_S3_BUCKET", ""),
		S3Region:    getEnv("STORAGE_S3_REGION", "us-east-1"),
		MaxFileSize: int64(getEnvInt("STORAGE_MAX_FILE_SIZE", 10485760)),
	}
}