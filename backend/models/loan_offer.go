package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// LoanOffer represents an available loan product
type LoanOffer struct {
	ID           string          `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name         string          `gorm:"not null" json:"name"`
	Description  string          `json:"description"`
	MinAmount    decimal.Decimal `gorm:"type:decimal(20,2)" json:"min_amount"`
	MaxAmount    decimal.Decimal `gorm:"type:decimal(20,2)" json:"max_amount"`
	InterestRate decimal.Decimal `gorm:"type:decimal(5,2)" json:"interest_rate"`
	TermMonths   int             `json:"term_months"`
	IsActive     bool            `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time       `json:"created_at"`
}

// BeforeCreate hook auto-generates UUID
func (l *LoanOffer) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	return nil
}
