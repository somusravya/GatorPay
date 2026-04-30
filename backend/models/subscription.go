package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscription struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Name       string    `gorm:"not null" json:"name"`
	Provider   string    `json:"provider"`
	Amount     float64   `gorm:"not null" json:"amount"`
	Currency   string    `gorm:"default:USD" json:"currency"`
	Frequency  string    `gorm:"default:monthly" json:"frequency"` // monthly, yearly, weekly
	Status     string    `gorm:"default:active" json:"status"`     // active, cancelled, detected, paused
	NextBill   time.Time `json:"next_bill"`
	Category   string    `json:"category"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
