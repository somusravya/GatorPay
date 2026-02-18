package handlers

import (
	"net/http"
	"strconv"

	"gatorpay-backend/services"
	"gatorpay-backend/utils"

	"github.com/gin-gonic/gin"
)

// WalletHandler handles wallet-related HTTP requests
type WalletHandler struct {
	walletService *services.WalletService
}

// NewWalletHandler creates a new WalletHandler
func NewWalletHandler(walletService *services.WalletService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}

// AddMoney handles adding money to the wallet
func (h *WalletHandler) AddMoney(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input services.AddMoneyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	wallet, err := h.walletService.AddMoney(userID.(string), input)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Money added successfully", wallet)
}

// Withdraw handles withdrawing money from the wallet
func (h *WalletHandler) Withdraw(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input services.WithdrawInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	wallet, err := h.walletService.Withdraw(userID.(string), input)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Withdrawal successful", wallet)
}

// GetTransactions returns paginated transactions
func (h *WalletHandler) GetTransactions(c *gin.Context) {
	userID, _ := c.Get("userID")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	response, err := h.walletService.GetTransactions(userID.(string), page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Transactions retrieved successfully", response)
}
