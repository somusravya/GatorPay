package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// KYC status enum values
const (
	KYCPending  = "pending"
	KYCVerified = "verified"
	KYCRejected = "rejected"
)

// Auth provider enum values
const (
	AuthProviderLocal  = "local"
	AuthProviderGoogle = "google"
)

// User represents a GatorPay user account
type User struct {
	ID            string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Email         string         `gorm:"uniqueIndex;not null" json:"email"`
	Username      string         `gorm:"uniqueIndex;not null" json:"username"`
	Phone         string         `gorm:"uniqueIndex;not null" json:"phone"`
	PasswordHash  string         `gorm:"column:password_hash" json:"-"`
	FirstName     string         `gorm:"not null" json:"first_name"`
	LastName      string         `gorm:"not null" json:"last_name"`
	AvatarURL     string         `json:"avatar_url"`
	AuthProvider  string         `gorm:"default:local" json:"auth_provider"`
	GoogleID      string         `gorm:"index" json:"google_id,omitempty"`
	EmailVerified bool           `gorm:"default:false" json:"email_verified"`
	KYCStatus     string         `gorm:"default:pending" json:"kyc_status"`
	CreditScore   int            `gorm:"default:650" json:"credit_score"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Wallet        *Wallet        `gorm:"foreignKey:UserID" json:"wallet,omitempty"`
}

// BeforeCreate hook auto-generates UUID before inserting
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// UserResponse is a sanitized user without sensitive fields
type UserResponse struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	Phone         string    `json:"phone"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	AvatarURL     string    `json:"avatar_url"`
	AuthProvider  string    `json:"auth_provider"`
	EmailVerified bool      `json:"email_verified"`
	KYCStatus     string    `json:"kyc_status"`
	CreditScore   int       `json:"credit_score"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ToResponse converts User to a sanitized response (no password)
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:            u.ID,
		Email:         u.Email,
		Username:      u.Username,
		Phone:         u.Phone,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		AvatarURL:     u.AvatarURL,
		AuthProvider:  u.AuthProvider,
		EmailVerified: u.EmailVerified,
		KYCStatus:     u.KYCStatus,
		CreditScore:   u.CreditScore,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}
