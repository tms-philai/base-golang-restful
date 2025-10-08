package config

type EmailConfig struct {
	Provider   string     `mapstructure:"provider"`
	From       string     `mapstructure:"from"`
	FromName   string     `mapstructure:"from_name"`
	ReplyTo    string     `mapstructure:"reply_to"`
	SMTP       SMTPConfig `mapstructure:"smtp"`
	MaxRetries int        `mapstructure:"max_retries"`
}

type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	UseTLS   bool   `mapstructure:"use_tls"`
}
