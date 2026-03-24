package services

import (
	"testing"
)

func TestValidateEmailValid(t *testing.T) {
	validEmails := []string{
		"user@example.com",
		"test@gator.edu",
		"john.doe@university.edu",
		"user+tag@gmail.com",
		"a@b.co",
	}

	for _, email := range validEmails {
		if err := validateEmail(email); err != nil {
			t.Errorf("expected email '%s' to be valid, got error: %v", email, err)
		}
	}
}

func TestValidateEmailInvalid(t *testing.T) {
	invalidEmails := []string{
		"",
		"invalid",
		"@example.com",
		"user@",
		"user@.com",
		"user@com",
		"user space@example.com",
	}

	for _, email := range invalidEmails {
		if err := validateEmail(email); err == nil {
			t.Errorf("expected email '%s' to be invalid, but got no error", email)
		}
	}
}

func TestValidatePhoneValid(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"3521234567", "3521234567"},
		{"(352) 123-4567", "3521234567"},
		{"352-123-4567", "3521234567"},
		{"352.123.4567", "3521234567"},
		{"352 123 4567", "3521234567"}, // spaces stripped, 10 digits
	}

	for _, tt := range tests {
		cleaned, err := validatePhone(tt.input)
		if err != nil {
			t.Errorf("expected phone '%s' to be valid, got error: %v", tt.input, err)
		}
		if cleaned != tt.expected {
			t.Errorf("for input '%s', expected '%s', got '%s'", tt.input, tt.expected, cleaned)
		}
	}
}

func TestValidatePhoneInvalid(t *testing.T) {
	invalidPhones := []string{
		"",
		"12345",
		"12345678901", // 11 digits
		"abcdefghij",
		"123",
	}

	for _, phone := range invalidPhones {
		_, err := validatePhone(phone)
		if err == nil {
			t.Errorf("expected phone '%s' to be invalid, but got no error", phone)
		}
	}
}

func TestMaskEmail(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"john@gator.edu", "j***@gator.edu"},
		{"a@example.com", "a***@example.com"},
		{"test.user@ufl.edu", "t***@ufl.edu"},
	}

	for _, tt := range tests {
		result := maskEmail(tt.input)
		if result != tt.expected {
			t.Errorf("maskEmail('%s') = '%s', expected '%s'", tt.input, result, tt.expected)
		}
	}
}

func TestMaskEmailInvalidFormat(t *testing.T) {
	result := maskEmail("invalid-email")
	if result != "invalid-email" {
		t.Errorf("expected 'invalid-email' returned as-is, got '%s'", result)
	}
}

func TestMaskEmailEmptyName(t *testing.T) {
	result := maskEmail("@gator.edu")
	// With empty name part, it should still handle gracefully
	if result == "" {
		t.Error("expected non-empty result for '@gator.edu'")
	}
}

func TestNewAuthService(t *testing.T) {
	// Just verifying constructor doesn't panic with nil args
	// In real usage, db/token/otp would be non-nil
	service := NewAuthService(nil, nil, nil)
	if service == nil {
		t.Error("expected non-nil AuthService")
	}
}
