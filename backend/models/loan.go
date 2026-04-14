package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Loan represents an issued loan
type Loan struct {
	ID              string          `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID          string          `gorm:"type:varchar(36);index" json:"user_id"`
	Amount          decimal.Decimal `gorm:"type:decimal(20,2)" json:"amount"`
	TermMonths      int             `json:"term_months"`
	InterestRate    decimal.Decimal `gorm:"type:decimal(5,2)" json:"interest_rate"`
	EMI             decimal.Decimal `gorm:"type:decimal(20,2)" json:"emi"`
	TotalPayable    decimal.Decimal `gorm:"type:decimal(20,2)" json:"total_payable"`
	AmountPaid      decimal.Decimal `gorm:"type:decimal(20,2);default:0" json:"amount_paid"`
	RemainingAmount decimal.Decimal `gorm:"type:decimal(20,2)" json:"remaining_amount"`
	Status          string          `json:"status" gorm:"default:'active'"` // active, paid
	NextPaymentDate time.Time       `json:"next_payment_date"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

func (l *Loan) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	return nil
}
