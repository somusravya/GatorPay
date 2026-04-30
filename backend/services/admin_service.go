package services

import (
	"time"

	"gatorpay-backend/models"

	"gorm.io/gorm"
)

// AdminService handles admin analytics and platform management
type AdminService struct {
	db *gorm.DB
}

// NewAdminService creates a new AdminService
func NewAdminService(db *gorm.DB) *AdminService {
	return &AdminService{db: db}
}

// GetMetrics returns platform-wide metrics
func (s *AdminService) GetMetrics() map[string]interface{} {
	var totalUsers int64
	s.db.Model(&models.User{}).Count(&totalUsers)

	var totalWallets int64
	s.db.Model(&models.Wallet{}).Count(&totalWallets)

	var totalTransactions int64
	s.db.Model(&models.Transaction{}).Count(&totalTransactions)

	var totalLoans int64
	s.db.Model(&models.Loan{}).Count(&totalLoans)

	var totalCards int64
	s.db.Model(&models.VirtualCard{}).Count(&totalCards)

	var pendingFraudAlerts int64
	s.db.Model(&models.FraudAlert{}).Where("status = ?", "pending").Count(&pendingFraudAlerts)

	var totalSubscriptions int64
	s.db.Model(&models.Subscription{}).Where("status IN ?", []string{"active", "detected"}).Count(&totalSubscriptions)

	var invoiceCount int64
	s.db.Model(&models.Invoice{}).Count(&invoiceCount)

	// Weekly new users
	weekAgo := time.Now().AddDate(0, 0, -7)
	var newUsersThisWeek int64
	s.db.Model(&models.User{}).Where("created_at >= ?", weekAgo).Count(&newUsersThisWeek)

	// Monthly active users (approximation)
	monthAgo := time.Now().AddDate(0, -1, 0)
	var activeUsersMonth int64
	s.db.Model(&models.Transaction{}).Where("created_at >= ?", monthAgo).
		Distinct("from_user_id").Count(&activeUsersMonth)

	return map[string]interface{}{
		"total_users":          totalUsers,
		"total_wallets":        totalWallets,
		"total_transactions":   totalTransactions,
		"total_loans":          totalLoans,
		"total_cards":          totalCards,
		"pending_fraud_alerts": pendingFraudAlerts,
		"new_users_this_week":  newUsersThisWeek,
		"active_users_month":   activeUsersMonth,
		"platform_gmv":         calculateGMV(s.db),
		"revenue_estimate":     calculateRevenue(s.db),
		"user_growth_rate":     calculateGrowthRate(totalUsers, newUsersThisWeek),
		"trading_volume":       calculateTradingVolume(s.db),
		"merchant_volume":      calculateMerchantVolume(s.db),
		"subscription_count":   totalSubscriptions,
		"invoice_count":        invoiceCount,
		"fraud_loss_rate":      calculateFraudLossRate(s.db),
		"generated_at":         time.Now(),
	}
}

// GetUsers returns paginated list of users for admin management
func (s *AdminService) GetUsers(page, limit int, search string) ([]models.UserResponse, int64, error) {
	var users []models.User
	var total int64

	query := s.db.Model(&models.User{})
	if search != "" {
		query = query.Where("LOWER(email) LIKE LOWER(?) OR LOWER(username) LIKE LOWER(?) OR LOWER(first_name) LIKE LOWER(?) OR LOWER(last_name) LIKE LOWER(?)",
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	var responses []models.UserResponse
	for _, u := range users {
		responses = append(responses, u.ToResponse())
	}
	return responses, total, nil
}

// GetFraudReviewQueue returns pending fraud alerts for admin review
func (s *AdminService) GetFraudReviewQueue() ([]models.FraudAlert, error) {
	var alerts []models.FraudAlert
	err := s.db.Preload("User").
		Where("status = ?", "pending").
		Order("risk_score desc").
		Limit(50).Find(&alerts).Error
	return alerts, err
}

func calculateGMV(db *gorm.DB) float64 {
	var count int64
	db.Model(&models.Transaction{}).Count(&count)
	// Rough estimation
	return float64(count) * 125.50
}

func calculateRevenue(db *gorm.DB) float64 {
	return calculateGMV(db) * 0.015 // 1.5% take rate
}

func calculateGrowthRate(total, newThisWeek int64) float64 {
	if total == 0 {
		return 0
	}
	return (float64(newThisWeek) / float64(total)) * 100
}

func calculateTradingVolume(db *gorm.DB) float64 {
	var count int64
	db.Model(&models.Transaction{}).Where("type IN ?", []string{"trade_buy", "trade_sell"}).Count(&count)
	return float64(count) * 325.75
}

func calculateMerchantVolume(db *gorm.DB) float64 {
	var count int64
	db.Model(&models.Invoice{}).Where("status IN ?", []string{"paid", "partial"}).Count(&count)
	return float64(count) * 89.25
}

func calculateFraudLossRate(db *gorm.DB) float64 {
	var total int64
	var confirmed int64
	db.Model(&models.FraudAlert{}).Count(&total)
	db.Model(&models.FraudAlert{}).Where("status = ?", "confirmed").Count(&confirmed)
	if total == 0 {
		return 0
	}
	return (float64(confirmed) / float64(total)) * 100
}
