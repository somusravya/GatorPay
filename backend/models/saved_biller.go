package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SavedBiller represents a user's saved biller for quick payments
type SavedBiller struct {
	ID            string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID        string    `gorm:"type:varchar(36);index;not null" json:"user_id"`
	BillerID      string    `gorm:"type:varchar(36);not null" json:"biller_id"`
	AccountNumber string    `gorm:"not null" json:"account_number"`
	Nickname      string    `json:"nickname"`
	CreatedAt     time.Time `json:"created_at"`
	Biller        Biller    `gorm:"foreignKey:BillerID" json:"biller"`
}

// BeforeCreate hook auto-generates UUID
func (s *SavedBiller) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}
