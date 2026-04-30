package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// BudgetGoal represents a savings goal
type BudgetGoal struct {
	ID            string          `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID        string          `gorm:"type:varchar(36);index;not null" json:"user_id"`
	Name          string          `gorm:"not null" json:"name"`
	Category      string          `gorm:"type:varchar(50)" json:"category"` // "emergency", "vacation", "education", etc.
	TargetAmount  decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"target_amount"`
	CurrentAmount decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"current_amount"`
	Deadline      *time.Time      `json:"deadline,omitempty"`
	Icon          string          `gorm:"type:varchar(10)" json:"icon"`
	Color         string          `gorm:"type:varchar(20)" json:"color"`
	Status        string          `gorm:"type:varchar(20);default:active" json:"status"` // "active", "completed", "paused"
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (b *BudgetGoal) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}

// AutoSaveRule represents an auto-save rule (round-up or auto-transfer)
type AutoSaveRule struct {
	ID        string          `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID    string          `gorm:"type:varchar(36);index;not null" json:"user_id"`
	GoalID    string          `gorm:"type:varchar(36)" json:"goal_id"`
	Type      string          `gorm:"type:varchar(20);not null" json:"type"` // "roundup", "auto_transfer"
	Amount    decimal.Decimal `gorm:"type:decimal(12,2)" json:"amount"`      // for auto_transfer
	Frequency string          `gorm:"type:varchar(20)" json:"frequency"`     // "daily", "weekly", "monthly"
	IsActive  bool            `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (a *AutoSaveRule) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

// CreateGoalRequest is the request body for creating a budget goal
type CreateGoalRequest struct {
	Name         string  `json:"name" binding:"required"`
	Category     string  `json:"category" binding:"required"`
	TargetAmount float64 `json:"target_amount" binding:"required,gt=0"`
	Deadline     string  `json:"deadline,omitempty"`
	Icon         string  `json:"icon"`
	Color        string  `json:"color"`
}

// CreateAutoSaveRequest is the request body for creating an auto-save rule
type CreateAutoSaveRequest struct {
	GoalID    string  `json:"goal_id"`
	Type      string  `json:"type" binding:"required,oneof=roundup auto_transfer"`
	Amount    float64 `json:"amount"`
	Frequency string  `json:"frequency"`
}
