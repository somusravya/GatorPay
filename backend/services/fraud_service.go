package services

import (
	"encoding/json"
	"math"
	"time"

	"gatorpay-backend/models"

	"gorm.io/gorm"
)

// FraudService handles fraud detection and risk scoring
type FraudService struct {
	db *gorm.DB
}

// NewFraudService creates a new FraudService
func NewFraudService(db *gorm.DB) *FraudService {
	return &FraudService{db: db}
}

// GetAlerts returns fraud alerts, optionally filtered by status
func (s *FraudService) GetAlerts(status string) ([]models.FraudAlert, error) {
	var alerts []models.FraudAlert
	query := s.db.Preload("User").Order("created_at desc")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Limit(50).Find(&alerts).Error
	return alerts, err
}

// ReviewAlert updates the status of a fraud alert
func (s *FraudService) ReviewAlert(alertID string, reviewerID string, action string) error {
	now := time.Now()
	return s.db.Model(&models.FraudAlert{}).Where("id = ?", alertID).Updates(map[string]interface{}{
		"status":      action,
		"reviewed_by": reviewerID,
		"reviewed_at": &now,
	}).Error
}

// AssessTransactionRisk evaluates the risk of a transaction
func (s *FraudService) AssessTransactionRisk(userID string, amount float64, txType string) (float64, string) {
	riskScore := 0.0
	factors := make(map[string]float64)

	// Factor 1: Amount-based risk
	if amount > 5000 {
		factors["high_amount"] = 30
		riskScore += 30
	} else if amount > 1000 {
		factors["medium_amount"] = 15
		riskScore += 15
	}

	// Factor 2: Velocity check - count recent transactions
	var recentCount int64
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	s.db.Model(&models.Transaction{}).
		Where("from_user_id = ? AND created_at >= ?", userID, oneHourAgo).
		Count(&recentCount)

	if recentCount > 10 {
		factors["high_velocity"] = 35
		riskScore += 35
	} else if recentCount > 5 {
		factors["medium_velocity"] = 15
		riskScore += 15
	}

	// Factor 3: Time-based risk (late night transactions)
	hour := time.Now().Hour()
	if hour >= 0 && hour < 6 {
		factors["odd_hours"] = 10
		riskScore += 10
	}

	// Factor 4: Account age risk
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err == nil {
		accountAge := time.Since(user.CreatedAt).Hours() / 24
		if accountAge < 7 {
			factors["new_account"] = 20
			riskScore += 20
		}
	}

	// Cap at 100
	riskScore = math.Min(riskScore, 100)

	// Determine action
	action := "allow"
	if riskScore >= 70 {
		action = "block"
	} else if riskScore >= 40 {
		action = "flag"
	}

	// Log risk event
	factorsJSON, _ := json.Marshal(factors)
	riskEvent := models.RiskEvent{
		UserID:    userID,
		EventType: txType,
		RiskScore: riskScore,
		Factors:   string(factorsJSON),
		Action:    action,
	}
	s.db.Create(&riskEvent)

	// Create fraud alert if flagged
	if action != "allow" {
		alert := models.FraudAlert{
			UserID:      userID,
			RiskScore:   riskScore,
			Type:        determineAlertType(factors),
			Status:      "pending",
			Description: "Suspicious activity detected",
			Details:     string(factorsJSON),
		}
		s.db.Create(&alert)
	}

	return riskScore, action
}

func determineAlertType(factors map[string]float64) string {
	if _, ok := factors["high_velocity"]; ok {
		return "velocity"
	}
	if _, ok := factors["high_amount"]; ok {
		return "amount_spike"
	}
	if _, ok := factors["odd_hours"]; ok {
		return "geo_anomaly"
	}
	return "suspicious_merchant"
}
