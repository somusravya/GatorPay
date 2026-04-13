package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Merchant is a user-registered business that can receive QR payments.
type Merchant struct {
	ID           string `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID       string `gorm:"type:varchar(36);uniqueIndex;not null" json:"user_id"`
	BusinessName string `json:"business_name"`
	Category     string `json:"category"`
	Status       string `json:"status"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (m *Merchant) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}

// MerchantQRCode stores a generated pay-by-QR payload for a merchant.
type MerchantQRCode struct {
	ID         string          `gorm:"type:varchar(36);primaryKey" json:"id"`
	MerchantID string          `gorm:"type:varchar(36);index;not null" json:"merchant_id"`
	CodeString string          `gorm:"uniqueIndex;not null" json:"code_string"`
	Base64PNG  string          `json:"qr_image"`
	Amount     decimal.Decimal `gorm:"type:decimal(20,2)" json:"amount"`
	IsDynamic  bool            `json:"is_dynamic"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (q *MerchantQRCode) BeforeCreate(tx *gorm.DB) error {
	if q.ID == "" {
		q.ID = uuid.New().String()
	}
	return nil
}

// VirtualCard represents a user's generated virtual card
type VirtualCard struct {
	ID         string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID     string    `gorm:"type:varchar(36);index" json:"user_id"`
	CardNumber string    `json:"card_number" gorm:"uniqueIndex"`
	CVV        string    `json:"-"`
	ExpiryDate string    `json:"expiry_date"` // MM/YY
	Name       string    `json:"name"`        // Card Name tag
	IsFrozen   bool      `json:"is_frozen" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (c *VirtualCard) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}
