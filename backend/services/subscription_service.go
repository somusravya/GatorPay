package services

import (
	"strings"
	"time"

	"gatorpay-backend/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// SubscriptionService handles subscription management
type SubscriptionService struct {
	db *gorm.DB
}

// NewSubscriptionService creates a new SubscriptionService
func NewSubscriptionService(db *gorm.DB) *SubscriptionService {
	return &SubscriptionService{db: db}
}

// GetSubscriptions returns all subscriptions for a user
func (s *SubscriptionService) GetSubscriptions(userID string) ([]models.Subscription, error) {
	var subs []models.Subscription
	err := s.db.Where("user_id = ?", userID).Order("next_renewal asc").Find(&subs).Error

	// If no subscriptions exist, provide detected ones from transaction patterns
	if len(subs) == 0 {
		subs = s.detectSubscriptions(userID)
	}

	return subs, err
}

// TrackSubscription manually tracks a subscription
func (s *SubscriptionService) TrackSubscription(userID string, req models.TrackSubscriptionRequest) (*models.Subscription, error) {
	nextRenewal := time.Now().AddDate(0, 1, 0) // Default: 1 month from now
	switch strings.ToLower(req.Frequency) {
	case "weekly":
		nextRenewal = time.Now().AddDate(0, 0, 7)
	case "yearly", "annual", "annually":
		nextRenewal = time.Now().AddDate(1, 0, 0)
	}

	sub := models.Subscription{
		UserID:      userID,
		Name:        req.Name,
		Category:    req.Category,
		Amount:      decimal.NewFromFloat(req.Amount),
		Frequency:   req.Frequency,
		NextRenewal: &nextRenewal,
		Icon:        req.Icon,
		Color:       req.Color,
		Provider:    req.Provider,
		Status:      "active",
	}

	if sub.Icon == "" {
		sub.Icon = "🔄"
	}
	if sub.Color == "" {
		sub.Color = "#8b5cf6"
	}

	if err := s.db.Create(&sub).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}

// SetAutoPay enables or disables auto-pay for a subscription
func (s *SubscriptionService) SetAutoPay(userID string, req models.AutoPayRequest) error {
	return s.db.Model(&models.Subscription{}).
		Where("id = ? AND user_id = ?", req.SubscriptionID, userID).
		Update("auto_pay", req.AutoPay).Error
}

// detectSubscriptions identifies repeated charges and falls back to demo detections.
func (s *SubscriptionService) detectSubscriptions(userID string) []models.Subscription {
	var transactions []models.Transaction
	since := time.Now().AddDate(0, -4, 0)
	s.db.Where("from_user_id = ? AND created_at >= ? AND type IN ?", userID, since, []string{"withdraw", "bill_pay"}).
		Order("created_at asc").
		Find(&transactions)

	type recurringCharge struct {
		name        string
		category    string
		provider    string
		amount      float64
		firstSeen   time.Time
		lastSeen    time.Time
		occurrences int
	}

	charges := map[string]*recurringCharge{}
	for _, tx := range transactions {
		amount, _ := tx.Amount.Float64()
		if amount <= 0 {
			continue
		}
		name, category := normalizeSubscriptionMerchant(tx.Description)
		if name == "" {
			continue
		}
		key := strings.ToLower(name) + ":" + decimal.NewFromFloat(amount).Round(0).String()
		charge, ok := charges[key]
		if !ok {
			charge = &recurringCharge{name: name, category: category, provider: name, amount: amount, firstSeen: tx.CreatedAt}
			charges[key] = charge
		}
		charge.occurrences++
		charge.lastSeen = tx.CreatedAt
		charge.amount = amount
	}

	var detected []models.Subscription
	for _, charge := range charges {
		if charge.occurrences < 2 || charge.lastSeen.Sub(charge.firstSeen) < 20*24*time.Hour {
			continue
		}
		nextRenewal := charge.lastSeen.AddDate(0, 1, 0)
		detected = append(detected, models.Subscription{
			ID:          "detected-" + strings.NewReplacer(" ", "-", ".", "", "&", "and").Replace(strings.ToLower(charge.name)),
			UserID:      userID,
			Name:        charge.name,
			Category:    charge.category,
			Amount:      decimal.NewFromFloat(charge.amount),
			Frequency:   "monthly",
			NextRenewal: &nextRenewal,
			Icon:        subscriptionIcon(charge.category),
			Color:       subscriptionColor(charge.category),
			Status:      "detected",
			Provider:    charge.provider,
		})
	}
	if len(detected) > 0 {
		return detected
	}

	now := time.Now()
	nextMonth := now.AddDate(0, 1, 0)
	nextWeek := now.AddDate(0, 0, 7)

	return []models.Subscription{
		{
			ID: "detected-1", UserID: userID, Name: "Netflix", Category: "streaming",
			Amount: decimal.NewFromFloat(15.99), Frequency: "monthly",
			NextRenewal: &nextMonth, Icon: "🎬", Color: "#e50914", Status: "active", Provider: "Netflix Inc.",
		},
		{
			ID: "detected-2", UserID: userID, Name: "Spotify", Category: "music",
			Amount: decimal.NewFromFloat(9.99), Frequency: "monthly",
			NextRenewal: &nextWeek, Icon: "🎵", Color: "#1db954", Status: "active", Provider: "Spotify AB",
		},
		{
			ID: "detected-3", UserID: userID, Name: "iCloud Storage", Category: "cloud",
			Amount: decimal.NewFromFloat(2.99), Frequency: "monthly",
			NextRenewal: &nextMonth, Icon: "☁️", Color: "#007aff", Status: "active", Provider: "Apple Inc.",
		},
		{
			ID: "detected-4", UserID: userID, Name: "Gym Membership", Category: "fitness",
			Amount: decimal.NewFromFloat(29.99), Frequency: "monthly",
			NextRenewal: &nextMonth, Icon: "💪", Color: "#ff6b6b", Status: "active", Provider: "Planet Fitness",
		},
		{
			ID: "detected-5", UserID: userID, Name: "Adobe Creative Cloud", Category: "software",
			Amount: decimal.NewFromFloat(54.99), Frequency: "monthly",
			NextRenewal: &nextMonth, Icon: "🎨", Color: "#ff0000", Status: "active", Provider: "Adobe Systems",
		},
	}
}

func normalizeSubscriptionMerchant(description string) (string, string) {
	desc := strings.ToLower(description)
	catalog := []struct {
		name     string
		category string
		terms    []string
	}{
		{"Netflix", "streaming", []string{"netflix"}},
		{"Spotify", "music", []string{"spotify"}},
		{"iCloud Storage", "cloud", []string{"icloud", "apple storage"}},
		{"Adobe Creative Cloud", "software", []string{"adobe"}},
		{"YouTube Premium", "streaming", []string{"youtube premium"}},
		{"Gym Membership", "fitness", []string{"gym", "fitness", "planet fitness"}},
		{"Amazon Prime", "shopping", []string{"amazon prime", "prime membership"}},
	}
	for _, item := range catalog {
		for _, term := range item.terms {
			if strings.Contains(desc, term) {
				return item.name, item.category
			}
		}
	}
	return "", ""
}

func subscriptionIcon(category string) string {
	icons := map[string]string{
		"streaming": "🎬",
		"music":     "🎵",
		"cloud":     "☁️",
		"fitness":   "💪",
		"software":  "🎨",
		"shopping":  "🛍️",
	}
	if icon, ok := icons[category]; ok {
		return icon
	}
	return "🔄"
}

func subscriptionColor(category string) string {
	colors := map[string]string{
		"streaming": "#e50914",
		"music":     "#1db954",
		"cloud":     "#007aff",
		"fitness":   "#ff6b6b",
		"software":  "#8b5cf6",
		"shopping":  "#f59e0b",
	}
	if color, ok := colors[category]; ok {
		return color
	}
	return "#8b5cf6"
}
