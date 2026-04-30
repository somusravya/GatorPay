package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FraudAlert represents a flagged fraud event
type FraudAlert struct {
	ID            string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID        string         `gorm:"type:varchar(36);index;not null" json:"user_id"`
	TransactionID string         `gorm:"type:varchar(36)" json:"transaction_id"`
	RiskScore     float64        `json:"risk_score"` // 0-100
	Type          string         `gorm:"type:varchar(50)" json:"type"`   // "velocity", "geo_anomaly", "amount_spike", "suspicious_merchant"
	Status        string         `gorm:"type:varchar(20);default:pending" json:"status"` // "pending", "reviewed", "dismissed", "confirmed"
	Description   string         `json:"description"`
	Details       string         `gorm:"type:text" json:"details"`
	ReviewedBy    string         `gorm:"type:varchar(36)" json:"reviewed_by,omitempty"`
	ReviewedAt    *time.Time     `json:"reviewed_at,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	User          *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (f *FraudAlert) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	return nil
}

// RiskEvent represents an audit log entry for risk engine evaluations
type RiskEvent struct {
	ID            string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID        string         `gorm:"type:varchar(36);index;not null" json:"user_id"`
	TransactionID string         `gorm:"type:varchar(36)" json:"transaction_id"`
	EventType     string         `gorm:"type:varchar(50)" json:"event_type"`
	RiskScore     float64        `json:"risk_score"`
	Factors       string         `gorm:"type:text" json:"factors"` // JSON string of risk factors
	Action        string         `gorm:"type:varchar(30)" json:"action"` // "allow", "flag", "block"
	CreatedAt     time.Time      `json:"created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (r *RiskEvent) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}
