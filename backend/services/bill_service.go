package services

import (
	"errors"

	"gatorpay-backend/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// BillService handles bill payment operations
type BillService struct {
	db            *gorm.DB
	rewardService *RewardService
}

// NewBillService creates a new BillService
func NewBillService(db *gorm.DB, rewardService *RewardService) *BillService {
	return &BillService{db: db, rewardService: rewardService}
}

// BillPayInput is the DTO for paying a bill
type BillPayInput struct {
	BillerID      string  `json:"biller_id" binding:"required"`
	AccountNumber string  `json:"account_number" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	SaveBiller    bool    `json:"save_biller"`
}

// BillPayResponse is the response after a successful bill payment
type BillPayResponse struct {
	PaymentID  string          `json:"payment_id"`
	Biller     models.Biller   `json:"biller"`
	Amount     decimal.Decimal `json:"amount"`
	NewBalance decimal.Decimal `json:"new_balance"`
}

// GetCategories returns distinct bill categories
func (s *BillService) GetCategories() ([]string, error) {
	var categories []string
	if err := s.db.Model(&models.Biller{}).
		Where("is_active = ?", true).
		Distinct("category").
		Pluck("category", &categories).Error; err != nil {
		return nil, errors.New("failed to fetch categories")
	}
	return categories, nil
}

// GetBillers returns billers, optionally filtered by category
func (s *BillService) GetBillers(category string) ([]models.Biller, error) {
	query := s.db.Where("is_active = ?", true)
	if category != "" {
		query = query.Where("category = ?", category)
	}

	var billers []models.Biller
	if err := query.Find(&billers).Error; err != nil {
		return nil, errors.New("failed to fetch billers")
	}
	return billers, nil
}

// PayBill processes a bill payment
func (s *BillService) PayBill(userID string, input BillPayInput) (*BillPayResponse, error) {
	amount := decimal.NewFromFloat(input.Amount)
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("amount must be greater than 0")
	}

	// Validate biller exists
	var biller models.Biller
	if err := s.db.Where("id = ? AND is_active = ?", input.BillerID, true).First(&biller).Error; err != nil {
		return nil, errors.New("biller not found")
	}

	var response BillPayResponse
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Get user wallet
		var wallet models.Wallet
		if err := tx.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			return errors.New("wallet not found")
		}

		if !wallet.IsActive {
			return errors.New("wallet is not active")
		}

		// Check balance
		if wallet.Balance.LessThan(amount) {
			return errors.New("insufficient balance")
		}

		// Debit wallet
		wallet.Balance = wallet.Balance.Sub(amount)
		if err := tx.Save(&wallet).Error; err != nil {
			return errors.New("failed to debit wallet")
		}

		// Create bill_pay transaction
		transaction := models.Transaction{
			WalletID:    wallet.ID,
			FromUserID:  &userID,
			Type:        models.TransactionTypeBillPay,
			Amount:      amount,
			Description: "Bill payment to " + biller.Name + " (Acct: " + input.AccountNumber + ")",
			Status:      models.TransactionStatusSuccess,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return errors.New("failed to create transaction")
		}

		// Create BillPayment record
		billPayment := models.BillPayment{
			UserID:        userID,
			BillerID:      biller.ID,
			AccountNumber: input.AccountNumber,
			Amount:        amount,
			Status:        "success",
		}
		if err := tx.Create(&billPayment).Error; err != nil {
			return errors.New("failed to create bill payment record")
		}

		// Save biller if requested
		if input.SaveBiller {
			var existing models.SavedBiller
			err := tx.Where("user_id = ? AND biller_id = ? AND account_number = ?",
				userID, biller.ID, input.AccountNumber).First(&existing).Error
			if err != nil {
				// Not found, create new saved biller
				savedBiller := models.SavedBiller{
					UserID:        userID,
					BillerID:      biller.ID,
					AccountNumber: input.AccountNumber,
					Nickname:      biller.Name,
				}
				tx.Create(&savedBiller)
			}
		}

		response = BillPayResponse{
			PaymentID:  billPayment.ID,
			Biller:     biller,
			Amount:     amount,
			NewBalance: wallet.Balance,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Award 2% cashback asynchronously
	go s.rewardService.AwardCashback(userID, amount, 0.02, "Cashback for bill payment to "+biller.Name)

	return &response, nil
}

// GetSavedBillers returns saved billers for a user
func (s *BillService) GetSavedBillers(userID string) ([]models.SavedBiller, error) {
	var savedBillers []models.SavedBiller
	if err := s.db.Where("user_id = ?", userID).
		Preload("Biller").
		Find(&savedBillers).Error; err != nil {
		return nil, errors.New("failed to fetch saved billers")
	}
	return savedBillers, nil
}

// RemoveSavedBiller removes a saved biller by ID for a user
func (s *BillService) RemoveSavedBiller(userID, savedBillerID string) error {
	result := s.db.Where("id = ? AND user_id = ?", savedBillerID, userID).
		Delete(&models.SavedBiller{})
	if result.Error != nil {
		return errors.New("failed to remove saved biller")
	}
	if result.RowsAffected == 0 {
		return errors.New("saved biller not found")
	}
	return nil
}
