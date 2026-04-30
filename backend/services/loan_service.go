package services

import (
	"errors"
	"fmt"
	"math"
	"time"

	"gatorpay-backend/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type LoanService struct {
	db *gorm.DB
}

func NewLoanService(db *gorm.DB) *LoanService {
	return &LoanService{db: db}
}

func (s *LoanService) GetOffers() ([]models.LoanOffer, error) {
	var offers []models.LoanOffer
	err := s.db.Find(&offers).Error
	return offers, err
}

func calculateEMI(principal decimal.Decimal, rate decimal.Decimal, months int) decimal.Decimal {
	rateFloat, _ := rate.Float64()
	if rateFloat == 0 {
		return principal.Div(decimal.NewFromInt(int64(months)))
	}
	r := (rateFloat / 100.0) / 12.0
	n := float64(months)

	pFloat, _ := principal.Float64()
	emiFloat := pFloat * r * math.Pow(1+r, n) / (math.Pow(1+r, n) - 1)
	return decimal.NewFromFloat(emiFloat).RoundBank(2)
}

// CheckEligibility evaluates a user's eligibility for loans based on credit score, balance, and existing loans
func (s *LoanService) CheckEligibility(userID string) (map[string]interface{}, error) {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	var wallet models.Wallet
	if err := s.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return nil, errors.New("wallet not found")
	}

	// Count active loans
	var activeLoans int64
	s.db.Model(&models.Loan{}).Where("user_id = ? AND status = ?", userID, "active").Count(&activeLoans)

	// Calculate total outstanding
	var totalOutstanding decimal.Decimal
	var loans []models.Loan
	s.db.Where("user_id = ? AND status = ?", userID, "active").Find(&loans)
	for _, l := range loans {
		totalOutstanding = totalOutstanding.Add(l.RemainingAmount)
	}

	eligible := true
	reasons := []string{}

	// Rule 1: Credit score must be >= 500
	if user.CreditScore < 500 {
		eligible = false
		reasons = append(reasons, fmt.Sprintf("Credit score too low (%d). Minimum required: 500.", user.CreditScore))
	}

	// Rule 2: Must have wallet balance >= $50
	minBalance := decimal.NewFromInt(50)
	if wallet.Balance.LessThan(minBalance) {
		eligible = false
		balStr, _ := wallet.Balance.Float64()
		reasons = append(reasons, fmt.Sprintf("Wallet balance too low ($%.2f). Minimum required: $50.00.", balStr))
	}

	// Rule 3: Max 3 active loans at a time
	if activeLoans >= 3 {
		eligible = false
		reasons = append(reasons, fmt.Sprintf("Too many active loans (%d). Maximum allowed: 3.", activeLoans))
	}

	// Rule 4: Total outstanding must be < $20,000
	maxOutstanding := decimal.NewFromInt(20000)
	if totalOutstanding.GreaterThanOrEqual(maxOutstanding) {
		eligible = false
		outStr, _ := totalOutstanding.Float64()
		reasons = append(reasons, fmt.Sprintf("Total outstanding debt too high ($%.2f). Maximum: $20,000.", outStr))
	}

	// Determine max eligible amount based on credit score
	maxAmount := decimal.NewFromInt(1000)
	if user.CreditScore >= 750 {
		maxAmount = decimal.NewFromInt(10000)
	} else if user.CreditScore >= 650 {
		maxAmount = decimal.NewFromInt(5000)
	} else if user.CreditScore >= 500 {
		maxAmount = decimal.NewFromInt(2000)
	}

	if eligible {
		reasons = append(reasons, "You meet all eligibility requirements.")
	}

	return map[string]interface{}{
		"eligible":          eligible,
		"credit_score":      user.CreditScore,
		"wallet_balance":    wallet.Balance,
		"active_loans":      activeLoans,
		"total_outstanding": totalOutstanding,
		"max_loan_amount":   maxAmount,
		"reasons":           reasons,
	}, nil
}

func (s *LoanService) ApplyForLoan(userID string, offerID string, amount decimal.Decimal, termMonths int) (*models.Loan, error) {
	// Check eligibility first
	eligibility, err := s.CheckEligibility(userID)
	if err != nil {
		return nil, err
	}
	if !eligibility["eligible"].(bool) {
		return nil, errors.New("you are not eligible for a loan at this time")
	}

	var offer models.LoanOffer
	if err := s.db.First(&offer, "id = ?", offerID).Error; err != nil {
		return nil, errors.New("invalid loan offer")
	}

	if amount.LessThan(offer.MinAmount) || amount.GreaterThan(offer.MaxAmount) {
		return nil, errors.New("amount outside allowed range")
	}

	if termMonths <= 0 || termMonths > offer.TermMonths {
		return nil, errors.New("invalid term")
	}

	emi := calculateEMI(amount, offer.InterestRate, termMonths)
	totalPayable := emi.Mul(decimal.NewFromInt(int64(termMonths)))

	loan := models.Loan{
		UserID:          userID,
		Amount:          amount,
		TermMonths:      termMonths,
		InterestRate:    offer.InterestRate,
		EMI:             emi,
		TotalPayable:    totalPayable,
		AmountPaid:      decimal.NewFromInt(0),
		RemainingAmount: totalPayable,
		Status:          "active",
		NextPaymentDate: time.Now().AddDate(0, 1, 0),
	}

	if err := s.db.Create(&loan).Error; err != nil {
		return nil, err
	}

	// Disburse: credit loan amount to wallet and create transaction log
	var wallet models.Wallet
	if err := s.db.Where("user_id = ?", userID).First(&wallet).Error; err == nil {
		wallet.Balance = wallet.Balance.Add(amount)
		s.db.Save(&wallet)

		// Transaction log for disbursement
		txn := models.Transaction{
			WalletID:    wallet.ID,
			Type:        "loan_disbursement",
			Amount:      amount,
			Description: fmt.Sprintf("Loan disbursement: %s (ID: %s)", offer.Name, loan.ID[:8]),
			Status:      models.TransactionStatusSuccess,
		}
		s.db.Create(&txn)
	}

	return &loan, nil
}

// CancelLoan cancels a loan if no payments have been made yet
func (s *LoanService) CancelLoan(loanID string, userID string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var loan models.Loan
		if err := tx.Where("id = ? AND user_id = ?", loanID, userID).First(&loan).Error; err != nil {
			return errors.New("loan not found")
		}

		if loan.Status == "cancelled" {
			return errors.New("loan is already cancelled")
		}

		if loan.Status == "paid" {
			return errors.New("cannot cancel a fully paid loan")
		}

		// Only allow cancel if zero payments made
		if loan.AmountPaid.GreaterThan(decimal.NewFromInt(0)) {
			return errors.New("cannot cancel a loan with payments already made")
		}

		// Reverse the disbursement: deduct loan amount from wallet
		var wallet models.Wallet
		if err := tx.Where("user_id = ?", userID).First(&wallet).Error; err == nil {
			if wallet.Balance.LessThan(loan.Amount) {
				return errors.New("insufficient wallet balance to reverse disbursement")
			}
			wallet.Balance = wallet.Balance.Sub(loan.Amount)
			tx.Save(&wallet)

			// Transaction log for reversal
			txn := models.Transaction{
				WalletID:    wallet.ID,
				Type:        "loan_reversal",
				Amount:      loan.Amount,
				Description: fmt.Sprintf("Loan cancelled and reversed (ID: %s)", loan.ID[:8]),
				Status:      models.TransactionStatusSuccess,
			}
			tx.Create(&txn)
		}

		loan.Status = "cancelled"
		return tx.Save(&loan).Error
	})
}

func (s *LoanService) GetUserLoans(userID string) ([]models.Loan, error) {
	var loans []models.Loan
	err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&loans).Error
	return loans, err
}

func (s *LoanService) GetLoan(loanID string, userID string) (*models.Loan, error) {
	var loan models.Loan
	err := s.db.Where("id = ? AND user_id = ?", loanID, userID).First(&loan).Error
	if err != nil {
		return nil, err
	}
	return &loan, nil
}

func (s *LoanService) MakeLoanPayment(loanID string, userID string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var loan models.Loan
		if err := tx.Where("id = ? AND user_id = ?", loanID, userID).First(&loan).Error; err != nil {
			return errors.New("loan not found")
		}

		if loan.Status != "active" {
			return errors.New("loan is not active")
		}

		paymentAmount := loan.EMI
		if loan.RemainingAmount.LessThan(loan.EMI) {
			paymentAmount = loan.RemainingAmount
		}

		var wallet models.Wallet
		if err := tx.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			return errors.New("wallet not found")
		}

		if wallet.Balance.LessThan(paymentAmount) {
			return errors.New("insufficient wallet balance")
		}

		wallet.Balance = wallet.Balance.Sub(paymentAmount)
		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		// Transaction log for EMI payment
		txn := models.Transaction{
			WalletID:    wallet.ID,
			Type:        "loan_emi_payment",
			Amount:      paymentAmount,
			Description: fmt.Sprintf("EMI payment for loan %s", loan.ID[:8]),
			Status:      models.TransactionStatusSuccess,
		}
		if err := tx.Create(&txn).Error; err != nil {
			return err
		}

		loan.AmountPaid = loan.AmountPaid.Add(paymentAmount)
		loan.RemainingAmount = loan.RemainingAmount.Sub(paymentAmount)
		loan.NextPaymentDate = loan.NextPaymentDate.AddDate(0, 1, 0)

		if loan.RemainingAmount.LessThanOrEqual(decimal.NewFromFloat(0.01)) {
			loan.Status = "paid"
			loan.RemainingAmount = decimal.NewFromInt(0)
		}

		return tx.Save(&loan).Error
	})
}
