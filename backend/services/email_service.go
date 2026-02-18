package services

import (
	"fmt"
	"log"
	"net/smtp"
	"strconv"

	"gatorpay-backend/config"
)

// EmailService handles sending emails (OTP codes, etc.)
type EmailService struct {
	cfg *config.Config
}

// NewEmailService creates a new EmailService
func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{cfg: cfg}
}

// IsConfigured returns true if SMTP is set up
func (s *EmailService) IsConfigured() bool {
	return s.cfg.SMTPHost != "" && s.cfg.SMTPUser != ""
}

// SendOTP sends an OTP code via email, or prints to console if SMTP is not configured
func (s *EmailService) SendOTP(toEmail, code, purpose string) error {
	if !s.IsConfigured() {
		// Dev mode: print to console
		log.Printf("📧 [DEV] OTP for %s (%s): %s", toEmail, purpose, code)
		return nil
	}

	subject := "GatorPay - Verification Code"
	if purpose == "login" {
		subject = "GatorPay - Login Verification"
	}

	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; background-color: #0f172a; color: #f1f5f9; padding: 40px;">
  <div style="max-width: 480px; margin: 0 auto; background: #1e293b; border-radius: 16px; padding: 40px; text-align: center;">
    <h1 style="color: #818cf8; margin-bottom: 8px;">🐊 GatorPay</h1>
    <p style="color: #94a3b8; margin-bottom: 32px;">Your verification code</p>
    <div style="background: #0f172a; border-radius: 12px; padding: 24px; margin-bottom: 24px;">
      <span style="font-size: 36px; font-weight: bold; letter-spacing: 12px; color: #818cf8;">%s</span>
    </div>
    <p style="color: #94a3b8; font-size: 14px;">This code expires in 5 minutes.<br>If you did not request this, please ignore this email.</p>
  </div>
</body>
</html>`, code)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		s.cfg.SMTPFrom, toEmail, subject, body)

	port, _ := strconv.Atoi(s.cfg.SMTPPort)
	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, port)
	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, s.cfg.SMTPHost)

	err := smtp.SendMail(addr, auth, s.cfg.SMTPFrom, []string{toEmail}, []byte(msg))
	if err != nil {
		log.Printf("❌ Failed to send email to %s: %v", toEmail, err)
		// Fallback to console in case of SMTP failure
		log.Printf("📧 [FALLBACK] OTP for %s (%s): %s", toEmail, purpose, code)
		return nil // Don't fail the request even if email fails
	}

	log.Printf("📧 OTP sent to %s for %s", toEmail, purpose)
	return nil
}
