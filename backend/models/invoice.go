package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Invoice represents a merchant invoice
type Invoice struct {
	ID             string          `gorm:"type:varchar(36);primaryKey" json:"id"`
	MerchantID     string          `gorm:"type:varchar(36);index;not null" json:"merchant_id"`
	CustomerEmail  string          `json:"customer_email"`
	CustomerName   string          `json:"customer_name"`
	Amount         decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"amount"`
	PaidAmount     decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"paid_amount"`
	Status         string          `gorm:"type:varchar(20);default:pending" json:"status"` // "pending", "paid", "partial", "overdue", "cancelled"
	DueDate        *time.Time      `json:"due_date,omitempty"`
	Description    string          `gorm:"type:text" json:"description"`
	Items          string          `gorm:"type:text" json:"items"` // JSON string of line items
	InvoiceNumber  string          `gorm:"uniqueIndex" json:"invoice_number"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (i *Invoice) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	return nil
}

// PaymentLink represents a shareable payment link
type PaymentLink struct {
	ID          string          `gorm:"type:varchar(36);primaryKey" json:"id"`
	MerchantID  string          `gorm:"type:varchar(36);index;not null" json:"merchant_id"`
	Amount      decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"amount"`
	Description string          `json:"description"`
	LinkCode    string          `gorm:"uniqueIndex;not null" json:"link_code"`
	IsActive    bool            `gorm:"default:true" json:"is_active"`
	UsesCount   int             `gorm:"default:0" json:"uses_count"`
	MaxUses     int             `gorm:"default:0" json:"max_uses"` // 0 = unlimited
	ExpiresAt   *time.Time      `json:"expires_at,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (p *PaymentLink) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

// CreateInvoiceRequest is the request body for creating an invoice
type CreateInvoiceRequest struct {
	CustomerEmail string  `json:"customer_email" binding:"required"`
	CustomerName  string  `json:"customer_name" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	Description   string  `json:"description"`
	Items         string  `json:"items"`
	DueDate       string  `json:"due_date"`
}

// CreatePaymentLinkRequest is the request body for creating a payment link
type CreatePaymentLinkRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
	MaxUses     int     `json:"max_uses"`
}
