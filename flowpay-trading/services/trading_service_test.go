package services_test

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	
	"flowpay-trading/config"
	"flowpay-trading/models"
	"flowpay-trading/services"
)

func setupTestDB() {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&models.TradingAccount{}, &models.Position{}, &models.Trade{})
	config.DB = db
}

func TestVerifyAccount(t *testing.T) {
	setupTestDB()
	ts := services.NewTradingService(&services.StockService{})

	err := ts.VerifyAccount("1", "1990-01-01", "123456789", "Q1", "A1", "Q2", "A2", true)
	assert.NoError(t, err)

	account, err := ts.GetAccount("1")
	assert.NoError(t, err)
	assert.Equal(t, "verified", account.Status)
	assert.Equal(t, 10000.00, account.BuyingPower)
	
	// Check Hash
	hash := sha256.Sum256([]byte("123456789"))
	expectedHash := hex.EncodeToString(hash[:])
	assert.Equal(t, expectedHash, account.SSNHash)
}

func TestExecuteTrade(t *testing.T) {
	setupTestDB()
	ts := services.NewTradingService(&services.StockService{})

	// Setup verified account
	ts.VerifyAccount("1", "1990-01-01", "123456789", "Q1", "A1", "Q2", "A2", true)

	err := ts.ExecuteTrade("1", "AAPL", 150.0, 10, "buy")
	assert.NoError(t, err)

	account, _ := ts.GetAccount("1")
	assert.Equal(t, 8500.00, account.BuyingPower) // 10000 - 1500

	port, _ := ts.GetPortfolio("1")
	positions := port["positions"].([]models.Position)
	assert.Len(t, positions, 1)
	assert.Equal(t, "AAPL", positions[0].Symbol)
	assert.Equal(t, 10, positions[0].Quantity)

	err = ts.ExecuteTrade("1", "AAPL", 200.0, 5, "sell")
	assert.NoError(t, err)

	account, _ = ts.GetAccount("1")
	assert.Equal(t, 9500.00, account.BuyingPower) // 8500 + 1000
}
