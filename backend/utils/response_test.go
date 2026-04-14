package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSuccessResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SuccessResponse(c, http.StatusOK, "test message", map[string]string{"key": "value"})

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	if body["success"] != true {
		t.Errorf("expected success to be true, got %v", body["success"])
	}
	if body["message"] != "test message" {
		t.Errorf("expected message to be 'test message', got %v", body["message"])
	}
	if body["data"] == nil {
		t.Error("expected data to be present")
	}
}

func TestSuccessResponseWithNilData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SuccessResponse(c, http.StatusCreated, "created", nil)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["success"] != true {
		t.Error("expected success to be true")
	}
}

func TestErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ErrorResponse(c, http.StatusBadRequest, "bad request")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)

	if body["success"] != false {
		t.Errorf("expected success to be false, got %v", body["success"])
	}
	if body["message"] != "bad request" {
		t.Errorf("expected message to be 'bad request', got %v", body["message"])
	}
	if body["data"] != nil {
		t.Error("expected data to be nil")
	}
}

func TestErrorResponseUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ErrorResponse(c, http.StatusUnauthorized, "unauthorized")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestErrorResponseInternalServer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ErrorResponse(c, http.StatusInternalServerError, "internal error")

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
