package handlers

import (
	"net/http"

	"gatorpay-backend/services"
	"gatorpay-backend/utils"

	"github.com/gin-gonic/gin"
)

// TransferHandler handles transfer-related HTTP requests
type TransferHandler struct {
	transferService *services.TransferService
}

// NewTransferHandler creates a new TransferHandler
func NewTransferHandler(transferService *services.TransferService) *TransferHandler {
	return &TransferHandler{transferService: transferService}
}

// SendMoney handles P2P money transfer
func (h *TransferHandler) SendMoney(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input services.TransferRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	response, err := h.transferService.SendMoney(userID.(string), input)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Money sent successfully", response)
}

// GetRecentContacts returns recent transfer recipients
func (h *TransferHandler) GetRecentContacts(c *gin.Context) {
	userID, _ := c.Get("userID")

	contacts, err := h.transferService.GetRecentContacts(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Recent contacts retrieved successfully", contacts)
}

// SearchUsers searches for users by query
func (h *TransferHandler) SearchUsers(c *gin.Context) {
	userID, _ := c.Get("userID")
	query := c.Query("query")

	results, err := h.transferService.SearchUsers(userID.(string), query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Search results", results)
}
