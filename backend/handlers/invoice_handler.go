package handlers

import (
	"net/http"

	"gatorpay-backend/models"
	"gatorpay-backend/services"
	"gatorpay-backend/utils"

	"github.com/gin-gonic/gin"
)

// InvoiceHandler handles merchant invoice API requests
type InvoiceHandler struct {
	service *services.InvoiceService
}

// NewInvoiceHandler creates a new InvoiceHandler
func NewInvoiceHandler(service *services.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{service: service}
}

// CreateInvoice creates a new invoice
func (h *InvoiceHandler) CreateInvoice(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	invoice, err := h.service.CreateInvoice(userID.(string), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create invoice")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Invoice created", invoice)
}

// GetInvoices returns invoices for the current merchant
func (h *InvoiceHandler) GetInvoices(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	invoices, err := h.service.GetInvoices(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch invoices")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Invoices retrieved", invoices)
}

// CreatePaymentLink creates a shareable payment link
func (h *InvoiceHandler) CreatePaymentLink(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.CreatePaymentLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	link, err := h.service.CreatePaymentLink(userID.(string), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create payment link")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Payment link created", link)
}

// GetPaymentLinks returns shareable payment links for the current merchant
func (h *InvoiceHandler) GetPaymentLinks(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	links, err := h.service.GetPaymentLinks(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch payment links")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payment links retrieved", links)
}
