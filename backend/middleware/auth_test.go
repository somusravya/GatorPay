package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gatorpay-backend/config"
	"gatorpay-backend/services"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAuthMiddlewareMissingHeader(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	tokenService := services.NewTokenService(cfg)

	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	router.Use(AuthMiddleware(tokenService))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	c.Request, _ = http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, c.Request)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["message"] != "Authorization header required" {
		t.Errorf("expected 'Authorization header required', got '%v'", body["message"])
	}
}

func TestAuthMiddlewareInvalidFormat(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	tokenService := services.NewTokenService(cfg)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.Use(AuthMiddleware(tokenService))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "NotBearer token123")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["message"] != "Invalid authorization format" {
		t.Errorf("expected 'Invalid authorization format', got '%v'", body["message"])
	}
}

func TestAuthMiddlewareInvalidToken(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	tokenService := services.NewTokenService(cfg)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.Use(AuthMiddleware(tokenService))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-jwt-token")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddlewareValidToken(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	tokenService := services.NewTokenService(cfg)

	// Generate a valid token
	token, _ := tokenService.GenerateToken("user-123")

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	var capturedUserID string
	router.Use(AuthMiddleware(tokenService))
	router.GET("/test", func(c *gin.Context) {
		userID, _ := c.Get("userID")
		capturedUserID = userID.(string)
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if capturedUserID != "user-123" {
		t.Errorf("expected userID 'user-123', got '%s'", capturedUserID)
	}
}

func TestAuthMiddlewareBearerOnly(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	tokenService := services.NewTokenService(cfg)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.Use(AuthMiddleware(tokenService))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}
