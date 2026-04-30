package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Invoice struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MerchantID  uuid.UUID `gorm:"type:uuid;not null" json:"merchant_id"`
	CustomerID  uuid.UUID `gorm:"type:uuid" json:"customer_id"`
	Amount      float64   `gorm:"not null" json:"amount"`
	Currency    string    `gorm:"default:USD" json:"currency"`
	Status      string    `gorm:"default:pending" json:"status"` // pending, paid, partial, overdue, cancelled
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	PaidAt      *time.Time `json:"paid_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (i *Invoice) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}
