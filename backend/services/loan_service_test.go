package services_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	"gatorpay-backend/models"
	"gatorpay-backend/services"
	"github.com/shopspring/decimal"
)

func setupTestDBL() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&models.Loan{}, &models.LoanOffer{}, &models.Wallet{})
	return db
}

func TestLoanApply(t *testing.T) {
	db := setupTestDBL()
	ls := services.NewLoanService(db)

	offer := models.LoanOffer{ID: "offer1", MinAmount: decimal.NewFromInt(100), MaxAmount: decimal.NewFromInt(5000), TermMonths: 12, InterestRate: decimal.NewFromInt(5)}
	db.Create(&offer)
	db.Create(&models.Wallet{UserID: "u1", Balance: decimal.NewFromInt(0)})

	loan, err := ls.ApplyForLoan("u1", "offer1", decimal.NewFromInt(1000), 12)
	assert.NoError(t, err)
	assert.NotNil(t, loan)
	
	var w models.Wallet
	db.First(&w, "user_id = ?", "u1")
	assert.True(t, w.Balance.Equal(decimal.NewFromInt(1000)))
}

func TestLoanApplyFailValues(t *testing.T) {
	db := setupTestDBL()
	ls := services.NewLoanService(db)

	offer := models.LoanOffer{ID: "offer1", MinAmount: decimal.NewFromInt(100), MaxAmount: decimal.NewFromInt(5000), TermMonths: 12, InterestRate: decimal.NewFromInt(5)}
	db.Create(&offer)

	_, err := ls.ApplyForLoan("u1", "offer1", decimal.NewFromInt(10), 12) // Below Min
	assert.Error(t, err)

	_, err2 := ls.ApplyForLoan("u1", "offer1", decimal.NewFromInt(1000), 24) // Above Term
	assert.Error(t, err2)
}
