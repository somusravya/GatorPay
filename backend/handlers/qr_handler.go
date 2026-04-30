package handlers

import (
	"net/http"

	"gatorpay-backend/services"
	"gatorpay-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type QRHandler struct {
	qrService *services.QRService
}

func NewQRHandler(qs *services.QRService) *QRHandler {
	return &QRHandler{qrService: qs}
}

type MerchantRegRequest struct {
	BusinessName string `json:"business_name" binding:"required"`
	Category     string `json:"category" binding:"required"`
}

func (h *QRHandler) RegisterMerchant(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req MerchantRegRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid payload")
		return
	}

	merchant, err := h.qrService.RegisterMerchant(userID.(string), req.BusinessName, req.Category)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Merchant registered", merchant)
}

type QRGenerateRequest struct {
	Amount    float64 `json:"amount"`
	IsDynamic bool    `json:"is_dynamic"`
}

func (h *QRHandler) GenerateQR(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req QRGenerateRequest
	c.ShouldBindJSON(&req)

	qr, err := h.qrService.GenerateQR(userID.(string), decimal.NewFromFloat(req.Amount), req.IsDynamic)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "QR Generated", qr)
}

type QRLookupRequest struct {
	CodeString string `json:"code_string" binding:"required"`
}

func (h *QRHandler) LookupQR(c *gin.Context) {
	var req QRLookupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "QR code string is required")
		return
	}

	result, err := h.qrService.LookupQR(req.CodeString)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "QR code verified", result)
}

type QRPayRequest struct {
	CodeString string  `json:"code_string" binding:"required"`
	Amount     float64 `json:"amount"`
}

func (h *QRHandler) PayViaQR(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req QRPayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid pay payload")
		return
	}

	err := h.qrService.PayViaQR(userID.(string), req.CodeString, decimal.NewFromFloat(req.Amount))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "QR Payment Successful", nil)
}
