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

func setupTestDBQ() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&models.Merchant{}, &models.MerchantQRCode{}, &models.Wallet{})
	return db
}

func TestQRGenerations(t *testing.T) {
	db := setupTestDBQ()
	qs := services.NewQRService(db)

	_, err := qs.RegisterMerchant("merch1_qr", "Gator Store", "Retail")
	assert.NoError(t, err)

	qr, err := qs.GenerateQR("merch1_qr", decimal.NewFromInt(50), false)
	assert.NoError(t, err)
	assert.NotNil(t, qr)
	assert.Contains(t, qr.Base64PNG, "data:image/png;base64")
	
	// Test Pay
	db.Create(&models.Wallet{UserID: "merch1_qr", Balance: decimal.Zero})
	db.Create(&models.Wallet{UserID: "u1_qr", Balance: decimal.NewFromInt(100)})
	err = qs.PayViaQR("u1_qr", qr.CodeString, decimal.NewFromInt(50))
	assert.NoError(t, err)
	
	// Verify Wallet Adjustments
	var payerWallet models.Wallet
	db.First(&payerWallet, "user_id = ?", "u1_qr")
	// Paid 50, but 1.5% cashback = 0.75 -> 50.75 remaining!
	assert.True(t, payerWallet.Balance.Equal(decimal.NewFromFloat(50.75)))
}
