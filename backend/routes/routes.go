package routes

import (
	"gatorpay-backend/handlers"
	"gatorpay-backend/middleware"
	"gatorpay-backend/services"

	"github.com/gin-gonic/gin"
)

// Setup configures all API routes
func Setup(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	walletHandler *handlers.WalletHandler,
	tokenService *services.TokenService,
) {
	api := router.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/verify-otp", authHandler.VerifyOTP)
		auth.POST("/resend-otp", authHandler.ResendOTP)
		auth.POST("/google", authHandler.GoogleAuth)

		// Protected auth route
		auth.GET("/me", middleware.AuthMiddleware(tokenService), authHandler.GetMe)
	}

	// Wallet routes (protected)
	wallet := api.Group("/wallet")
	wallet.Use(middleware.AuthMiddleware(tokenService))
	{
		wallet.POST("/add", walletHandler.AddMoney)
		wallet.POST("/withdraw", walletHandler.Withdraw)
		wallet.GET("/transactions", walletHandler.GetTransactions)
	}
}
