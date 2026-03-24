package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// TestNewAuthHandler verifies handler construction
func TestNewAuthHandler(t *testing.T) {
	handler := NewAuthHandler(nil)
	if handler == nil {
		t.Error("expected non-nil AuthHandler")
	}
}

// TestRegisterInvalidInput tests Register with invalid JSON body
func TestRegisterInvalidInput(t *testing.T) {
	handler := NewAuthHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/auth/register", bytes.NewBufferString(`{invalid-json}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Register(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["success"] != false {
		t.Error("expected success to be false")
	}
}

// TestRegisterMissingFields tests Register with missing required fields
func TestRegisterMissingFields(t *testing.T) {
	handler := NewAuthHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/auth/register", bytes.NewBufferString(`{"email": "test@test.com"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Register(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestLoginInvalidInput tests Login with invalid JSON
func TestLoginInvalidInput(t *testing.T) {
	handler := NewAuthHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/auth/login", bytes.NewBufferString(`not-json`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestLoginMissingPassword tests Login without password field
func TestLoginMissingPassword(t *testing.T) {
	handler := NewAuthHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/auth/login", bytes.NewBufferString(`{"email": "test@test.com"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestVerifyOTPInvalidInput tests VerifyOTP with invalid JSON
func TestVerifyOTPInvalidInput(t *testing.T) {
	handler := NewAuthHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/auth/verify-otp", bytes.NewBufferString(`{invalid}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.VerifyOTP(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestResendOTPInvalidInput tests ResendOTP with invalid JSON
func TestResendOTPInvalidInput(t *testing.T) {
	handler := NewAuthHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/auth/resend-otp", bytes.NewBufferString(`{invalid}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.ResendOTP(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestGetMeNotAuthenticated tests GetMe without userID in context
func TestGetMeNotAuthenticated(t *testing.T) {
	handler := NewAuthHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/auth/me", nil)

	handler.GetMe(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["message"] != "User not authenticated" {
		t.Errorf("expected 'User not authenticated', got '%v'", body["message"])
	}
}

// TestGoogleAuthInvalidInput tests GoogleAuth with invalid JSON
func TestGoogleAuthInvalidInput(t *testing.T) {
	handler := NewAuthHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/auth/google", bytes.NewBufferString(`{bad}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.GoogleAuth(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
