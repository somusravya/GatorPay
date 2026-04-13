package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestNewWalletHandler verifies handler construction
func TestNewWalletHandler(t *testing.T) {
	handler := NewWalletHandler(nil)
	if handler == nil {
		t.Error("expected non-nil WalletHandler")
	}
}

// TestAddMoneyInvalidInput tests AddMoney with invalid JSON
func TestAddMoneyInvalidInput(t *testing.T) {
	handler := NewWalletHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "user-1")
	c.Request, _ = http.NewRequest("POST", "/wallet/add", bytes.NewBufferString(`{invalid}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.AddMoney(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["success"] != false {
		t.Error("expected success to be false")
	}
}

// TestAddMoneyMissingFields tests AddMoney with missing required fields
func TestAddMoneyMissingFields(t *testing.T) {
	handler := NewWalletHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "user-1")
	c.Request, _ = http.NewRequest("POST", "/wallet/add", bytes.NewBufferString(`{"amount": 100}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.AddMoney(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestWithdrawInvalidInput tests Withdraw with invalid JSON
func TestWithdrawInvalidInput(t *testing.T) {
	handler := NewWalletHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "user-1")
	c.Request, _ = http.NewRequest("POST", "/wallet/withdraw", bytes.NewBufferString(`bad-json`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Withdraw(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestWithdrawMissingBankAccount tests Withdraw without bank_account
func TestWithdrawMissingBankAccount(t *testing.T) {
	handler := NewWalletHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "user-1")
	c.Request, _ = http.NewRequest("POST", "/wallet/withdraw", bytes.NewBufferString(`{"amount": 50}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Withdraw(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
