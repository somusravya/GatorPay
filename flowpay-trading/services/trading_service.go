package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"flowpay-trading/config"
	"flowpay-trading/models"
)

type TradingService struct {
	stockService *StockService
}

func NewTradingService(ss *StockService) *TradingService {
	return &TradingService{
		stockService: ss,
	}
}

func hashString(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (s *TradingService) VerifyAccount(userID string, dob string, ssn string, q1, a1, q2, a2 string, riskAck bool) error {
	// Parse DOB
	dobTime, err := time.Parse("2006-01-02", dob)
	if err != nil {
		return errors.New("invalid date of birth format, use YYYY-MM-DD")
	}

	// Calculate age
	age := time.Now().Year() - dobTime.Year()
	if time.Now().YearDay() < dobTime.YearDay() {
		age--
	}
	if age < 18 {
		return errors.New("must be at least 18 years old to trade")
	}

	// SSN Rules
	if len(ssn) != 9 {
		return errors.New("SSN must be exactly 9 digits")
	}

	if !riskAck {
		return errors.New("you must acknowledge the risk of trading")
	}

	var account models.TradingAccount
	result := config.DB.Where("user_id = ?", userID).First(&account)

	hashSSN := hashString(ssn)
	hashA1 := hashString(a1)
	hashA2 := hashString(a2)

	if result.Error != nil {
		// Create new
		account = models.TradingAccount{
			UserID:      userID,
			DOB:         dobTime,
			SSNHash:     hashSSN,
			SecQ1:       q1,
			SecA1Hash:   hashA1,
			SecQ2:       q2,
			SecA2Hash:   hashA2,
			RiskAck:     true,
			Status:      "verified",
			BuyingPower: 10000.00, // Demo starting amount
		}
		config.DB.Create(&account)
	} else {
		// Update existing
		account.DOB = dobTime
		account.SSNHash = hashSSN
		account.SecQ1 = q1
		account.SecA1Hash = hashA1
		account.SecQ2 = q2
		account.SecA2Hash = hashA2
		account.RiskAck = true
		account.Status = "verified"
		config.DB.Save(&account)
	}

	return nil
}

func (s *TradingService) GetAccount(userID string) (*models.TradingAccount, error) {
	var account models.TradingAccount
	err := config.DB.Where("user_id = ?", userID).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (s *TradingService) ExecuteTrade(userID string, symbol string, quote float64, qty int, tradeType string) error {
	// Verify account
	account, err := s.GetAccount(userID)
	if err != nil || account.Status != "verified" {
		return errors.New("trading account not verified")
	}

	// Validate TradeType
	if tradeType != "buy" && tradeType != "sell" {
		return errors.New("invalid trade type")
	}

	totalAmount := quote * float64(qty)

	var position models.Position
	posErr := config.DB.Where("user_id = ? AND symbol = ?", userID, symbol).First(&position).Error

	if tradeType == "buy" {
		if account.BuyingPower < totalAmount {
			return errors.New("insufficient buying power")
		}

		account.BuyingPower -= totalAmount

		if posErr != nil {
			position = models.Position{
				UserID:   userID,
				Symbol:   symbol,
				Quantity: qty,
				AvgCost:  quote,
			}
			config.DB.Create(&position)
		} else {
			totalCost := float64(position.Quantity)*position.AvgCost + totalAmount
			position.Quantity += qty
			position.AvgCost = totalCost / float64(position.Quantity)
			config.DB.Save(&position)
		}
	} else if tradeType == "sell" {
		if posErr != nil || position.Quantity < qty {
			return errors.New("insufficient position to sell")
		}

		account.BuyingPower += totalAmount
		position.Quantity -= qty
		if position.Quantity == 0 {
			config.DB.Delete(&position)
		} else {
			config.DB.Save(&position)
		}
	}

	// Save changes
	config.DB.Save(account)

	// Record Trade
	trade := models.Trade{
		UserID:      userID,
		Symbol:      symbol,
		Type:        tradeType,
		Quantity:    qty,
		Price:       quote,
		TotalAmount: totalAmount,
	}
	config.DB.Create(&trade)

	return nil
}

func (s *TradingService) GetPortfolio(userID string) (map[string]interface{}, error) {
	var positions []models.Position
	config.DB.Where("user_id = ?", userID).Find(&positions)

	account, _ := s.GetAccount(userID)
	buyingPower := 0.0
	if account != nil {
		buyingPower = account.BuyingPower
	}

	return map[string]interface{}{
		"positions":    positions,
		"buying_power": buyingPower,
	}, nil
}
