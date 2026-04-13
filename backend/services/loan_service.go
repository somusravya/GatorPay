package services

import (
	"errors"

	"gatorpay-backend/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// LoanService handles loan offers and user loans.
type LoanService struct {
	db *gorm.DB
}

// NewLoanService creates a LoanService.
func NewLoanService(db *gorm.DB) *LoanService {
	return &LoanService{db: db}
}

// GetOffers returns active loan products.
func (s *LoanService) GetOffers() ([]models.LoanOffer, error) {
	var offers []models.LoanOffer
	if err := s.db.Where("is_active = ?", true).Find(&offers).Error; err != nil {
		return nil, err
	}
	return offers, nil
}

// ApplyForLoan creates a loan after validating against the offer limits.
func (s *LoanService) ApplyForLoan(userID, offerID string, amount decimal.Decimal, termMonths int) (*models.Loan, error) {
	var offer models.LoanOffer
	if err := s.db.Where("id = ? AND is_active = ?", offerID, true).First(&offer).Error; err != nil {
		return nil, errors.New("invalid offer")
	}
	if amount.LessThan(offer.MinAmount) || amount.GreaterThan(offer.MaxAmount) {
		return nil, errors.New("amount outside allowed range")
	}
	if termMonths != offer.TermMonths {
		return nil, errors.New("term does not match offer")
	}

	monthlyRate := offer.InterestRate.Div(decimal.NewFromInt(1200))
	n := decimal.NewFromInt(int64(termMonths))
	one := decimal.NewFromInt(1)
	pow := one.Add(monthlyRate).Pow(n)
	em := amount.Mul(monthlyRate).Mul(pow).Div(pow.Sub(one)).Round(2)

	loan := models.Loan{
		UserID:             userID,
		OfferID:            offerID,
		Principal:          amount,
		TermMonths:         termMonths,
		InterestRate:       offer.InterestRate,
		OutstandingBalance: amount,
		MonthlyPayment:     em,
		Status:             "active",
	}
	if err := s.db.Create(&loan).Error; err != nil {
		return nil, err
	}
	return &loan, nil
}

// GetUserLoans lists loans for a user.
func (s *LoanService) GetUserLoans(userID string) ([]models.Loan, error) {
	var loans []models.Loan
	if err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&loans).Error; err != nil {
		return nil, err
	}
	return loans, nil
}

// GetLoan returns one loan if owned by user.
func (s *LoanService) GetLoan(loanID, userID string) (*models.Loan, error) {
	var loan models.Loan
	if err := s.db.Where("id = ? AND user_id = ?", loanID, userID).First(&loan).Error; err != nil {
		return nil, err
	}
	return &loan, nil
}

// MakeLoanPayment applies one EMI toward outstanding balance.
func (s *LoanService) MakeLoanPayment(loanID, userID string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var loan models.Loan
		if err := tx.Where("id = ? AND user_id = ?", loanID, userID).First(&loan).Error; err != nil {
			return errors.New("loan not found")
		}
		if loan.Status != "active" {
			return errors.New("loan not active")
		}
		payment := loan.MonthlyPayment
		if loan.OutstandingBalance.LessThanOrEqual(decimal.Zero) {
			loan.Status = "closed"
			return tx.Save(&loan).Error
		}
		newBal := loan.OutstandingBalance.Sub(payment)
		if newBal.LessThan(decimal.Zero) {
			newBal = decimal.Zero
		}
		loan.OutstandingBalance = newBal
		if newBal.Equal(decimal.Zero) {
			loan.Status = "closed"
		}
		return tx.Save(&loan).Error
	})
}
