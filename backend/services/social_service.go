package services

import (
	"gatorpay-backend/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// SocialService handles social payments feed
type SocialService struct {
	db *gorm.DB
}

// NewSocialService creates a new SocialService
func NewSocialService(db *gorm.DB) *SocialService {
	return &SocialService{db: db}
}

// GetFeed returns the social feed for a user
func (s *SocialService) GetFeed(userID string) ([]models.SocialFeedItem, error) {
	var items []models.SocialFeedItem

	// Get public feed items and friend feed items
	err := s.db.Preload("User").Preload("Recipient").
		Where("privacy = ? OR user_id = ? OR recipient_id = ?", "public", userID, userID).
		Order("created_at desc").Limit(50).
		Find(&items).Error

	return items, err
}

// CreatePost creates a new social feed post
func (s *SocialService) CreatePost(userID string, req models.CreateSocialPostRequest) (*models.SocialFeedItem, error) {
	item := models.SocialFeedItem{
		UserID:      userID,
		RecipientID: req.RecipientID,
		Note:        req.Note,
		Emoji:       req.Emoji,
		Amount:      decimal.NewFromFloat(req.Amount),
		Privacy:     req.Privacy,
	}

	if item.Privacy == "" {
		item.Privacy = "friends"
	}
	if item.Emoji == "" {
		item.Emoji = "💸"
	}

	if err := s.db.Create(&item).Error; err != nil {
		return nil, err
	}

	// Reload with associations
	s.db.Preload("User").Preload("Recipient").First(&item, "id = ?", item.ID)

	return &item, nil
}

// ReactToPost adds a reaction (like/comment) to a feed item
func (s *SocialService) ReactToPost(userID string, req models.ReactToPostRequest) (*models.FeedReaction, error) {
	reaction := models.FeedReaction{
		FeedItemID: req.FeedItemID,
		UserID:     userID,
		Type:       req.Type,
		Content:    req.Content,
		Emoji:      req.Emoji,
	}

	if err := s.db.Create(&reaction).Error; err != nil {
		return nil, err
	}

	// Update counts
	if req.Type == "like" {
		s.db.Model(&models.SocialFeedItem{}).Where("id = ?", req.FeedItemID).
			Update("like_count", gorm.Expr("like_count + 1"))
	} else {
		s.db.Model(&models.SocialFeedItem{}).Where("id = ?", req.FeedItemID).
			Update("comment_count", gorm.Expr("comment_count + 1"))
	}

	// Reload with user
	s.db.Preload("User").First(&reaction, "id = ?", reaction.ID)

	return &reaction, nil
}

// GetFriends returns the friend list for a user
func (s *SocialService) GetFriends(userID string) ([]models.Friendship, error) {
	var friends []models.Friendship
	err := s.db.Preload("Friend").
		Where("user_id = ? AND status = ?", userID, "accepted").
		Find(&friends).Error
	return friends, err
}

// AddFriend sends a friend request
func (s *SocialService) AddFriend(userID string, req models.AddFriendRequest) (*models.Friendship, error) {
	// Check if already friends
	var existing models.Friendship
	err := s.db.Where("user_id = ? AND friend_id = ?", userID, req.FriendID).First(&existing).Error
	if err == nil {
		return &existing, nil
	}

	friendship := models.Friendship{
		UserID:   userID,
		FriendID: req.FriendID,
		Status:   "accepted", // Auto-accept for demo
	}

	if err := s.db.Create(&friendship).Error; err != nil {
		return nil, err
	}

	// Create reciprocal friendship
	reciprocal := models.Friendship{
		UserID:   req.FriendID,
		FriendID: userID,
		Status:   "accepted",
	}
	s.db.Create(&reciprocal)

	s.db.Preload("Friend").First(&friendship, "id = ?", friendship.ID)

	return &friendship, nil
}
