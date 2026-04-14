package models

import (
	"time"

	"gorm.io/gorm"
)

// TradingAccount represents the user's KYC verified trading profile
type TradingAccount struct {
	gorm.Model
	UserID      string     `json:"user_id" gorm:"uniqueIndex"`
	DOB         time.Time  `json:"dob"`
	SSNHash     string     `json:"-"` // Hashed SSN
	SecQ1       string     `json:"sec_q1"`
	SecA1Hash   string     `json:"-"` // Hashed answer 1
	SecQ2       string     `json:"sec_q2"`
	SecA2Hash   string     `json:"-"` // Hashed answer 2
	RiskAck     bool       `json:"risk_ack"`
	Status      string     `json:"status" gorm:"default:'pending'"` // 'pending' or 'verified'
	BuyingPower float64    `json:"buying_power" gorm:"default:0.0"`
}
