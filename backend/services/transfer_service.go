package services

import (
	"errors"

	"gatorpay-backend/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// TransferService handles P2P transfers and user search
type TransferService struct {
	db            *gorm.DB
	rewardService *RewardService
}

// NewTransferService creates a new TransferService
func NewTransferService(db *gorm.DB, rewardService *RewardService) *TransferService {
	return &TransferService{db: db, rewardService: rewardService}
}

// TransferRequest is the DTO for sending money
type TransferRequest struct {
	Recipient string  `json:"recipient" binding:"required"` // username, email, or phone
	Amount    float64 `json:"amount" binding:"required"`
	Note      string  `json:"note"`
}

// TransferResponse is the response after a successful transfer
type TransferResponse struct {
	TransactionID string              `json:"transaction_id"`
	Recipient     models.UserResponse `json:"recipient"`
	Amount        decimal.Decimal     `json:"amount"`
	Note          string              `json:"note"`
	NewBalance    decimal.Decimal     `json:"new_balance"`
}

// SendMoney transfers money from one user to another
func (s *TransferService) SendMoney(senderID string, input TransferRequest) (*TransferResponse, error) {
	amount := decimal.NewFromFloat(input.Amount)
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("amount must be greater than 0")
	}

	// Resolve recipient by username, email, or phone
	var recipient models.User
	err := s.db.Where("username = ? OR email = ? OR phone = ?", input.Recipient, input.Recipient, input.Recipient).
		First(&recipient).Error
	if err != nil {
		return nil, errors.New("recipient not found")
	}

	// Prevent self-transfer
	if recipient.ID == senderID {
		return nil, errors.New("cannot send money to yourself")
	}

	var response TransferResponse
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Get sender wallet
		var senderWallet models.Wallet
		if err := tx.Where("user_id = ?", senderID).First(&senderWallet).Error; err != nil {
			return errors.New("sender wallet not found")
		}

		if !senderWallet.IsActive {
			return errors.New("sender wallet is not active")
		}

		// Check sufficient balance
		if senderWallet.Balance.LessThan(amount) {
			return errors.New("insufficient balance")
		}

		// Get recipient wallet
		var recipientWallet models.Wallet
		if err := tx.Where("user_id = ?", recipient.ID).First(&recipientWallet).Error; err != nil {
			return errors.New("recipient wallet not found")
		}

		// Debit sender
		senderWallet.Balance = senderWallet.Balance.Sub(amount)
		if err := tx.Save(&senderWallet).Error; err != nil {
			return errors.New("failed to debit sender")
		}

		// Credit recipient
		recipientWallet.Balance = recipientWallet.Balance.Add(amount)
		if err := tx.Save(&recipientWallet).Error; err != nil {
			return errors.New("failed to credit recipient")
		}

		description := "Transfer to " + recipient.Username
		if input.Note != "" {
			description += " - " + input.Note
		}

		// Create p2p_send transaction (sender's record)
		sendTx := models.Transaction{
			WalletID:    senderWallet.ID,
			FromUserID:  &senderID,
			ToUserID:    &recipient.ID,
			Type:        models.TransactionTypeP2PSend,
			Amount:      amount,
			Description: description,
			Status:      models.TransactionStatusSuccess,
		}
		if err := tx.Create(&sendTx).Error; err != nil {
			return errors.New("failed to create send transaction")
		}

		// Create p2p_receive transaction (recipient's record)
		var sender models.User
		tx.Where("id = ?", senderID).First(&sender)

		receiveDesc := "Received from " + sender.Username
		if input.Note != "" {
			receiveDesc += " - " + input.Note
		}

		receiveTx := models.Transaction{
			WalletID:    recipientWallet.ID,
			FromUserID:  &senderID,
			ToUserID:    &recipient.ID,
			Type:        models.TransactionTypeP2PReceive,
			Amount:      amount,
			Description: receiveDesc,
			Status:      models.TransactionStatusSuccess,
		}
		if err := tx.Create(&receiveTx).Error; err != nil {
			return errors.New("failed to create receive transaction")
		}

		response = TransferResponse{
			TransactionID: sendTx.ID,
			Recipient:     recipient.ToResponse(),
			Amount:        amount,
			Note:          input.Note,
			NewBalance:    senderWallet.Balance,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Award 1% cashback asynchronously
	go s.rewardService.AwardCashback(senderID, amount, 0.01, "Cashback for P2P transfer to "+recipient.Username)

	return &response, nil
}

// GetRecentContacts returns the 10 most recent unique recipients
func (s *TransferService) GetRecentContacts(userID string) ([]models.UserResponse, error) {
	var transactions []models.Transaction
	if err := s.db.Where("from_user_id = ? AND type = ?", userID, models.TransactionTypeP2PSend).
		Order("created_at DESC").
		Preload("ToUser").
		Find(&transactions).Error; err != nil {
		return nil, errors.New("failed to fetch recent contacts")
	}

	// Deduplicate recipients using a map
	seen := make(map[string]bool)
	var contacts []models.UserResponse

	for _, tx := range transactions {
		if tx.ToUser != nil && tx.ToUserID != nil && !seen[*tx.ToUserID] {
			seen[*tx.ToUserID] = true
			contacts = append(contacts, tx.ToUser.ToResponse())
			if len(contacts) >= 10 {
				break
			}
		}
	}

	return contacts, nil
}

// SearchUsers searches for users by username, email, or full name
func (s *TransferService) SearchUsers(currentUserID, query string) ([]models.UserResponse, error) {
	if query == "" {
		return []models.UserResponse{}, nil
	}

	pattern := "%" + query + "%"
	var users []models.User
	if err := s.db.Where(
		"(username ILIKE ? OR email ILIKE ? OR CONCAT(first_name, ' ', last_name) ILIKE ?) AND id != ?",
		pattern, pattern, pattern, currentUserID,
	).Limit(10).Find(&users).Error; err != nil {
		return nil, errors.New("failed to search users")
	}

	var results []models.UserResponse
	for _, u := range users {
		results = append(results, u.ToResponse())
	}

	return results, nil
}
