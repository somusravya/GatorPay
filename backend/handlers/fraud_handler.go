package handlers

import (
	"net/http"

	"gatorpay-backend/services"
	"gatorpay-backend/utils"

	"github.com/gin-gonic/gin"
)

// FraudHandler handles fraud detection API requests
type FraudHandler struct {
	service *services.FraudService
}

// NewFraudHandler creates a new FraudHandler
func NewFraudHandler(service *services.FraudService) *FraudHandler {
	return &FraudHandler{service: service}
}

// GetAlerts returns fraud alerts
func (h *FraudHandler) GetAlerts(c *gin.Context) {
	status := c.Query("status")

	alerts, err := h.service.GetAlerts(status)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch fraud alerts")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Fraud alerts retrieved", alerts)
}

// ReviewAlert handles fraud alert review
func (h *FraudHandler) ReviewAlert(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req struct {
		AlertID string `json:"alert_id" binding:"required"`
		Action  string `json:"action" binding:"required,oneof=reviewed dismissed confirmed"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	err := h.service.ReviewAlert(req.AlertID, userID.(string), req.Action)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to review alert")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Alert reviewed", nil)
}
