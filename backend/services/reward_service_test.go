package services

import (
	"testing"
)

func TestGetOffers(t *testing.T) {
	service := NewRewardService(nil)
	offers := service.GetOffers()

	if len(offers) != 4 {
		t.Errorf("expected 4 offers, got %d", len(offers))
	}
}

func TestGetOffersContent(t *testing.T) {
	service := NewRewardService(nil)
	offers := service.GetOffers()

	expectedTitles := []string{"Send & Save", "Bill Pay Bonus", "Points Multiplier", "Refer a Friend"}
	for i, title := range expectedTitles {
		if offers[i].Title != title {
			t.Errorf("offer %d: expected title '%s', got '%s'", i, title, offers[i].Title)
		}
	}
}

func TestGetOffersAllActive(t *testing.T) {
	service := NewRewardService(nil)
	offers := service.GetOffers()

	for i, offer := range offers {
		if !offer.IsActive {
			t.Errorf("offer %d (%s) expected to be active", i, offer.Title)
		}
	}
}

func TestGetOffersHaveIDs(t *testing.T) {
	service := NewRewardService(nil)
	offers := service.GetOffers()

	for i, offer := range offers {
		if offer.ID == "" {
			t.Errorf("offer %d expected to have an ID", i)
		}
	}
}

func TestGetOffersHaveDescriptions(t *testing.T) {
	service := NewRewardService(nil)
	offers := service.GetOffers()

	for i, offer := range offers {
		if offer.Description == "" {
			t.Errorf("offer %d expected to have a description", i)
		}
	}
}

func TestGetOffersHaveDiscounts(t *testing.T) {
	service := NewRewardService(nil)
	offers := service.GetOffers()

	for i, offer := range offers {
		if offer.Discount == "" {
			t.Errorf("offer %d expected to have a discount", i)
		}
	}
}

func TestGetOffersHaveIcons(t *testing.T) {
	service := NewRewardService(nil)
	offers := service.GetOffers()

	for i, offer := range offers {
		if offer.Icon == "" {
			t.Errorf("offer %d expected to have an icon", i)
		}
	}
}

func TestNewRewardService(t *testing.T) {
	service := NewRewardService(nil)
	if service == nil {
		t.Error("expected non-nil RewardService")
	}
}
