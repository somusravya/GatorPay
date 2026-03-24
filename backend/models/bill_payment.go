package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// BillPayment represents a bill payment record
type BillPayment struct {
	ID            string          `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID        string          `gorm:"type:varchar(36);index;not null" json:"user_id"`
	BillerID      string          `gorm:"type:varchar(36);not null" json:"biller_id"`
	AccountNumber string          `gorm:"not null" json:"account_number"`
	Amount        decimal.Decimal `gorm:"type:decimal(20,2);not null" json:"amount"`
	Status        string          `gorm:"default:success" json:"status"`
	CreatedAt     time.Time       `json:"created_at"`
	Biller        Biller          `gorm:"foreignKey:BillerID" json:"biller"`
}

// BeforeCreate hook auto-generates UUID
func (bp *BillPayment) BeforeCreate(tx *gorm.DB) error {
	if bp.ID == "" {
		bp.ID = uuid.New().String()
	}
	return nil
}
