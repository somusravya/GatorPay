package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Transaction type enum values
const (
	TransactionTypeDeposit    = "deposit"
	TransactionTypeWithdraw   = "withdraw"
	TransactionTypeP2PSend    = "p2p_send"
	TransactionTypeP2PReceive = "p2p_receive"
	TransactionTypeBillPay    = "bill_pay"
	TransactionTypeCashback   = "cashback"
)

// Transaction status enum values
const (
	TransactionStatusSuccess = "success"
	TransactionStatusFailed  = "failed"
)

// Transaction represents a wallet transaction record
type Transaction struct {
	ID          string          `gorm:"type:varchar(36);primaryKey" json:"id"`
	WalletID    string          `gorm:"type:varchar(36);index;not null" json:"wallet_id"`
	FromUserID  *string         `gorm:"type:varchar(36);index" json:"from_user_id"`
	ToUserID    *string         `gorm:"type:varchar(36);index" json:"to_user_id"`
	Type        string          `gorm:"not null" json:"type"`
	Amount      decimal.Decimal `gorm:"type:decimal(20,2);not null" json:"amount"`
	Description string          `json:"description"`
	Status      string          `gorm:"default:success" json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	FromUser    *User           `gorm:"foreignKey:FromUserID" json:"from_user,omitempty"`
	ToUser      *User           `gorm:"foreignKey:ToUserID" json:"to_user,omitempty"`
}

// BeforeCreate hook auto-generates UUID before inserting
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}
