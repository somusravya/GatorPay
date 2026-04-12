package services

import (
	"errors"
	"log"

	"gatorpay-backend/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// RewardService handles rewards and cashback logic
type RewardService struct {
	db *gorm.DB
}

// NewRewardService creates a new RewardService
func NewRewardService(db *gorm.DB) *RewardService {
	return &RewardService{db: db}
}

// RewardSummary represents aggregated reward stats
type RewardSummary struct {
	TotalPoints       int             `json:"total_points"`
	TotalCashback     decimal.Decimal `json:"total_cashback"`
	LifetimeEarnings  decimal.Decimal `json:"lifetime_earnings"`
	TotalTransactions int64           `json:"total_transactions"`
}

// Offer represents a promotional offer
type Offer struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Discount    string `json:"discount"`
	Icon        string `json:"icon"`
	IsActive    bool   `json:"is_active"`
}

// AwardCashback calculates and awards cashback to a user
func (s *RewardService) AwardCashback(userID string, amount decimal.Decimal, rate float64, description string) {
	cashbackAmount := amount.Mul(decimal.NewFromFloat(rate))
	points := int(amount.IntPart()) * 10 // 10 points per dollar

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Create reward record
		reward := models.Reward{
			UserID:      userID,
			Type:        models.RewardTypeCashback,
			Amount:      cashbackAmount,
			Points:      points,
			Description: description,
		}
		if err := tx.Create(&reward).Error; err != nil {
			return errors.New("failed to create reward record")
		}

		// Credit wallet
		var wallet models.Wallet
		if err := tx.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			return errors.New("wallet not found")
		}

		wallet.Balance = wallet.Balance.Add(cashbackAmount)
		if err := tx.Save(&wallet).Error; err != nil {
			return errors.New("failed to credit wallet")
		}

		// Create cashback transaction
		transaction := models.Transaction{
			WalletID:    wallet.ID,
			ToUserID:    &userID,
			Type:        models.TransactionTypeCashback,
			Amount:      cashbackAmount,
			Description: description,
			Status:      models.TransactionStatusSuccess,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return errors.New("failed to create cashback transaction")
		}

		// Link transaction to reward
		tx.Model(&reward).Update("transaction_id", transaction.ID)

		return nil
	})

	if err != nil {
		log.Printf("⚠️ Failed to award cashback to user %s: %v", userID, err)
	} else {
		log.Printf("💰 Awarded %s cashback + %d points to user %s", cashbackAmount.String(), points, userID)
	}
}

// GetSummary returns aggregated reward stats for a user
func (s *RewardService) GetSummary(userID string) (*RewardSummary, error) {
	var summary RewardSummary

	// Total points
	s.db.Model(&models.Reward{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(points), 0)").
		Scan(&summary.TotalPoints)

	// Total cashback
	s.db.Model(&models.Reward{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&summary.TotalCashback)

	// Lifetime earnings (same as total cashback for now)
	summary.LifetimeEarnings = summary.TotalCashback

	// Total reward transactions
	s.db.Model(&models.Reward{}).
		Where("user_id = ?", userID).
		Count(&summary.TotalTransactions)

	return &summary, nil
}

// GetHistory returns paginated reward history for a user
func (s *RewardService) GetHistory(userID string, page, limit int) ([]models.Reward, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	var total int64
	s.db.Model(&models.Reward{}).Where("user_id = ?", userID).Count(&total)

	var rewards []models.Reward
	offset := (page - 1) * limit
	if err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&rewards).Error; err != nil {
		return nil, 0, errors.New("failed to fetch reward history")
	}

	return rewards, total, nil
}

// GetOffers returns a list of promotional offers
func (s *RewardService) GetOffers() []Offer {
	return []Offer{
		{
			ID:          "1",
			Title:       "Send & Save",
			Description: "Earn 1% cashback on every P2P transfer",
			Discount:    "1% Cashback",
			Icon:        "💸",
			IsActive:    true,
		},
		{
			ID:          "2",
			Title:       "Bill Pay Bonus",
			Description: "Get 2% cashback when you pay bills through GatorPay",
			Discount:    "2% Cashback",
			Icon:        "📄",
			IsActive:    true,
		},
		{
			ID:          "3",
			Title:       "Points Multiplier",
			Description: "Earn 10 GatorPoints for every dollar spent",
			Discount:    "10x Points",
			Icon:        "⭐",
			IsActive:    true,
		},
		{
			ID:          "4",
			Title:       "Refer a Friend",
			Description: "Get $5 bonus when your friend signs up and makes their first transfer",
			Discount:    "$5 Bonus",
			Icon:        "🎁",
			IsActive:    true,
		},
	}
}
