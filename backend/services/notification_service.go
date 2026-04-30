package services

import (
	"gatorpay-backend/models"

	"gorm.io/gorm"
)

// NotificationService handles in-app notifications
type NotificationService struct {
	db *gorm.DB
}

// NewNotificationService creates a new NotificationService
func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

// GetNotifications lists notifications for a user
func (s *NotificationService) GetNotifications(userID string, filterType string) ([]models.Notification, error) {
	var notifications []models.Notification
	query := s.db.Where("user_id = ?", userID).Order("created_at desc").Limit(50)
	if filterType != "" {
		query = query.Where("type = ?", filterType)
	}
	err := query.Find(&notifications).Error
	return notifications, err
}

// CreateNotification creates a new notification
func (s *NotificationService) CreateNotification(userID, nType, title, body, icon, actionURL string) (*models.Notification, error) {
	n := models.Notification{
		UserID:    userID,
		Type:      nType,
		Title:     title,
		Body:      body,
		Icon:      icon,
		ActionURL: actionURL,
	}
	if err := s.db.Create(&n).Error; err != nil {
		return nil, err
	}
	return &n, nil
}

// MarkRead marks a notification as read
func (s *NotificationService) MarkRead(notificationID, userID string) error {
	return s.db.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true).Error
}

// MarkAllRead marks all notifications as read for a user
func (s *NotificationService) MarkAllRead(userID string) error {
	return s.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}

// GetUnreadCount returns the count of unread notifications
func (s *NotificationService) GetUnreadCount(userID string) int64 {
	var count int64
	s.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count)
	return count
}

// GetPreferences returns notification preferences for a user
func (s *NotificationService) GetPreferences(userID string) (*models.NotificationPreference, error) {
	var pref models.NotificationPreference
	err := s.db.Where("user_id = ?", userID).First(&pref).Error
	if err != nil {
		// Create default preferences
		pref = models.NotificationPreference{
			UserID:         userID,
			PaymentAlerts:  true,
			TradeAlerts:    true,
			LoanReminders:  true,
			PromoOffers:    true,
			SecurityAlerts: true,
			PriceAlerts:    false,
		}
		s.db.Create(&pref)
		return &pref, nil
	}
	return &pref, nil
}

// UpdatePreferences updates notification preferences
func (s *NotificationService) UpdatePreferences(userID string, updates map[string]interface{}) error {
	if _, err := s.GetPreferences(userID); err != nil {
		return err
	}
	return s.db.Model(&models.NotificationPreference{}).
		Where("user_id = ?", userID).
		Updates(updates).Error
}
