package config

import "os"

// Config holds all configuration for the application
type Config struct {
	Port        string
	DBDSN       string
	JWTSecret   string
	CORSOrigins string
	FrontendURL string

	// SMTP settings
	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string
	SMTPFrom string
}

// Load reads configuration from environment variables with sensible defaults
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DBDSN:       getEnv("DB_DSN", ""),
		JWTSecret:   getEnv("JWT_SECRET", "gatorpay-super-secret-key-change-in-production"),
		CORSOrigins: getEnv("CORS_ORIGINS", "http://localhost:4200"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:4200"),
		SMTPHost:    getEnv("SMTP_HOST", ""),
		SMTPPort:    getEnv("SMTP_PORT", "587"),
		SMTPUser:    getEnv("SMTP_USER", ""),
		SMTPPass:    getEnv("SMTP_PASS", ""),
		SMTPFrom:    getEnv("SMTP_FROM", "noreply@gatorpay.app"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
