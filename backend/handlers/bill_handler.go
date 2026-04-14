package handlers

import (
	"net/http"

	"gatorpay-backend/services"
	"gatorpay-backend/utils"

	"github.com/gin-gonic/gin"
)

// BillHandler handles bill-related HTTP requests
type BillHandler struct {
	billService *services.BillService
}

// NewBillHandler creates a new BillHandler
func NewBillHandler(billService *services.BillService) *BillHandler {
	return &BillHandler{billService: billService}
}

// GetCategories returns all bill categories
func (h *BillHandler) GetCategories(c *gin.Context) {
	categories, err := h.billService.GetCategories()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Categories retrieved successfully", categories)
}

// GetBillers returns billers, optionally filtered by category
func (h *BillHandler) GetBillers(c *gin.Context) {
	category := c.Query("category")

	billers, err := h.billService.GetBillers(category)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Billers retrieved successfully", billers)
}

// PayBill handles bill payment
func (h *BillHandler) PayBill(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input services.BillPayInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	response, err := h.billService.PayBill(userID.(string), input)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Bill paid successfully", response)
}

// GetSavedBillers returns saved billers for the current user
func (h *BillHandler) GetSavedBillers(c *gin.Context) {
	userID, _ := c.Get("userID")

	savedBillers, err := h.billService.GetSavedBillers(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Saved billers retrieved successfully", savedBillers)
}

// RemoveSavedBiller removes a saved biller
func (h *BillHandler) RemoveSavedBiller(c *gin.Context) {
	userID, _ := c.Get("userID")
	savedBillerID := c.Param("id")

	if err := h.billService.RemoveSavedBiller(userID.(string), savedBillerID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Saved biller removed successfully", nil)
}
