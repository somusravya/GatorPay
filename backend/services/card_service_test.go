package services_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	"gatorpay-backend/models"
	"gatorpay-backend/services"
)

func TestCreateCard(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&models.VirtualCard{})
	cs := services.NewCardService(db)

	card, err := cs.CreateCard("u1", "My Subscriptions")
	assert.NoError(t, err)
	assert.NotNil(t, card)
	assert.Equal(t, "My Subscriptions", card.Name)
	assert.False(t, card.IsFrozen)

	// Mask checking
	cards, err := cs.GetCards("u1")
	assert.NoError(t, err)
	assert.Len(t, cards, 1)
	
	// Toggle
	frozenCard, err := cs.ToggleFreeze(card.ID, "u1")
	assert.NoError(t, err)
	assert.True(t, frozenCard.IsFrozen)
}
