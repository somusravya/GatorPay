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

// --- Validation helpers ---

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func validateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func validatePhone(phone string) (string, error) {
	// Strip non-digit characters
	var digits strings.Builder
	for _, r := range phone {
		if unicode.IsDigit(r) {
			digits.WriteRune(r)
		}
	}
	cleaned := digits.String()
	if len(cleaned) != 10 {
		return "", errors.New("phone number must be exactly 10 digits")
	}
	return cleaned, nil
}

// --- Auth methods ---

// Register creates a new user with wallet, then sends OTP (no JWT returned yet)
func (s *AuthService) Register(input RegisterInput) (*OTPSentResponse, error) {
	// Validate email format
	if err := validateEmail(input.Email); err != nil {
		return nil, err
	}

	// Validate phone
	cleanPhone, err := validatePhone(input.Phone)
	if err != nil {
		return nil, err
	}
	input.Phone = cleanPhone

	// Check email uniqueness
	var existingUser models.User
	if err := s.db.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already registered")
	}

	// Check username uniqueness
	if err := s.db.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already taken")
	}

	// Check phone uniqueness
	if err := s.db.Where("phone = ?", input.Phone).First(&existingUser).Error; err == nil {
		return nil, errors.New("phone number already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	var user models.User

	// Create user and wallet in a DB transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		user = models.User{
			Email:         input.Email,
			Username:      input.Username,
			Phone:         input.Phone,
			PasswordHash:  string(hashedPassword),
			FirstName:     input.FirstName,
			LastName:      input.LastName,
			AuthProvider:  models.AuthProviderLocal,
			EmailVerified: false,
			KYCStatus:     models.KYCPending,
			CreditScore:   650,
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		wallet := models.Wallet{
			UserID:   user.ID,
			Balance:  decimal.NewFromInt(0),
			Currency: "USD",
			IsActive: true,
		}
		return tx.Create(&wallet).Error
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique") {
			return nil, errors.New("user with this information already exists")
		}
		return nil, errors.New("failed to create user")
	}

	// Send OTP for registration verification
	if err := s.otpService.GenerateAndSend(user.ID, user.Email, models.OTPPurposeRegister); err != nil {
		return nil, errors.New("failed to send verification code")
	}

	return &OTPSentResponse{
		UserID:  user.ID,
		Email:   maskEmail(user.Email),
		Purpose: models.OTPPurposeRegister,
	}, nil
}

// Login validates credentials then sends OTP (no JWT returned yet)
func (s *AuthService) Login(input LoginInput) (*OTPSentResponse, error) {
	// Validate email format
	if err := validateEmail(input.Email); err != nil {
		return nil, err
	}

	var user models.User
	if err := s.db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Compare bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Send OTP for login verification
	if err := s.otpService.GenerateAndSend(user.ID, user.Email, models.OTPPurposeLogin); err != nil {
		return nil, errors.New("failed to send verification code")
	}

	return &OTPSentResponse{
		UserID:  user.ID,
		Email:   maskEmail(user.Email),
		Purpose: models.OTPPurposeLogin,
	}, nil
}

// VerifyOTP verifies the OTP code and completes authentication
func (s *AuthService) VerifyOTP(input VerifyOTPInput) (*AuthResponse, error) {
	// Verify OTP
	if err := s.otpService.Verify(input.UserID, input.Code, input.Purpose); err != nil {
		return nil, err
	}

	// Load user with wallet
	var user models.User
	if err := s.db.Preload("Wallet").Where("id = ?", input.UserID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// If this is a registration OTP, mark email as verified
	if input.Purpose == models.OTPPurposeRegister {
		user.EmailVerified = true
		s.db.Save(&user)
	}

	// Generate JWT
	token, err := s.tokenService.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResponse{
		Token:  token,
		User:   user.ToResponse(),
		Wallet: user.Wallet,
	}, nil
}

// ResendOTP resends a new OTP code
func (s *AuthService) ResendOTP(input ResendOTPInput) (*OTPSentResponse, error) {
	var user models.User
	if err := s.db.Where("id = ?", input.UserID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	if err := s.otpService.GenerateAndSend(user.ID, user.Email, input.Purpose); err != nil {
		return nil, errors.New("failed to send verification code")
	}

	return &OTPSentResponse{
		UserID:  user.ID,
		Email:   maskEmail(user.Email),
		Purpose: input.Purpose,
	}, nil
}

// GoogleAuth handles Google OAuth login/registration (no OTP required)
func (s *AuthService) GoogleAuth(input GoogleAuthInput) (*AuthResponse, error) {
	var user models.User

	// Check if GoogleID already exists → login
	if err := s.db.Preload("Wallet").Where("google_id = ?", input.GoogleID).First(&user).Error; err == nil {
		token, err := s.tokenService.GenerateToken(user.ID)
		if err != nil {
			return nil, errors.New("failed to generate token")
		}
		return &AuthResponse{
			Token:  token,
			User:   user.ToResponse(),
			Wallet: user.Wallet,
		}, nil
	}

	// Check if email exists → link account
	if err := s.db.Preload("Wallet").Where("email = ?", input.Email).First(&user).Error; err == nil {
		user.GoogleID = input.GoogleID
		user.AuthProvider = models.AuthProviderGoogle
		user.EmailVerified = true
		if input.Avatar != "" {
			user.AvatarURL = input.Avatar
		}
		s.db.Save(&user)

		token, err := s.tokenService.GenerateToken(user.ID)
		if err != nil {
			return nil, errors.New("failed to generate token")
		}
		return &AuthResponse{
			Token:  token,
			User:   user.ToResponse(),
			Wallet: user.Wallet,
		}, nil
	}

	// Create new user
	nameParts := strings.SplitN(input.Name, " ", 2)
	firstName := nameParts[0]
	lastName := ""
	if len(nameParts) > 1 {
		lastName = nameParts[1]
	}

	var wallet models.Wallet
	err := s.db.Transaction(func(tx *gorm.DB) error {
		user = models.User{
			Email:         input.Email,
			Username:      strings.Split(input.Email, "@")[0],
			Phone:         "",
			FirstName:     firstName,
			LastName:      lastName,
			AvatarURL:     input.Avatar,
			AuthProvider:  models.AuthProviderGoogle,
			GoogleID:      input.GoogleID,
			EmailVerified: true,
			KYCStatus:     models.KYCPending,
			CreditScore:   650,
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		wallet = models.Wallet{
			UserID:   user.ID,
			Balance:  decimal.NewFromInt(0),
			Currency: "USD",
			IsActive: true,
		}
		return tx.Create(&wallet).Error
	})
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	token, err := s.tokenService.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResponse{
		Token:  token,
		User:   user.ToResponse(),
		Wallet: &wallet,
	}, nil
}

// GetMe returns the current user with their wallet
func (s *AuthService) GetMe(userID string) (*AuthResponse, error) {
	var user models.User
	if err := s.db.Preload("Wallet").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return &AuthResponse{
		User:   user.ToResponse(),
		Wallet: user.Wallet,
	}, nil
}

// maskEmail masks the middle of an email for display (e.g. j***@gator.edu)
func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}
	name := parts[0]
	if len(name) <= 1 {
		return name + "***@" + parts[1]
	}
	return string(name[0]) + "***@" + parts[1]
}
