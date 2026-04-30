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
	// Sprint 4 handlers
	insightHandler *handlers.InsightHandler,
	budgetHandler *handlers.BudgetHandler,
	subscriptionHandler *handlers.SubscriptionHandler,
	fraudHandler *handlers.FraudHandler,
	notificationHandler *handlers.NotificationHandler,
	socialHandler *handlers.SocialHandler,
	adminHandler *handlers.AdminHandler,
	invoiceHandler *handlers.InvoiceHandler,
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
		loans.GET("/eligibility", loanHandler.CheckEligibility)
		loans.POST("/apply", loanHandler.ApplyForLoan)
		loans.GET("", loanHandler.GetUserLoans)
		loans.GET("/:id", loanHandler.GetLoan)
		loans.POST("/:id/pay", loanHandler.PayEMI)
		loans.POST("/:id/cancel", loanHandler.CancelLoan)
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
		// Sprint 4: Merchant invoicing & payment links
		merchant.POST("/invoices", invoiceHandler.CreateInvoice)
		merchant.GET("/invoices", invoiceHandler.GetInvoices)
		merchant.POST("/payment-links", invoiceHandler.CreatePaymentLink)
		merchant.GET("/payment-links", invoiceHandler.GetPaymentLinks)
	}

	qr := api.Group("/qr")
	qr.Use(middleware.AuthMiddleware(tokenService))
	{
		qr.POST("/generate", qrHandler.GenerateQR)
		qr.POST("/lookup", qrHandler.LookupQR)
		qr.POST("/pay", qrHandler.PayViaQR)
	}

	// ═══════════════════════════════════════════
	// SPRINT 4: New Route Groups
	// ═══════════════════════════════════════════

	// AI Financial Insights routes (protected)
	insights := api.Group("/insights")
	insights.Use(middleware.AuthMiddleware(tokenService))
	{
		insights.GET("/summary", insightHandler.GetSummary)
	}

	// Budgeting & Auto-Save routes (protected)
	budget := api.Group("/budget")
	budget.Use(middleware.AuthMiddleware(tokenService))
	{
		budget.POST("/goals", budgetHandler.CreateGoal)
		budget.GET("/goals", budgetHandler.GetGoals)
	}

	autosave := api.Group("/autosave")
	autosave.Use(middleware.AuthMiddleware(tokenService))
	{
		autosave.POST("/rules", budgetHandler.CreateAutoSaveRule)
		autosave.POST("/roundup/execute", budgetHandler.ExecuteRoundup)
	}

	// Subscription routes (protected)
	subscriptions := api.Group("/subscriptions")
	subscriptions.Use(middleware.AuthMiddleware(tokenService))
	{
		subscriptions.GET("", subscriptionHandler.GetSubscriptions)
		subscriptions.POST("/track", subscriptionHandler.TrackSubscription)
		subscriptions.POST("/autopay", subscriptionHandler.SetAutoPay)
	}

	// Fraud Detection routes (protected)
	fraud := api.Group("/fraud")
	fraud.Use(middleware.AuthMiddleware(tokenService))
	{
		fraud.GET("/alerts", fraudHandler.GetAlerts)
		fraud.POST("/review", fraudHandler.ReviewAlert)
	}

	// Notifications routes (protected)
	notifications := api.Group("/notifications")
	notifications.Use(middleware.AuthMiddleware(tokenService))
	{
		notifications.GET("", notificationHandler.GetNotifications)
		notifications.GET("/preferences", notificationHandler.GetPreferences)
		notifications.PUT("/:id/read", notificationHandler.MarkRead)
		notifications.PUT("/preferences", notificationHandler.UpdatePreferences)
	}

	// Social Payments routes (protected)
	social := api.Group("/social")
	social.Use(middleware.AuthMiddleware(tokenService))
	{
		social.GET("/feed", socialHandler.GetFeed)
		social.POST("/post", socialHandler.CreatePost)
		social.POST("/react", socialHandler.ReactToPost)
		social.GET("/friends", socialHandler.GetFriends)
		social.POST("/friends/add", socialHandler.AddFriend)
	}

	// Admin Analytics routes (protected)
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(tokenService))
	{
		admin.GET("/metrics", adminHandler.GetMetrics)
		admin.GET("/users", adminHandler.GetUsers)
		admin.GET("/fraud/review", adminHandler.GetFraudReview)
	}
}
