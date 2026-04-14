package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// VirtualCard represents a user's generated virtual card
type VirtualCard struct {
	ID         string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID     string    `gorm:"type:varchar(36);index" json:"user_id"`
	CardNumber string    `json:"card_number" gorm:"uniqueIndex"`
	CVV        string    `json:"-"`
	ExpiryDate string    `json:"expiry_date"` // MM/YY
	Name       string    `json:"name"`        // Card Name tag
	IsFrozen   bool      `json:"is_frozen" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (c *VirtualCard) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}
