package services

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"

	"gatorpay-backend/models"

	"github.com/shopspring/decimal"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

type QRService struct {
	db *gorm.DB
}

func NewQRService(db *gorm.DB) *QRService {
	return &QRService{db: db}
}

func (s *QRService) RegisterMerchant(userID string, bizName, category string) (*models.Merchant, error) {
	merchant := models.Merchant{
		UserID:       userID,
		BusinessName: bizName,
		Category:     category,
		Status:       "active",
	}

	if err := s.db.Create(&merchant).Error; err != nil {
		return nil, errors.New("merchant already exists or invalid data")
	}
	return &merchant, nil
}

func (s *QRService) GenerateQR(merchantUserID string, amount decimal.Decimal, isDynamic bool) (*models.MerchantQRCode, error) {
	var merchant models.Merchant
	if err := s.db.Where("user_id = ?", merchantUserID).First(&merchant).Error; err != nil {
		return nil, errors.New("merchant not found")
	}

	// For mock uniqueness
	codeString := fmt.Sprintf("GATORPAY-%s-%d", merchant.ID[:8], rand.Intn(999999))

	png, err := qrcode.Encode(codeString, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}
	base64PNG := base64.StdEncoding.EncodeToString(png)

	qr := models.MerchantQRCode{
		MerchantID: merchant.ID,
		CodeString: codeString,
		Base64PNG:  "data:image/png;base64," + base64PNG,
		Amount:     amount,
		IsDynamic:  isDynamic,
	}

	if err := s.db.Create(&qr).Error; err != nil {
		return nil, err
	}
	return &qr, nil
}

// LookupQR returns QR code and merchant details for recipient verification before payment
func (s *QRService) LookupQR(codeString string) (map[string]interface{}, error) {
	var qr models.MerchantQRCode
	if err := s.db.Where("code_string = ?", codeString).First(&qr).Error; err != nil {
		return nil, errors.New("invalid or unrecognized QR code")
	}

	var merchant models.Merchant
	if err := s.db.Where("id = ?", qr.MerchantID).First(&merchant).Error; err != nil {
		return nil, errors.New("merchant not found for this QR code")
	}

	return map[string]interface{}{
		"code_string":   qr.CodeString,
		"business_name": merchant.BusinessName,
		"category":      merchant.Category,
		"merchant_id":   merchant.ID,
		"amount":        qr.Amount,
		"is_dynamic":    qr.IsDynamic,
		"status":        merchant.Status,
	}, nil
}

func (s *QRService) PayViaQR(payerUserID string, codeString string, overrideAmount decimal.Decimal) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var qr models.MerchantQRCode
		if err := tx.Where("code_string = ?", codeString).First(&qr).Error; err != nil {
			return errors.New("invalid QR code")
		}

		var merchant models.Merchant
		if err := tx.Where("id = ?", qr.MerchantID).First(&merchant).Error; err != nil {
			return errors.New("merchant no longer exists")
		}

		payAmount := qr.Amount
		if overrideAmount.GreaterThan(decimal.NewFromInt(0)) && !qr.IsDynamic {
			payAmount = overrideAmount
		}

		if payAmount.LessThanOrEqual(decimal.NewFromInt(0)) {
			return errors.New("invalid payment amount")
		}

		if merchant.UserID == payerUserID {
			return errors.New("cannot pay yourself")
		}

		var payerWallet models.Wallet
		if err := tx.Where("user_id = ?", payerUserID).First(&payerWallet).Error; err != nil {
			return errors.New("payer wallet not found")
		}

		if payerWallet.Balance.LessThan(payAmount) {
			return errors.New("insufficient balance")
		}

		var merchantWallet models.Wallet
		if err := tx.Where("user_id = ?", merchant.UserID).First(&merchantWallet).Error; err != nil {
			merchantWallet = models.Wallet{UserID: merchant.UserID, Balance: decimal.NewFromInt(0)}
			tx.Create(&merchantWallet)
		}

		merchantFee := payAmount.Mul(decimal.NewFromFloat(0.015))
		payerCashback := payAmount.Mul(decimal.NewFromFloat(0.015))
		netMerchant := payAmount.Sub(merchantFee)

		payerWallet.Balance = payerWallet.Balance.Sub(payAmount).Add(payerCashback)
		merchantWallet.Balance = merchantWallet.Balance.Add(netMerchant)

		if err := tx.Save(&payerWallet).Error; err != nil {
			return err
		}
		if err := tx.Save(&merchantWallet).Error; err != nil {
			return err
		}

		// Transaction log for payer
		payerTxn := models.Transaction{
			WalletID:    payerWallet.ID,
			Type:        "qr_payment",
			Amount:      payAmount,
			Description: fmt.Sprintf("QR payment to %s (cashback: $%.2f)", merchant.BusinessName, payerCashback.InexactFloat64()),
			Status:      models.TransactionStatusSuccess,
		}
		tx.Create(&payerTxn)

		// Transaction log for merchant
		merchantTxn := models.Transaction{
			WalletID:    merchantWallet.ID,
			Type:        "qr_received",
			Amount:      netMerchant,
			Description: fmt.Sprintf("QR payment received (fee: $%.2f)", merchantFee.InexactFloat64()),
			Status:      models.TransactionStatusSuccess,
		}
		tx.Create(&merchantTxn)

		return nil
	})
}
