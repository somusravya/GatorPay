package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Subscription represents a tracked recurring charge
type Subscription struct {
	ID          string          `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID      string          `gorm:"type:varchar(36);index;not null" json:"user_id"`
	Name        string          `gorm:"not null" json:"name"`
	Category    string          `gorm:"type:varchar(50)" json:"category"` // "streaming", "software", "fitness", etc.
	Amount      decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"amount"`
	Frequency   string          `gorm:"type:varchar(20);default:monthly" json:"frequency"` // "monthly", "yearly", "weekly"
	NextRenewal *time.Time      `json:"next_renewal"`
	Icon        string          `gorm:"type:varchar(10)" json:"icon"`
	Color       string          `gorm:"type:varchar(20)" json:"color"`
	AutoPay     bool            `gorm:"default:false" json:"auto_pay"`
	Status      string          `gorm:"type:varchar(20);default:active" json:"status"` // "active", "paused", "cancelled"
	Provider    string          `json:"provider"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

// TrackSubscriptionRequest is the request body for tracking a subscription
type TrackSubscriptionRequest struct {
	Name      string  `json:"name" binding:"required"`
	Category  string  `json:"category" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	Frequency string  `json:"frequency" binding:"required"`
	Provider  string  `json:"provider"`
	Icon      string  `json:"icon"`
	Color     string  `json:"color"`
}

// AutoPayRequest is the request for setting auto-pay
type AutoPayRequest struct {
	SubscriptionID string `json:"subscription_id" binding:"required"`
	AutoPay        bool   `json:"auto_pay"`
}
