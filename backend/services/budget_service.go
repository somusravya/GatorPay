package services

import (
	"time"

	"gatorpay-backend/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// BudgetService handles budgeting goals and auto-save rules
type BudgetService struct {
	db *gorm.DB
}

// NewBudgetService creates a new BudgetService
func NewBudgetService(db *gorm.DB) *BudgetService {
	return &BudgetService{db: db}
}

// CreateGoal creates a new savings goal
func (s *BudgetService) CreateGoal(userID string, req models.CreateGoalRequest) (*models.BudgetGoal, error) {
	goal := models.BudgetGoal{
		UserID:       userID,
		Name:         req.Name,
		Category:     req.Category,
		TargetAmount: decimal.NewFromFloat(req.TargetAmount),
		Icon:         req.Icon,
		Color:        req.Color,
		Status:       "active",
	}

	if req.Icon == "" {
		icons := map[string]string{
			"emergency":  "🛡️",
			"vacation":   "✈️",
			"education":  "📚",
			"car":        "🚗",
			"home":       "🏠",
			"retirement": "🏖️",
			"wedding":    "💍",
		}
		if icon, ok := icons[req.Category]; ok {
			goal.Icon = icon
		} else {
			goal.Icon = "🎯"
		}
	}

	if req.Color == "" {
		goal.Color = "#3b82f6"
	}

	if req.Deadline != "" {
		t, err := time.Parse("2006-01-02", req.Deadline)
		if err == nil {
			goal.Deadline = &t
		}
	}

	if err := s.db.Create(&goal).Error; err != nil {
		return nil, err
	}
	return &goal, nil
}

// GetGoals returns all budget goals for a user
func (s *BudgetService) GetGoals(userID string) ([]models.BudgetGoal, error) {
	var goals []models.BudgetGoal
	err := s.db.Where("user_id = ?", userID).Order("created_at desc").Find(&goals).Error
	return goals, err
}

// CreateAutoSaveRule creates a new auto-save rule
func (s *BudgetService) CreateAutoSaveRule(userID string, req models.CreateAutoSaveRequest) (*models.AutoSaveRule, error) {
	rule := models.AutoSaveRule{
		UserID:    userID,
		GoalID:    req.GoalID,
		Type:      req.Type,
		Amount:    decimal.NewFromFloat(req.Amount),
		Frequency: req.Frequency,
		IsActive:  true,
	}

	if err := s.db.Create(&rule).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

// ExecuteRoundup simulates a round-up savings execution
func (s *BudgetService) ExecuteRoundup(userID string) (map[string]interface{}, error) {
	// Get active round-up rules
	var rules []models.AutoSaveRule
	s.db.Where("user_id = ? AND type = ? AND is_active = ?", userID, "roundup", true).Find(&rules)

	totalSaved := decimal.NewFromInt(0)
	transactionsRounded := 0

	// Get recent transactions and apply round-up
	var transactions []models.Transaction
	oneWeekAgo := time.Now().AddDate(0, 0, -7)
	s.db.Where("from_user_id = ? AND created_at >= ?", userID, oneWeekAgo).Find(&transactions)

	for range transactions {
		// Simulate round-up: each transaction rounds up to nearest dollar
		roundup := decimal.NewFromFloat(0.50) // Average round-up
		totalSaved = totalSaved.Add(roundup)
		transactionsRounded++
	}

	// If no transactions, simulate
	if transactionsRounded == 0 {
		totalSaved = decimal.NewFromFloat(3.75)
		transactionsRounded = 5
	}

	// Update goal progress if rule is associated with a goal
	for _, rule := range rules {
		if rule.GoalID != "" {
			s.db.Model(&models.BudgetGoal{}).Where("id = ?", rule.GoalID).
				Update("current_amount", gorm.Expr("current_amount + ?", totalSaved))
		}
	}

	return map[string]interface{}{
		"total_saved":          totalSaved.StringFixed(2),
		"transactions_rounded": transactionsRounded,
		"rules_applied":        len(rules),
	}, nil
}
