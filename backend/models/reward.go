package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Reward type constants
const (
	RewardTypeCashback = "cashback"
)

// Reward represents a cashback or reward record
type Reward struct {
	ID            string          `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID        string          `gorm:"type:varchar(36);index;not null" json:"user_id"`
	Type          string          `gorm:"not null" json:"type"`
	Amount        decimal.Decimal `gorm:"type:decimal(20,2);not null" json:"amount"`
	Points        int             `gorm:"default:0" json:"points"`
	TransactionID string          `gorm:"type:varchar(36)" json:"transaction_id"`
	Description   string          `json:"description"`
	CreatedAt     time.Time       `json:"created_at"`
}

// BeforeCreate hook auto-generates UUID
func (r *Reward) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}
