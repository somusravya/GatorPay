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
	transferHandler *handlers.TransferHandler,
	billHandler *handlers.BillHandler,
	rewardHandler *handlers.RewardHandler,
	tokenService *services.TokenService,
	loanHandler *handlers.LoanHandler,
	cardHandler *handlers.CardHandler,
	qrHandler *handlers.QRHandler,
	statementHandler *handlers.StatementHandler,
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
		wallet.GET("/statement", statementHandler.GetStatement)
	}

	// Transfer routes (protected)
	transfer := api.Group("/transfer")
	transfer.Use(middleware.AuthMiddleware(tokenService))
	{
		transfer.POST("/send", transferHandler.SendMoney)
		transfer.GET("/contacts", transferHandler.GetRecentContacts)
		transfer.GET("/search", transferHandler.SearchUsers)
	}

	// Bills routes (protected)
	bills := api.Group("/bills")
	bills.Use(middleware.AuthMiddleware(tokenService))
	{
		bills.GET("/categories", billHandler.GetCategories)
		bills.GET("/billers", billHandler.GetBillers)
		bills.POST("/pay", billHandler.PayBill)
		bills.GET("/saved", billHandler.GetSavedBillers)
		bills.DELETE("/saved/:id", billHandler.RemoveSavedBiller)
	}

	// Rewards routes (protected)
	rewards := api.Group("/rewards")
	rewards.Use(middleware.AuthMiddleware(tokenService))
	{
		rewards.GET("", rewardHandler.GetSummary)
		rewards.GET("/history", rewardHandler.GetHistory)
		rewards.GET("/offers", rewardHandler.GetOffers)
	}

	// Loans routes (protected)
	loans := api.Group("/loans")
	loans.Use(middleware.AuthMiddleware(tokenService))
	{
		loans.GET("/offers", loanHandler.GetOffers)
		loans.POST("/apply", loanHandler.ApplyForLoan)
		loans.GET("", loanHandler.GetUserLoans)
		loans.GET("/:id", loanHandler.GetLoan)
		loans.POST("/:id/pay", loanHandler.PayEMI)
	}

	// Cards routes (protected)
	cards := api.Group("/cards")
	cards.Use(middleware.AuthMiddleware(tokenService))
	{
		cards.POST("", cardHandler.CreateCard)
		cards.GET("", cardHandler.GetCards)
		cards.GET("/:id", cardHandler.GetCardDetails)
		cards.POST("/:id/otp", cardHandler.RequestOTP)
		cards.POST("/:id/details", cardHandler.GetCardDetails)
		cards.POST("/:id/freeze", cardHandler.FreezeCard)
	}

	// Merchant & QR routes (protected)
	merchant := api.Group("/merchant")
	merchant.Use(middleware.AuthMiddleware(tokenService))
	{
		merchant.POST("/register", qrHandler.RegisterMerchant)
	}

	qr := api.Group("/qr")
	qr.Use(middleware.AuthMiddleware(tokenService))
	{
		qr.POST("/generate", qrHandler.GenerateQR)
		qr.POST("/pay", qrHandler.PayViaQR)
	}
}
