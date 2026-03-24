package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	// Clear env vars to test defaults
	os.Unsetenv("PORT")
	os.Unsetenv("DB_DSN")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("CORS_ORIGINS")
	os.Unsetenv("FRONTEND_URL")
	os.Unsetenv("SMTP_HOST")
	os.Unsetenv("SMTP_PORT")
	os.Unsetenv("SMTP_USER")
	os.Unsetenv("SMTP_PASS")
	os.Unsetenv("SMTP_FROM")

	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("expected default port '8080', got '%s'", cfg.Port)
	}
	if cfg.DBDSN != "" {
		t.Errorf("expected default DB_DSN '', got '%s'", cfg.DBDSN)
	}
	if cfg.JWTSecret != "gatorpay-super-secret-key-change-in-production" {
		t.Errorf("expected default JWT secret, got '%s'", cfg.JWTSecret)
	}
	if cfg.CORSOrigins != "http://localhost:4200" {
		t.Errorf("expected default CORS origins, got '%s'", cfg.CORSOrigins)
	}
	if cfg.FrontendURL != "http://localhost:4200" {
		t.Errorf("expected default frontend URL, got '%s'", cfg.FrontendURL)
	}
	if cfg.SMTPHost != "" {
		t.Errorf("expected empty SMTP host, got '%s'", cfg.SMTPHost)
	}
	if cfg.SMTPPort != "587" {
		t.Errorf("expected default SMTP port '587', got '%s'", cfg.SMTPPort)
	}
	if cfg.SMTPFrom != "noreply@gatorpay.app" {
		t.Errorf("expected default SMTP from, got '%s'", cfg.SMTPFrom)
	}
}

func TestLoadFromEnv(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("JWT_SECRET", "my-custom-secret")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("JWT_SECRET")

	cfg := Load()

	if cfg.Port != "9090" {
		t.Errorf("expected port '9090', got '%s'", cfg.Port)
	}
	if cfg.JWTSecret != "my-custom-secret" {
		t.Errorf("expected JWT secret 'my-custom-secret', got '%s'", cfg.JWTSecret)
	}
}

func TestGetEnvWithFallback(t *testing.T) {
	os.Unsetenv("TEST_NONEXISTENT_KEY")
	result := getEnv("TEST_NONEXISTENT_KEY", "fallback_value")
	if result != "fallback_value" {
		t.Errorf("expected 'fallback_value', got '%s'", result)
	}
}

func TestGetEnvWithValue(t *testing.T) {
	os.Setenv("TEST_EXISTING_KEY", "real_value")
	defer os.Unsetenv("TEST_EXISTING_KEY")

	result := getEnv("TEST_EXISTING_KEY", "fallback_value")
	if result != "real_value" {
		t.Errorf("expected 'real_value', got '%s'", result)
	}
}
