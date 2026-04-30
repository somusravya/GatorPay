package main

import (
	"log"

	"flowpay-trading/config"
	"flowpay-trading/handlers"
	"flowpay-trading/middleware"
	"flowpay-trading/models"
	"flowpay-trading/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// Initialize DB
	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.TradingAccount{}, &models.Trade{}, &models.Position{})

	// Initialize Services & Handlers
	stockService := services.NewStockService()
	tradingService := services.NewTradingService(stockService)
	
	stockHandler := handlers.NewStockHandler(stockService)
	tradingHandler := handlers.NewTradingHandler(tradingService)

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "trading service up"})
		})
		
		// Stock read-only routes
		stocks := v1.Group("/stocks")
		{
			stocks.GET("/search", stockHandler.Search)
			stocks.GET("/market-summary", stockHandler.MarketSummary)
			stocks.GET("/:symbol/quote", stockHandler.GetQuote)
			stocks.GET("/:symbol/details", stockHandler.GetDetails)
			stocks.GET("/:symbol/chart", stockHandler.GetChart)
		}

		// Trading Routes (Protected)
		trading := v1.Group("/trading")
		trading.Use(middleware.AuthMiddleware())
		{
			trading.POST("/verify", tradingHandler.Verify)
			trading.GET("/account", tradingHandler.GetAccount)
			trading.POST("/trade", tradingHandler.ExecuteTrade)
			trading.GET("/portfolio", tradingHandler.GetPortfolio)
			trading.GET("/orders", tradingHandler.GetOrderHistory)
		}
	}

	log.Println("Starting Trading Microservice on port 8081...")
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to start service: %v", err)
	}
}
