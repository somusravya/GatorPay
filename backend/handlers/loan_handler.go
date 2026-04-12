package handlers

import (
	"log"
	"net/http"

	"gatorpay-backend/services"
	"gatorpay-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type LoanHandler struct {
	loanService *services.LoanService
}

func NewLoanHandler(ls *services.LoanService) *LoanHandler {
	return &LoanHandler{loanService: ls}
}

func (h *LoanHandler) GetOffers(c *gin.Context) {
	offers, err := h.loanService.GetOffers()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to load offers")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Loaded offers", offers)
}

type ApplyLoanRequest struct {
	OfferID    string  `json:"offer_id" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	TermMonths int     `json:"term_months" binding:"required"`
}

func (h *LoanHandler) ApplyForLoan(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req ApplyLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	loan, err := h.loanService.ApplyForLoan(userID.(string), req.OfferID, decimal.NewFromFloat(req.Amount), req.TermMonths)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Loan created", loan)
}

func (h *LoanHandler) GetUserLoans(c *gin.Context) {
	userID, _ := c.Get("userID")
	loans, err := h.loanService.GetUserLoans(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve loans")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Loans retrieved", loans)
}

func (h *LoanHandler) GetLoan(c *gin.Context) {
	userID, _ := c.Get("userID")
	loanID := c.Param("id")

	loan, err := h.loanService.GetLoan(loanID, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Loan not found")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Loan retrieved", loan)
}

func (h *LoanHandler) PayEMI(c *gin.Context) {
	userID, _ := c.Get("userID")
	loanID := c.Param("id")

	err := h.loanService.MakeLoanPayment(loanID, userID.(string))
	if err != nil {
		log.Printf("[Loan EMI] Payment error for loan %s: %v", loanID, err)
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "EMI payment successful", nil)
}
