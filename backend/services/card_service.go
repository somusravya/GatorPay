package services

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"gatorpay-backend/models"

	"gorm.io/gorm"
)

type CardService struct {
	db *gorm.DB
}

func NewCardService(db *gorm.DB) *CardService {
	rand.Seed(time.Now().UnixNano())
	return &CardService{db: db}
}

func generate16Digit() string {
	return fmt.Sprintf("4%015d", rand.Int63n(1000000000000000))
}

func generateCVV() string {
	return fmt.Sprintf("%03d", rand.Intn(1000))
}

func (s *CardService) CreateCard(userID string, name string) (*models.VirtualCard, error) {
	var count int64
	s.db.Model(&models.VirtualCard{}).Where("user_id = ?", userID).Count(&count)
	if count >= 3 {
		return nil, errors.New("maximum of 3 virtual cards allowed")
	}

	expDate := time.Now().AddDate(3, 0, 0).Format("01/06")

	card := models.VirtualCard{
		UserID:     userID,
		CardNumber: generate16Digit(),
		CVV:        generateCVV(),
		ExpiryDate: expDate,
		Name:       name,
		IsFrozen:   false,
	}

	if err := s.db.Create(&card).Error; err != nil {
		return nil, err
	}
	return &card, nil
}

func (s *CardService) GetCards(userID string) ([]map[string]interface{}, error) {
	var cards []models.VirtualCard
	if err := s.db.Where("user_id = ?", userID).Find(&cards).Error; err != nil {
		return nil, err
	}

	var maskedCards []map[string]interface{}
	for _, c := range cards {
		maskedCards = append(maskedCards, map[string]interface{}{
			"id":          c.ID,
			"card_number": "**** **** **** " + c.CardNumber[12:],
			"expiry_date": c.ExpiryDate,
			"name":        c.Name,
			"is_frozen":   c.IsFrozen,
		})
	}
	return maskedCards, nil
}

func (s *CardService) GetFullCardDetails(cardID string, userID string) (*models.VirtualCard, error) {
	var card models.VirtualCard
	if err := s.db.Where("id = ? AND user_id = ?", cardID, userID).First(&card).Error; err != nil {
		return nil, errors.New("card not found")
	}
	return &card, nil
}

func (s *CardService) ToggleFreeze(cardID string, userID string) (*models.VirtualCard, error) {
	var card models.VirtualCard
	if err := s.db.Where("id = ? AND user_id = ?", cardID, userID).First(&card).Error; err != nil {
		return nil, errors.New("card not found")
	}
	card.IsFrozen = !card.IsFrozen
	s.db.Save(&card)
	return &card, nil
}
