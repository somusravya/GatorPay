package handlers

import (
	"net/http"

	"gatorpay-backend/services"
	"gatorpay-backend/utils"

	"github.com/gin-gonic/gin"
)

// InsightHandler handles insight API requests
type InsightHandler struct {
	service *services.InsightService
}

// NewInsightHandler creates a new InsightHandler
func NewInsightHandler(service *services.InsightService) *InsightHandler {
	return &InsightHandler{service: service}
}

// GetSummary returns the AI financial insights summary for the current user
func (h *InsightHandler) GetSummary(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	summary, err := h.service.GetSummary(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate insights")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Financial insights generated", summary)
}
