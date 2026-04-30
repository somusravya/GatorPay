package handlers

import (
	"net/http"

	"gatorpay-backend/models"
	"gatorpay-backend/services"
	"gatorpay-backend/utils"

	"github.com/gin-gonic/gin"
)

// SubscriptionHandler handles subscription API requests
type SubscriptionHandler struct {
	service *services.SubscriptionService
}

// NewSubscriptionHandler creates a new SubscriptionHandler
func NewSubscriptionHandler(service *services.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

// GetSubscriptions returns all subscriptions for the current user
func (h *SubscriptionHandler) GetSubscriptions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	subs, err := h.service.GetSubscriptions(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch subscriptions")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Subscriptions retrieved", subs)
}

// TrackSubscription manually tracks a new subscription
func (h *SubscriptionHandler) TrackSubscription(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.TrackSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	sub, err := h.service.TrackSubscription(userID.(string), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to track subscription")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Subscription tracked", sub)
}

// SetAutoPay enables/disables auto-pay for a subscription
func (h *SubscriptionHandler) SetAutoPay(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.AutoPayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	err := h.service.SetAutoPay(userID.(string), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update auto-pay")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Auto-pay updated", nil)
}
