package models

import (
	"gorm.io/gorm"
)

// Trade represents a historical buy/sell order
type Trade struct {
	gorm.Model
	UserID      string  `json:"user_id" gorm:"index"`
	Symbol      string  `json:"symbol" gorm:"index"`
	Type        string  `json:"type"` // "buy" or "sell"
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`        // Unit price of execution
	TotalAmount float64 `json:"total_amount"` // Total cost/proceeds
}

// Position represents current stock holdings
type Position struct {
	gorm.Model
	UserID   string  `json:"user_id" gorm:"index:idx_user_symbol,unique"`
	Symbol   string  `json:"symbol" gorm:"index:idx_user_symbol,unique"`
	Quantity int     `json:"quantity"`
	AvgCost  float64 `json:"avg_cost"`
}
