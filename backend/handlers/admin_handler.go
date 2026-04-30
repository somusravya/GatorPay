package handlers

import (
	"net/http"
	"strconv"

	"gatorpay-backend/services"
	"gatorpay-backend/utils"

	"github.com/gin-gonic/gin"
)

// AdminHandler handles admin analytics and management API requests
type AdminHandler struct {
	service *services.AdminService
}

// NewAdminHandler creates a new AdminHandler
func NewAdminHandler(service *services.AdminService) *AdminHandler {
	return &AdminHandler{service: service}
}

// GetMetrics returns platform-wide analytics metrics
func (h *AdminHandler) GetMetrics(c *gin.Context) {
	metrics := h.service.GetMetrics()
	utils.SuccessResponse(c, http.StatusOK, "Platform metrics retrieved", metrics)
}

// GetUsers returns paginated user list for admin
func (h *AdminHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")

	users, total, err := h.service.GetUsers(page, limit, search)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Users retrieved", gin.H{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetFraudReview returns the fraud review queue
func (h *AdminHandler) GetFraudReview(c *gin.Context) {
	alerts, err := h.service.GetFraudReviewQueue()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch fraud review queue")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Fraud review queue retrieved", alerts)
}
