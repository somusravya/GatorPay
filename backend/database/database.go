package database

import (
	"log"

	"gatorpay-backend/models"

	"github.com/shopspring/decimal"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database connection
var DB *gorm.DB

// Connect initializes the PostgreSQL database connection
func Connect(dsn string) {
	if dsn == "" {
		log.Fatal("DB_DSN is required. Set it in .env or as an environment variable.")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("✅ Database connected successfully (PostgreSQL)")
}

// Seed populates the database with initial data
func Seed() {
	seedBillers()
	seedLoanOffers()
}