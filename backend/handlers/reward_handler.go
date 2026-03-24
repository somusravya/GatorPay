package handlers

import (
	"net/http"
	"strconv"

	"gatorpay-backend/services"
	"gatorpay-backend/utils"

	"github.com/gin-gonic/gin"
)

// RewardHandler handles reward-related HTTP requests
type RewardHandler struct {
	rewardService *services.RewardService
}

// NewRewardHandler creates a new RewardHandler
func NewRewardHandler(rewardService *services.RewardService) *RewardHandler {
	return &RewardHandler{rewardService: rewardService}
}

// GetSummary returns reward summary for the current user
func (h *RewardHandler) GetSummary(c *gin.Context) {
	userID, _ := c.Get("userID")

	summary, err := h.rewardService.GetSummary(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Reward summary retrieved successfully", summary)
}

// GetHistory returns paginated reward history
func (h *RewardHandler) GetHistory(c *gin.Context) {
	userID, _ := c.Get("userID")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	rewards, total, err := h.rewardService.GetHistory(userID.(string), page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Reward history retrieved successfully", gin.H{
		"rewards": rewards,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

// GetOffers returns available promotional offers
func (h *RewardHandler) GetOffers(c *gin.Context) {
	offers := h.rewardService.GetOffers()
	utils.SuccessResponse(c, http.StatusOK, "Offers retrieved successfully", offers)
}
