package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestNewBillHandler verifies handler construction
func TestNewBillHandler(t *testing.T) {
	handler := NewBillHandler(nil)
	if handler == nil {
		t.Error("expected non-nil BillHandler")
	}
}

// TestPayBillInvalidInput tests PayBill with invalid JSON
func TestPayBillInvalidInput(t *testing.T) {
	handler := NewBillHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "user-1")
	c.Request, _ = http.NewRequest("POST", "/bills/pay", bytes.NewBufferString(`{invalid}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.PayBill(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["success"] != false {
		t.Error("expected success to be false")
	}
}

// TestPayBillMissingFields tests PayBill with missing required fields
func TestPayBillMissingFields(t *testing.T) {
	handler := NewBillHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "user-1")
	c.Request, _ = http.NewRequest("POST", "/bills/pay", bytes.NewBufferString(`{"biller_id": "b1"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.PayBill(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
