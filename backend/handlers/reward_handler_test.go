package handlers

import (
	"testing"
)

// TestNewRewardHandler verifies handler construction
func TestNewRewardHandler(t *testing.T) {
	handler := NewRewardHandler(nil)
	if handler == nil {
		t.Error("expected non-nil RewardHandler")
	}
}
