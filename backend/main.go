package main

import (
	"log"
	"strings"

	"gatorpay-backend/config"
	"gatorpay-backend/database"
	"gatorpay-backend/handlers"
	"gatorpay-backend/middleware"
	"gatorpay-backend/routes"
	"gatorpay-backend/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (ignore error if not present)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Connect to database
	database.Connect(cfg.DBDSN)
	database.Migrate()
	database.Seed()

	// Initialize services
	tokenService := services.NewTokenService(cfg)
	emailService := services.NewEmailService(cfg)
	otpService := services.NewOTPService(database.DB, emailService)
	authService := services.NewAuthService(database.DB, tokenService, otpService)
	walletService := services.NewWalletService(database.DB)
	rewardService := services.NewRewardService(database.DB)
	transferService := services.NewTransferService(database.DB, rewardService)
	billService := services.NewBillService(database.DB, rewardService)
	loanService := services.NewLoanService(database.DB)
	cardService := services.NewCardService(database.DB)
	qrService := services.NewQRService(database.DB)
	statementService := services.NewStatementService(database.DB)

	// Sprint 4 services
	insightService := services.NewInsightService(database.DB)
	budgetService := services.NewBudgetService(database.DB)
	subscriptionService := services.NewSubscriptionService(database.DB)
	fraudService := services.NewFraudService(database.DB)
	notificationService := services.NewNotificationService(database.DB)
	socialService := services.NewSocialService(database.DB)
	adminService := services.NewAdminService(database.DB)
	invoiceService := services.NewInvoiceService(database.DB)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	walletHandler := handlers.NewWalletHandler(walletService)
	transferHandler := handlers.NewTransferHandler(transferService)
	billHandler := handlers.NewBillHandler(billService)
	rewardHandler := handlers.NewRewardHandler(rewardService)
	loanHandler := handlers.NewLoanHandler(loanService)
	cardHandler := handlers.NewCardHandler(cardService, otpService)
	qrHandler := handlers.NewQRHandler(qrService)
	statementHandler := handlers.NewStatementHandler(statementService)

	// Sprint 4 handlers
	insightHandler := handlers.NewInsightHandler(insightService)
	budgetHandler := handlers.NewBudgetHandler(budgetService)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)
	fraudHandler := handlers.NewFraudHandler(fraudService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	socialHandler := handlers.NewSocialHandler(socialService)
	adminHandler := handlers.NewAdminHandler(adminService)
	invoiceHandler := handlers.NewInvoiceHandler(invoiceService)

	// Setup Gin router
	router := gin.Default()

	// CORS configuration (must be before rate limiter for OPTIONS preflight)
	corsOrigins := strings.Split(cfg.CORSOrigins, ",")
	router.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Rate limiting middleware (Sprint 4)
	router.Use(middleware.RateLimiterMiddleware())

	// Setup routes
	routes.Setup(router, authHandler, walletHandler, transferHandler, billHandler, rewardHandler, tokenService, loanHandler, cardHandler, qrHandler, statementHandler,
		// Sprint 4 handlers
		insightHandler, budgetHandler, subscriptionHandler, fraudHandler, notificationHandler, socialHandler, adminHandler, invoiceHandler)

	// Start server
	log.Printf("🚀 GatorPay Backend running on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
