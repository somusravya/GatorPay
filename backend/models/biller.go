package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Biller represents a bill payment provider
type Biller struct {
	ID        string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Category  string    `gorm:"not null" json:"category"`
	Icon      string    `json:"icon"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// BeforeCreate hook auto-generates UUID
func (b *Biller) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}
