package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Loan is an active user loan instance tied to a LoanOffer.
type Loan struct {
	ID                   string          `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID               string          `gorm:"type:varchar(36);index;not null" json:"user_id"`
	OfferID              string          `gorm:"type:varchar(36);not null" json:"offer_id"`
	Principal            decimal.Decimal `gorm:"type:decimal(20,2);not null" json:"principal"`
	TermMonths           int             `json:"term_months"`
	InterestRate         decimal.Decimal `gorm:"type:decimal(5,2)" json:"interest_rate"`
	OutstandingBalance   decimal.Decimal `gorm:"type:decimal(20,2);not null" json:"outstanding_balance"`
	MonthlyPayment       decimal.Decimal `gorm:"type:decimal(20,2);not null" json:"monthly_payment"`
	Status               string          `gorm:"default:active" json:"status"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

func (l *Loan) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	return nil
}
