package services

import (
	"testing"
	"time"

	"gatorpay-backend/config"
)

func TestGenerateToken(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret-key-for-testing"}
	ts := NewTokenService(cfg)

	token, err := ts.GenerateToken("user-123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
}

func TestValidateToken(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret-key-for-testing"}
	ts := NewTokenService(cfg)

	token, _ := ts.GenerateToken("user-456")

	userID, err := ts.ValidateToken(token)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if userID != "user-456" {
		t.Errorf("expected user ID 'user-456', got '%s'", userID)
	}
}

func TestValidateTokenInvalidSignature(t *testing.T) {
	cfg1 := &config.Config{JWTSecret: "secret-1"}
	cfg2 := &config.Config{JWTSecret: "secret-2"}

	ts1 := NewTokenService(cfg1)
	ts2 := NewTokenService(cfg2)

	token, _ := ts1.GenerateToken("user-789")

	_, err := ts2.ValidateToken(token)
	if err == nil {
		t.Error("expected error for invalid signature, got nil")
	}
}

func TestValidateTokenMalformed(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	ts := NewTokenService(cfg)

	_, err := ts.ValidateToken("this.is.not.a.valid.jwt")
	if err == nil {
		t.Error("expected error for malformed token, got nil")
	}
}

func TestValidateTokenEmptyString(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	ts := NewTokenService(cfg)

	_, err := ts.ValidateToken("")
	if err == nil {
		t.Error("expected error for empty token, got nil")
	}
}

func TestGenerateTokenDifferentUserIDs(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	ts := NewTokenService(cfg)

	token1, _ := ts.GenerateToken("user-1")
	token2, _ := ts.GenerateToken("user-2")

	if token1 == token2 {
		t.Error("expected different tokens for different user IDs")
	}

	userID1, _ := ts.ValidateToken(token1)
	userID2, _ := ts.ValidateToken(token2)

	if userID1 != "user-1" {
		t.Errorf("expected 'user-1', got '%s'", userID1)
	}
	if userID2 != "user-2" {
		t.Errorf("expected 'user-2', got '%s'", userID2)
	}
}

func TestNewTokenService(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test"}
	ts := NewTokenService(cfg)
	if ts == nil {
		t.Error("expected non-nil TokenService")
	}
}

func TestTokenServiceRoundTrip(t *testing.T) {
	cfg := &config.Config{JWTSecret: "roundtrip-secret"}
	ts := NewTokenService(cfg)

	// Generate and validate in quick succession
	start := time.Now()
	token, _ := ts.GenerateToken("test-user")
	userID, err := ts.ValidateToken(token)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("round-trip validation failed: %v", err)
	}
	if userID != "test-user" {
		t.Errorf("expected 'test-user', got '%s'", userID)
	}
	if elapsed > 1*time.Second {
		t.Errorf("token round-trip took too long: %v", elapsed)
	}
}
