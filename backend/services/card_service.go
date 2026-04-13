package services

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"gatorpay-backend/models"

	"gorm.io/gorm"
)

// CardService manages virtual cards.
type CardService struct {
	db *gorm.DB
}

// NewCardService creates a CardService.
func NewCardService(db *gorm.DB) *CardService {
	return &CardService{db: db}
}

func randomDigits(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('0' + rand.Intn(10))
	}
	return string(b)
}

// CreateCard issues a new virtual card for the user.
func (s *CardService) CreateCard(userID, name string) (*models.VirtualCard, error) {
	card := models.VirtualCard{
		UserID:     userID,
		CardNumber: fmt.Sprintf("4%s", randomDigits(15)),
		CVV:        randomDigits(3),
		ExpiryDate: time.Now().AddDate(3, 0, 0).Format("01/06"),
		Name:       name,
		IsFrozen:   false,
	}
	if err := s.db.Create(&card).Error; err != nil {
		return nil, err
	}
	return &card, nil
}

// GetCards lists cards for a user (masked numbers acceptable for list).
func (s *CardService) GetCards(userID string) ([]models.VirtualCard, error) {
	var cards []models.VirtualCard
	if err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&cards).Error; err != nil {
		return nil, err
	}
	return cards, nil
}

// GetFullCardDetails returns a card with sensitive fields for an owner.
func (s *CardService) GetFullCardDetails(cardID, userID string) (*models.VirtualCard, error) {
	var card models.VirtualCard
	if err := s.db.Where("id = ? AND user_id = ?", cardID, userID).First(&card).Error; err != nil {
		return nil, errors.New("card not found")
	}
	return &card, nil
}

// ToggleFreeze flips freeze state.
func (s *CardService) ToggleFreeze(cardID, userID string) (*models.VirtualCard, error) {
	var card models.VirtualCard
	if err := s.db.Where("id = ? AND user_id = ?", cardID, userID).First(&card).Error; err != nil {
		return nil, errors.New("card not found")
	}
	card.IsFrozen = !card.IsFrozen
	if err := s.db.Save(&card).Error; err != nil {
		return nil, err
	}
	return &card, nil
}
