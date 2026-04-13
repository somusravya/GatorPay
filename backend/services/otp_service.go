package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"gatorpay-backend/models"

	"gorm.io/gorm"
)

// OTPService handles OTP generation, storage, and verification
type OTPService struct {
	db           *gorm.DB
	emailService *EmailService
}

// NewOTPService creates a new OTPService
func NewOTPService(db *gorm.DB, emailService *EmailService) *OTPService {
	return &OTPService{db: db, emailService: emailService}
}

// GenerateAndSend creates a new 6-digit OTP, stores it, and sends via email
func (s *OTPService) GenerateAndSend(userID, email, purpose string) error {
	// Invalidate any existing unused OTPs for the same user+purpose
	s.db.Model(&models.OTPCode{}).
		Where("user_id = ? AND purpose = ? AND used = false", userID, purpose).
		Update("used", true)

	// Generate secure 6-digit code
	code, err := generateSecureCode()
	if err != nil {
		return errors.New("failed to generate OTP")
	}

	// Store OTP with 5-minute expiry
	otp := models.OTPCode{
		UserID:    userID,
		Code:      code,
		Purpose:   purpose,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Used:      false,
	}
	if err := s.db.Create(&otp).Error; err != nil {
		return errors.New("failed to store OTP")
	}

	// Send via email (or console)
	return s.emailService.SendOTP(email, code, purpose)
}

// Verify checks the OTP code for a given user and purpose
func (s *OTPService) Verify(userID, code, purpose string) error {
	var otp models.OTPCode
	err := s.db.Where(
		"user_id = ? AND purpose = ? AND used = false AND verified = false",
		userID, purpose,
	).Order("created_at DESC").First(&otp).Error

	if err != nil {
		return errors.New("invalid verification code")
	}

	if otp.IsExpired() {
		s.db.Model(&otp).Update("used", true)
		return errors.New("verification code has expired, please request a new one")
	}

	// Check max attempts
	if otp.Attempts >= 3 {
		s.db.Model(&otp).Update("used", true)
		return errors.New("maximum attempts exceeded, please request a new code")
	}

	// Increment attempts
	s.db.Model(&otp).Update("attempts", otp.Attempts+1)

	// Verify code
	if otp.Code != code {
		return errors.New("invalid verification code")
	}

	// Mark as verified and used
	s.db.Model(&otp).Updates(map[string]interface{}{"used": true, "verified": true})
	return nil
}

// GenerateOTP creates a standalone OTP without sending email (for transfer verification)
func (s *OTPService) GenerateOTP(userID, purpose string) (string, error) {
	// Invalidate existing OTPs
	s.db.Model(&models.OTPCode{}).
		Where("user_id = ? AND purpose = ? AND used = false", userID, purpose).
		Update("used", true)

	code, err := generateSecureCode()
	if err != nil {
		return "", errors.New("failed to generate OTP")
	}

	otp := models.OTPCode{
		UserID:    userID,
		Code:      code,
		Purpose:   purpose,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Used:      false,
	}
	if err := s.db.Create(&otp).Error; err != nil {
		return "", errors.New("failed to store OTP")
	}

	return code, nil
}

// generateSecureCode returns a cryptographically random 6-digit string
func generateSecureCode() (string, error) {
	b := make([]byte, 3) // 3 bytes = enough entropy for 6 digits
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// Convert to 6-digit number (000000-999999)
	num := (int(b[0])<<16 | int(b[1])<<8 | int(b[2])) % 1000000
	return fmt.Sprintf("%06d", num), nil
}
