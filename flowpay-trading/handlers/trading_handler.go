package handlers

import (
	"net/http"

	"flowpay-trading/services"
	"flowpay-trading/utils"

	"github.com/gin-gonic/gin"
)

type TradingHandler struct {
	tradingService *services.TradingService
}

func NewTradingHandler(ts *services.TradingService) *TradingHandler {
	return &TradingHandler{tradingService: ts}
}

type VerifyRequest struct {
	DOB     string `json:"dob" binding:"required"`
	SSN     string `json:"ssn" binding:"required"`
	SecQ1   string `json:"sec_q1" binding:"required"`
	SecA1   string `json:"sec_a1" binding:"required"`
	SecQ2   string `json:"sec_q2" binding:"required"`
	SecA2   string `json:"sec_a2" binding:"required"`
	RiskAck bool   `json:"risk_ack"`
}

func (h *TradingHandler) Verify(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err := h.tradingService.VerifyAccount(userID.(string), req.DOB, req.SSN, req.SecQ1, req.SecA1, req.SecQ2, req.SecA2, req.RiskAck)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Account verified successfully"})
}

func (h *TradingHandler) GetAccount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	account, err := h.tradingService.GetAccount(userID.(string))
	if err != nil {
		// Return pending structure if not found
		utils.SuccessResponse(c, http.StatusOK, gin.H{
			"status": "pending",
		})
		return
	}

	utils.SuccessResponse(c, http.StatusOK, account)
}

type TradeRequest struct {
	Symbol   string  `json:"symbol" binding:"required"`
	Type     string  `json:"type" binding:"required"` // "buy" or "sell"
	Quantity int     `json:"quantity" binding:"required"`
	Price    float64 `json:"price" binding:"required"` // In real-world, price would be server-verified strictly, but we accept it for mock.
}

func (h *TradingHandler) ExecuteTrade(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req TradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err := h.tradingService.ExecuteTrade(userID.(string), req.Symbol, req.Price, req.Quantity, req.Type)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Trade executed successfully"})
}

func (h *TradingHandler) GetPortfolio(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	data, err := h.tradingService.GetPortfolio(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to load portfolio")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, data)
}

func (h *TradingHandler) GetOrderHistory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	trades, err := h.tradingService.GetOrderHistory(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to load order history")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, trades)
}
