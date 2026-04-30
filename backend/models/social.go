package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// SocialFeedItem represents a payment feed entry
type SocialFeedItem struct {
	ID            string          `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID        string          `gorm:"type:varchar(36);index;not null" json:"user_id"`
	RecipientID   string          `gorm:"type:varchar(36)" json:"recipient_id"`
	TransactionID string          `gorm:"type:varchar(36)" json:"transaction_id"`
	Note          string          `gorm:"type:text" json:"note"`
	Emoji         string          `gorm:"type:varchar(10)" json:"emoji"`
	Amount        decimal.Decimal `gorm:"type:decimal(12,2)" json:"amount"`
	Privacy       string          `gorm:"type:varchar(20);default:friends" json:"privacy"` // "public", "friends", "private"
	LikeCount     int             `gorm:"default:0" json:"like_count"`
	CommentCount  int             `gorm:"default:0" json:"comment_count"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`
	User          *User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Recipient     *User           `gorm:"foreignKey:RecipientID" json:"recipient,omitempty"`
}

func (s *SocialFeedItem) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

// Friendship represents a friend connection between two users
type Friendship struct {
	ID        string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID    string         `gorm:"type:varchar(36);index;not null" json:"user_id"`
	FriendID  string         `gorm:"type:varchar(36);index;not null" json:"friend_id"`
	Status    string         `gorm:"type:varchar(20);default:pending" json:"status"` // "pending", "accepted", "blocked"
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Friend    *User          `gorm:"foreignKey:FriendID" json:"friend,omitempty"`
}

func (f *Friendship) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	return nil
}

// FeedReaction represents a reaction (like/comment) on a feed item
type FeedReaction struct {
	ID         string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	FeedItemID string         `gorm:"type:varchar(36);index;not null" json:"feed_item_id"`
	UserID     string         `gorm:"type:varchar(36);index;not null" json:"user_id"`
	Type       string         `gorm:"type:varchar(20);not null" json:"type"` // "like", "comment"
	Content    string         `gorm:"type:text" json:"content,omitempty"`    // for comments
	Emoji      string         `gorm:"type:varchar(10)" json:"emoji,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	User       *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (fr *FeedReaction) BeforeCreate(tx *gorm.DB) error {
	if fr.ID == "" {
		fr.ID = uuid.New().String()
	}
	return nil
}

// CreateSocialPostRequest is the request for creating a social feed post
type CreateSocialPostRequest struct {
	RecipientID string  `json:"recipient_id"`
	Note        string  `json:"note" binding:"required"`
	Emoji       string  `json:"emoji"`
	Amount      float64 `json:"amount"`
	Privacy     string  `json:"privacy"`
}

// ReactToPostRequest is the request for reacting to a feed item
type ReactToPostRequest struct {
	FeedItemID string `json:"feed_item_id" binding:"required"`
	Type       string `json:"type" binding:"required,oneof=like comment"`
	Content    string `json:"content"`
	Emoji      string `json:"emoji"`
}

// AddFriendRequest is the request for adding a friend
type AddFriendRequest struct {
	FriendID string `json:"friend_id" binding:"required"`
}
