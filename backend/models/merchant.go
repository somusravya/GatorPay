package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Merchant struct {
	ID           string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID       string    `gorm:"type:varchar(36);uniqueIndex" json:"user_id"`
	BusinessName string    `json:"business_name"`
	Category     string    `json:"category"`
	Status       string    `json:"status" gorm:"default:'active'"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (m *Merchant) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}

type MerchantQRCode struct {
	ID         string          `gorm:"type:varchar(36);primaryKey" json:"id"`
	MerchantID string          `gorm:"type:varchar(36);index" json:"merchant_id"`
	CodeString string          `json:"code_string" gorm:"uniqueIndex"`
	Base64PNG  string          `json:"base64_png"`
	Amount     decimal.Decimal `gorm:"type:decimal(20,2)" json:"amount,omitempty"`
	IsDynamic  bool            `json:"is_dynamic" gorm:"default:false"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

func (q *MerchantQRCode) BeforeCreate(tx *gorm.DB) error {
	if q.ID == "" {
		q.ID = uuid.New().String()
	}
	return nil
}
