package services

import (
	"errors"
	"regexp"
	"strings"
	"unicode"

	"gatorpay-backend/models"

	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService handles user authentication business logic
type AuthService struct {
	db           *gorm.DB
	tokenService *TokenService
	otpService   *OTPService
}

// NewAuthService creates a new AuthService
func NewAuthService(db *gorm.DB, tokenService *TokenService, otpService *OTPService) *AuthService {
	return &AuthService{db: db, tokenService: tokenService, otpService: otpService}
}

// --- DTOs ---

// RegisterInput is the DTO for user registration
type RegisterInput struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required,min=8"`
	Username  string `json:"username" binding:"required,min=3"`
	Phone     string `json:"phone" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// LoginInput is the DTO for user login
type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// VerifyOTPInput is the DTO for OTP verification
type VerifyOTPInput struct {
	UserID  string `json:"user_id" binding:"required"`
	Code    string `json:"code" binding:"required,len=6"`
	Purpose string `json:"purpose" binding:"required"`
}

// ResendOTPInput is the DTO for resending OTP
type ResendOTPInput struct {
	UserID  string `json:"user_id" binding:"required"`
	Purpose string `json:"purpose" binding:"required"`
}

// GoogleAuthInput is the DTO for Google OAuth
type GoogleAuthInput struct {
	GoogleID string `json:"google_id" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Avatar   string `json:"avatar"`
}

// --- Responses ---

// AuthResponse is the full auth response (after OTP verified)
type AuthResponse struct {
	Token  string              `json:"token"`
	User   models.UserResponse `json:"user"`
	Wallet *models.Wallet      `json:"wallet"`
}

// OTPSentResponse is returned when OTP is sent (step 1)
type OTPSentResponse struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	Purpose string `json:"purpose"`
}
