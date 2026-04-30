package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InsightReport stores generated AI insight reports for a user
type InsightReport struct {
	ID                  string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID              string         `gorm:"type:varchar(36);index;not null" json:"user_id"`
	HealthScore         int            `gorm:"default:50" json:"health_score"` // 0-100
	TotalIncome         float64        `json:"total_income"`
	TotalSpending       float64        `json:"total_spending"`
	SavingsRate         float64        `json:"savings_rate"`
	TopCategory         string         `json:"top_category"`
	AnomalyCount        int            `json:"anomaly_count"`
	RecommendationCount int            `json:"recommendation_count"`
	Period              string         `gorm:"type:varchar(20)" json:"period"` // "monthly", "weekly"
	GeneratedAt         time.Time      `json:"generated_at"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}

func (i *InsightReport) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	return nil
}

// SpendingCategory represents a spending category breakdown
type SpendingCategory struct {
	Category   string  `json:"category"`
	Amount     float64 `json:"amount"`
	Percentage float64 `json:"percentage"`
	Trend      string  `json:"trend"` // "up", "down", "stable"
	Icon       string  `json:"icon"`
}

// SpendingAnomaly represents a detected spending anomaly
type SpendingAnomaly struct {
	Category    string  `json:"category"`
	Amount      float64 `json:"amount"`
	Average     float64 `json:"average"`
	Deviation   float64 `json:"deviation"`
	Description string  `json:"description"`
	Severity    string  `json:"severity"` // "low", "medium", "high"
}

// MonthlyForecast represents a monthly spending forecast
type MonthlyForecast struct {
	Month     string  `json:"month"`
	Predicted float64 `json:"predicted"`
	Actual    float64 `json:"actual"`
	Trend     string  `json:"trend"`
}

// InsightRecommendation represents a smart savings recommendation
type InsightRecommendation struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Impact      string  `json:"impact"`
	Savings     float64 `json:"potential_savings"`
	Category    string  `json:"category"`
	Priority    string  `json:"priority"` // "high", "medium", "low"
	Icon        string  `json:"icon"`
}

// InsightSummaryResponse is the full response for GET /insights/summary
type InsightSummaryResponse struct {
	HealthScore     int                     `json:"health_score"`
	HealthGrade     string                  `json:"health_grade"`
	TotalIncome     float64                 `json:"total_income"`
	TotalSpending   float64                 `json:"total_spending"`
	SavingsRate     float64                 `json:"savings_rate"`
	Categories      []SpendingCategory      `json:"categories"`
	Anomalies       []SpendingAnomaly       `json:"anomalies"`
	Forecasts       []MonthlyForecast       `json:"forecasts"`
	Recommendations []InsightRecommendation `json:"recommendations"`
	Period          string                  `json:"period"`
	GeneratedAt     time.Time               `json:"generated_at"`
}
