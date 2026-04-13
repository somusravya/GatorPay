package handlers

import (
	"net/http"

	"gatorpay-backend/services"
	"gatorpay-backend/utils"

	"github.com/gin-gonic/gin"
)

type CardHandler struct {
	cardService *services.CardService
	otpService  *services.OTPService
}

func NewCardHandler(cs *services.CardService, os *services.OTPService) *CardHandler {
	return &CardHandler{cardService: cs, otpService: os}
}

type CreateCardRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *CardHandler) CreateCard(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	card, err := h.cardService.CreateCard(userID.(string), req.Name)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Card created", card)
}

func (h *CardHandler) GetCards(c *gin.Context) {
	userID, _ := c.Get("userID")
	cards, err := h.cardService.GetCards(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to load cards")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Cards loaded", cards)
}

type OTPRequest struct {
	OTP string `json:"otp" binding:"required"`
}

func (h *CardHandler) GetCardDetails(c *gin.Context) {
	userID, _ := c.Get("userID")
	cardID := c.Param("id")

	var req OTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "OTP required")
		return
	}

	if req.OTP != "123456" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid OTP")
		return
	}

	card, err := h.cardService.GetFullCardDetails(cardID, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Card unlocked", card)
}

func (h *CardHandler) RequestOTP(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "OTP 123456 sent to email", nil)
}

func (h *CardHandler) FreezeCard(c *gin.Context) {
	userID, _ := c.Get("userID")
	cardID := c.Param("id")

	card, err := h.cardService.ToggleFreeze(cardID, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Card status updated", card)
}
