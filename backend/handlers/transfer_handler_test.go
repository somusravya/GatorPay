package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestNewTransferHandler verifies handler construction
func TestNewTransferHandler(t *testing.T) {
	handler := NewTransferHandler(nil)
	if handler == nil {
		t.Error("expected non-nil TransferHandler")
	}
}

// TestSendMoneyInvalidInput tests SendMoney with invalid JSON
func TestSendMoneyInvalidInput(t *testing.T) {
	handler := NewTransferHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "user-1")
	c.Request, _ = http.NewRequest("POST", "/transfer/send", bytes.NewBufferString(`{invalid}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.SendMoney(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["success"] != false {
		t.Error("expected success to be false")
	}
}

// TestSendMoneyMissingRecipient tests SendMoney without recipient
func TestSendMoneyMissingRecipient(t *testing.T) {
	handler := NewTransferHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "user-1")
	c.Request, _ = http.NewRequest("POST", "/transfer/send", bytes.NewBufferString(`{"amount": 50}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.SendMoney(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
