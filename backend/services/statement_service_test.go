package services_test

import (
	"bytes"
	"gatorpay-backend/models"
	"gatorpay-backend/services"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestStatements(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&models.Transaction{})
	ss := services.NewStatementService(db)

	u1 := "u1"
	db.Create(&models.Transaction{FromUserID: &u1, Amount: decimal.NewFromInt(100), Type: "CREDIT", Description: "Paycheck"})

	csvData, _, err := ss.GenerateStatement("u1", "csv", "", "")
	assert.NoError(t, err)
	assert.True(t, bytes.Contains(csvData, []byte("Paycheck")))

	pdfData, _, err := ss.GenerateStatement("u1", "pdf", "", "")
	assert.NoError(t, err)
	assert.True(t, bytes.Contains(pdfData, []byte("%PDF")))
}
