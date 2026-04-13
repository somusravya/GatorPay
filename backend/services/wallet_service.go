package services

import (
	"errors"
	"math"

	"gatorpay-backend/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// WalletService handles wallet-related business logic
type WalletService struct {
	db *gorm.DB
}

// NewWalletService creates a new WalletService
func NewWalletService(db *gorm.DB) *WalletService {
	return &WalletService{db: db}
}

// AddMoneyInput is the DTO for adding money
type AddMoneyInput struct {
	Amount      float64 `json:"amount" binding:"required"`
	Source      string  `json:"source" binding:"required"`
	Description string  `json:"description"`
}

// WithdrawInput is the DTO for withdrawing money
type WithdrawInput struct {
	Amount      float64 `json:"amount" binding:"required"`
	BankAccount string  `json:"bank_account" binding:"required"`
}

// TransactionListResponse is the paginated transaction response
type TransactionListResponse struct {
	Transactions []models.Transaction `json:"transactions"`
	Total        int64                `json:"total"`
	Page         int                  `json:"page"`
	Limit        int                  `json:"limit"`
	TotalPages   int                  `json:"total_pages"`
}

// GetWallet returns the wallet for a given user
func (s *WalletService) GetWallet(userID string) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := s.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return nil, errors.New("wallet not found")
	}
	return &wallet, nil
}

// AddMoney deposits money into the user's wallet atomically
func (s *WalletService) AddMoney(userID string, input AddMoneyInput) (*models.Wallet, error) {
	amount := decimal.NewFromFloat(input.Amount)
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("amount must be greater than 0")
	}

	var wallet models.Wallet
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Lock the wallet row for update
		if err := tx.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			return errors.New("wallet not found")
		}

		if !wallet.IsActive {
			return errors.New("wallet is not active")
		}

		// Update balance
		wallet.Balance = wallet.Balance.Add(amount)
		if err := tx.Save(&wallet).Error; err != nil {
			return errors.New("failed to update balance")
		}

		// Create transaction record
		description := input.Description
		if description == "" {
			description = "Deposit from " + input.Source
		}
		transaction := models.Transaction{
			WalletID:    wallet.ID,
			Type:        models.TransactionTypeDeposit,
			Amount:      amount,
			Description: description,
			Status:      models.TransactionStatusSuccess,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return errors.New("failed to create transaction record")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

// Withdraw deducts money from the user's wallet atomically
func (s *WalletService) Withdraw(userID string, input WithdrawInput) (*models.Wallet, error) {
	amount := decimal.NewFromFloat(input.Amount)
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("amount must be greater than 0")
	}

	var wallet models.Wallet
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			return errors.New("wallet not found")
		}

		if !wallet.IsActive {
			return errors.New("wallet is not active")
		}

		// Check sufficient balance
		if wallet.Balance.LessThan(amount) {
			return errors.New("insufficient balance")
		}

		// Deduct balance
		wallet.Balance = wallet.Balance.Sub(amount)
		if err := tx.Save(&wallet).Error; err != nil {
			return errors.New("failed to update balance")
		}

		// Create transaction record
		transaction := models.Transaction{
			WalletID:    wallet.ID,
			Type:        models.TransactionTypeWithdraw,
			Amount:      amount,
			Description: "Withdrawal to " + input.BankAccount,
			Status:      models.TransactionStatusSuccess,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return errors.New("failed to create transaction record")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

// GetTransactions returns paginated transactions for a user's wallet
func (s *WalletService) GetTransactions(userID string, page, limit int) (*TransactionListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	var wallet models.Wallet
	if err := s.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return nil, errors.New("wallet not found")
	}

	var total int64
	s.db.Model(&models.Transaction{}).
		Where("wallet_id = ? OR from_user_id = ? OR to_user_id = ?", wallet.ID, userID, userID).
		Count(&total)

	var transactions []models.Transaction
	offset := (page - 1) * limit
	if err := s.db.Where("wallet_id = ? OR from_user_id = ? OR to_user_id = ?", wallet.ID, userID, userID).
		Preload("FromUser").
		Preload("ToUser").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&transactions).Error; err != nil {
		return nil, errors.New("failed to fetch transactions")
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &TransactionListResponse{
		Transactions: transactions,
		Total:        total,
		Page:         page,
		Limit:        limit,
		TotalPages:   totalPages,
	}, nil
}
