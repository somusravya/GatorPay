package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OTP purpose enum values
const (
	OTPPurposeRegister = "register"
	OTPPurposeLogin    = "login"
)

// OTPCode stores verification codes for 2-step authentication
type OTPCode struct {
	ID        string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID    string    `gorm:"type:varchar(36);index;not null" json:"user_id"`
	Code      string    `gorm:"type:varchar(6);not null" json:"-"`
	Purpose   string    `gorm:"type:varchar(20);not null" json:"purpose"` // register, login
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

// BeforeCreate hook auto-generates UUID
func (o *OTPCode) BeforeCreate(tx *gorm.DB) error {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	return nil
}

// IsExpired checks if the OTP has expired
func (o *OTPCode) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}
