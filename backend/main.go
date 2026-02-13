package main

import (
	"log"
	"strings"

	"gatorpay-backend/config"
	"gatorpay-backend/database"
	"gatorpay-backend/handlers"
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

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	walletHandler := handlers.NewWalletHandler(walletService)

	// Setup Gin router
	router := gin.Default()

	// CORS configuration
	corsOrigins := strings.Split(cfg.CORSOrigins, ",")
	router.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Setup routes
	routes.Setup(router, authHandler, walletHandler, tokenService)

	// Start server
	log.Printf("🚀 GatorPay Backend running on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
