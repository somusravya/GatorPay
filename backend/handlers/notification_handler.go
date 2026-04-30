package handlers

import (
	"net/http"

	"gatorpay-backend/services"
	"gatorpay-backend/utils"

	"github.com/gin-gonic/gin"
)

// NotificationHandler handles notification API requests
type NotificationHandler struct {
	service *services.NotificationService
}

// NewNotificationHandler creates a new NotificationHandler
func NewNotificationHandler(service *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// GetNotifications returns notifications for the current user
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	filterType := c.Query("type")
	notifications, err := h.service.GetNotifications(userID.(string), filterType)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch notifications")
		return
	}

	unreadCount := h.service.GetUnreadCount(userID.(string))

	utils.SuccessResponse(c, http.StatusOK, "Notifications retrieved", gin.H{
		"notifications": notifications,
		"unread_count":  unreadCount,
	})
}

// MarkRead marks a notification as read
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	notificationID := c.Param("id")
	if notificationID == "all" {
		err := h.service.MarkAllRead(userID.(string))
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to mark all read")
			return
		}
		utils.SuccessResponse(c, http.StatusOK, "All notifications marked as read", nil)
		return
	}

	err := h.service.MarkRead(notificationID, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to mark notification read")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Notification marked as read", nil)
}

// GetPreferences returns notification preferences for the current user
func (h *NotificationHandler) GetPreferences(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	preferences, err := h.service.GetPreferences(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch preferences")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Preferences retrieved", preferences)
}

// UpdatePreferences updates notification preferences
func (h *NotificationHandler) UpdatePreferences(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	err := h.service.UpdatePreferences(userID.(string), updates)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update preferences")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Preferences updated", nil)
}
