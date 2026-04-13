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

// Migrate runs auto-migrations for all models
func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Wallet{},
		&models.Transaction{},
		&models.Biller{},
		&models.LoanOffer{},
		&models.Loan{},
		&models.VirtualCard{},
		&models.Merchant{},
		&models.MerchantQRCode{},
		&models.OTPCode{},
		&models.Reward{},
		&models.SavedBiller{},
		&models.BillPayment{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("✅ Database migrated successfully")
}

// Seed populates the database with initial data
func Seed() {
	seedBillers()
	seedLoanOffers()
}

func seedBillers() {
	var count int64
	DB.Model(&models.Biller{}).Count(&count)
	if count > 0 {
		return
	}

	billers := []models.Biller{
		{Name: "Florida Power & Light", Category: "electricity", Icon: "⚡"},
		{Name: "AT&T", Category: "internet", Icon: "🌐"},
		{Name: "T-Mobile", Category: "phone", Icon: "📱"},
		{Name: "Netflix", Category: "streaming", Icon: "🎬"},
		{Name: "Spotify", Category: "streaming", Icon: "🎵"},
		{Name: "State Farm Insurance", Category: "insurance", Icon: "🛡️"},
		{Name: "Gainesville Water", Category: "water", Icon: "💧"},
		{Name: "Xfinity", Category: "cable", Icon: "📺"},
	}

	for _, b := range billers {
		DB.Create(&b)
	}
	log.Println("✅ Seeded 8 billers")
}

func seedLoanOffers() {
	var count int64
	DB.Model(&models.LoanOffer{}).Count(&count)
	if count > 0 {
		return
	}

	offers := []models.LoanOffer{
		{
			Name:         "Quick Cash Loan",
			Description:  "Short-term loan for immediate needs",
			MinAmount:    decimal.NewFromInt(100),
			MaxAmount:    decimal.NewFromInt(1000),
			InterestRate: decimal.NewFromFloat(5.99),
			TermMonths:   3,
		},
		{
			Name:         "Student Loan",
			Description:  "Low-interest loan for Gator students",
			MinAmount:    decimal.NewFromInt(500),
			MaxAmount:    decimal.NewFromInt(5000),
			InterestRate: decimal.NewFromFloat(3.49),
			TermMonths:   12,
		},
		{
			Name:         "Personal Loan",
			Description:  "Flexible personal financing",
			MinAmount:    decimal.NewFromInt(1000),
			MaxAmount:    decimal.NewFromInt(10000),
			InterestRate: decimal.NewFromFloat(7.99),
			TermMonths:   24,
		},
	}

	for _, o := range offers {
		DB.Create(&o)
	}
	log.Println("✅ Seeded 3 loan offers")
}
