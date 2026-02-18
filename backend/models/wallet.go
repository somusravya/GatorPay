package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Wallet represents a user's financial wallet
type Wallet struct {
	ID        string          `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID    string          `gorm:"type:varchar(36);uniqueIndex;not null" json:"user_id"`
	Balance   decimal.Decimal `gorm:"type:decimal(20,2);default:0" json:"balance"`
	Currency  string          `gorm:"default:USD" json:"currency"`
	IsActive  bool            `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// BeforeCreate hook auto-generates UUID before inserting
func (w *Wallet) BeforeCreate(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = uuid.New().String()
	}
	return nil
}
