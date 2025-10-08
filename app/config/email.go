package config

type EmailConfig struct {
	Provider string
	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string
	From     string
	FromName string
	Enabled  bool
}

func loadEmailConfig() EmailConfig {
	return EmailConfig{
		Provider: getEnv("EMAIL_PROVIDER", "smtp"),
		SMTPHost: getEnv("EMAIL_SMTP_HOST", "localhost"),
		SMTPPort: getEnvInt("EMAIL_SMTP_PORT", 587),
		SMTPUser: getEnv("EMAIL_SMTP_USER", ""),
		SMTPPass: getEnv("EMAIL_SMTP_PASS", ""),
		From:     getEnv("EMAIL_FROM", "noreply@example.com"),
		FromName: getEnv("EMAIL_FROM_NAME", "Base Golang API"),
		Enabled:  getEnvBool("EMAIL_ENABLED", false),
	}
}