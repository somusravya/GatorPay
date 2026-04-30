package services_test

import (
	"gatorpay-backend/models"
	"gatorpay-backend/services"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func setupTestDBL() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	db.AutoMigrate(&models.User{}, &models.Loan{}, &models.LoanOffer{}, &models.Wallet{})
	return db
}

func TestLoanApply(t *testing.T) {
	db := setupTestDBL()
	ls := services.NewLoanService(db)

	offer := models.LoanOffer{ID: "offer1", MinAmount: decimal.NewFromInt(100), MaxAmount: decimal.NewFromInt(5000), TermMonths: 12, InterestRate: decimal.NewFromInt(5)}
	db.Create(&offer)
	db.Create(&models.User{ID: "u1", Email: "u1@example.com", Username: "u1", Phone: "1111111111", FirstName: "Test", LastName: "User", CreditScore: 700})
	db.Create(&models.Wallet{UserID: "u1", Balance: decimal.NewFromInt(100)})

	loan, err := ls.ApplyForLoan("u1", "offer1", decimal.NewFromInt(1000), 12)
	assert.NoError(t, err)
	assert.NotNil(t, loan)

	var w models.Wallet
	db.First(&w, "user_id = ?", "u1")
	assert.True(t, w.Balance.Equal(decimal.NewFromInt(1100)))
}

func TestLoanApplyFailValues(t *testing.T) {
	db := setupTestDBL()
	ls := services.NewLoanService(db)

	offer := models.LoanOffer{ID: "offer1", MinAmount: decimal.NewFromInt(100), MaxAmount: decimal.NewFromInt(5000), TermMonths: 12, InterestRate: decimal.NewFromInt(5)}
	db.Create(&offer)
	db.Create(&models.User{ID: "u1", Email: "u1@example.com", Username: "u1", Phone: "1111111111", FirstName: "Test", LastName: "User", CreditScore: 700})
	db.Create(&models.Wallet{UserID: "u1", Balance: decimal.NewFromInt(100)})

	_, err := ls.ApplyForLoan("u1", "offer1", decimal.NewFromInt(10), 12) // Below Min
	assert.Error(t, err)

	_, err2 := ls.ApplyForLoan("u1", "offer1", decimal.NewFromInt(1000), 24) // Above Term
	assert.Error(t, err2)
}
