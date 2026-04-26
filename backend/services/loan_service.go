package services

import (
	"errors"
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

func (s *LoanService) ApplyForLoan(userID string, offerID string, amount decimal.Decimal, termMonths int) (*models.Loan, error) {
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

	var wallet models.Wallet
	if err := s.db.Where("user_id = ?", userID).First(&wallet).Error; err == nil {
		wallet.Balance = wallet.Balance.Add(amount)
		s.db.Save(&wallet)
	}

	return &loan, nil
}

func (s *LoanService) GetUserLoans(userID string) ([]models.Loan, error) {
	var loans []models.Loan
	err := s.db.Where("user_id = ?", userID).Find(&loans).Error
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
