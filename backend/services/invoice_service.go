package services

import (
	"fmt"
	"math/rand"
	"time"

	"gatorpay-backend/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// InvoiceService handles merchant invoicing and payment links
type InvoiceService struct {
	db *gorm.DB
}

// NewInvoiceService creates a new InvoiceService
func NewInvoiceService(db *gorm.DB) *InvoiceService {
	return &InvoiceService{db: db}
}

// CreateInvoice creates a new invoice
func (s *InvoiceService) CreateInvoice(merchantID string, req models.CreateInvoiceRequest) (*models.Invoice, error) {
	inv := models.Invoice{
		MerchantID:    merchantID,
		CustomerEmail: req.CustomerEmail,
		CustomerName:  req.CustomerName,
		Amount:        decimal.NewFromFloat(req.Amount),
		Description:   req.Description,
		Items:         req.Items,
		InvoiceNumber: generateInvoiceNumber(),
		Status:        "pending",
	}

	if req.DueDate != "" {
		t, err := time.Parse("2006-01-02", req.DueDate)
		if err == nil {
			inv.DueDate = &t
		}
	}

	if err := s.db.Create(&inv).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

// GetInvoices returns invoices for a merchant
func (s *InvoiceService) GetInvoices(merchantID string) ([]models.Invoice, error) {
	var invoices []models.Invoice
	err := s.db.Where("merchant_id = ?", merchantID).Order("created_at desc").Find(&invoices).Error
	return invoices, err
}

// CreatePaymentLink creates a shareable payment link
func (s *InvoiceService) CreatePaymentLink(merchantID string, req models.CreatePaymentLinkRequest) (*models.PaymentLink, error) {
	link := models.PaymentLink{
		MerchantID:  merchantID,
		Amount:      decimal.NewFromFloat(req.Amount),
		Description: req.Description,
		LinkCode:    generateLinkCode(),
		MaxUses:     req.MaxUses,
		IsActive:    true,
	}

	if err := s.db.Create(&link).Error; err != nil {
		return nil, err
	}
	return &link, nil
}

// GetPaymentLinks returns payment links for a merchant
func (s *InvoiceService) GetPaymentLinks(merchantID string) ([]models.PaymentLink, error) {
	var links []models.PaymentLink
	err := s.db.Where("merchant_id = ?", merchantID).Order("created_at desc").Find(&links).Error
	return links, err
}

func generateInvoiceNumber() string {
	return fmt.Sprintf("INV-%d-%04d", time.Now().Year(), rand.Intn(10000))
}

func generateLinkCode() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 12)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}
	return string(code)
}
