package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notification represents an in-app notification
type Notification struct {
	ID        string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID    string         `gorm:"type:varchar(36);index;not null" json:"user_id"`
	Type      string         `gorm:"type:varchar(50);not null" json:"type"` // "payment", "alert", "promo", "system", "trade", "loan"
	Title     string         `gorm:"not null" json:"title"`
	Body      string         `gorm:"type:text" json:"body"`
	Icon      string         `gorm:"type:varchar(10)" json:"icon"`
	IsRead    bool           `gorm:"default:false" json:"is_read"`
	ActionURL string         `json:"action_url,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	return nil
}

// NotificationPreference stores user notification preferences
type NotificationPreference struct {
	ID              string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID          string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"user_id"`
	PaymentAlerts   bool           `gorm:"default:true" json:"payment_alerts"`
	TradeAlerts     bool           `gorm:"default:true" json:"trade_alerts"`
	LoanReminders   bool           `gorm:"default:true" json:"loan_reminders"`
	PromoOffers     bool           `gorm:"default:true" json:"promo_offers"`
	SecurityAlerts  bool           `gorm:"default:true" json:"security_alerts"`
	PriceAlerts     bool           `gorm:"default:false" json:"price_alerts"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (np *NotificationPreference) BeforeCreate(tx *gorm.DB) error {
	if np.ID == "" {
		np.ID = uuid.New().String()
	}
	return nil
}
